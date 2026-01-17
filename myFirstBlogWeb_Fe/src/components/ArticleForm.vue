<template>
  <div>
    <form @submit.prevent="submit">
      <div class="form-group">
        <label>标题</label>
        <input :disabled="isSubmitting" type="text" v-model="formData.title" placeholder="请输入文章标题" required>
      </div>
      <div class="form-group">
        <label class="label-strong">内容</label>
        <textarea ref="contentArea" :disabled="isSubmitting" v-model="formData.content" placeholder="请输入文章内容" rows="8" required></textarea>
      </div>
          <div class="form-group">
            <label>图片</label>
            <div class="image-controls">
              <input ref="fileInput" type="file" accept="image/*" multiple @change="onFileChange" :disabled="isSubmitting" style="display:none" />
              <button class="btn btn-plus" type="button" @click="() => fileInput.click()">＋ 选择图片</button>
            </div>
              <div class="image-list">
                <div v-for="(imgObj, idx) in images" :key="imgObj.url" class="image-tile">
                  <img :src="imgObj.url" alt="img" />
                        <!-- 缩略图下方移除宽度/插入/删除按钮，交互放到模态中统一管理 -->
                        <div class="image-actions">
                        </div>
                </div>
              </div>
            <small>你也可以粘贴图片到内容区 (Ctrl+V)</small>
          </div>
      <!-- Modal for selecting width before insertion -->
      <div v-if="showModal" class="modal-backdrop">
        <div class="modal">
          <h3>选择图片显示宽度</h3>
          <div class="modal-images">
            <div v-for="(imgObj, idx) in pendingImages" :key="imgObj.url" style="width:140px; text-align:center">
              <img :src="imgObj.url" style="max-width:100%; max-height:120px; object-fit:cover" />
              <div style="margin-top:8px">
                <select v-model.number="imgObj.width">
                  <option v-for="w in [100,75,50,25]" :key="w" :value="w">{{ w }}%</option>
                </select>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-subtle" type="button" @click="cancelPending">取消</button>
            <button class="btn btn-success" type="button" @click="confirmInsertPending">确定插入</button>
          </div>
        </div>
      </div>
      <!-- 本地校验错误 -->
      <div v-if="localError" class="alert alert-error">{{ localError }}</div>
      <!-- 来自父组件的网络错误 -->
      <div v-if="error" class="alert alert-error">{{ error }}</div>
      <button type="submit" class="btn btn-success btn-block" :disabled="isSubmitting">
        <i :class="isEditMode ? 'fas fa-save' : 'fas fa-plus'"></i>
        {{ isEditMode ? '保存修改' : (isSubmitting ? '提交中...' : '发布文章') }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits, onMounted, onBeforeUnmount, nextTick } from 'vue'
import api from '@/services/api'

const props = defineProps({
  article: Object,
  isEditMode: Boolean,
  error: { type: String, default: '' },
  isSubmitting: { type: Boolean, default: false }
})

const emit = defineEmits(['submitted'])

const formData = ref({
  title: '',
  content: ''
})
const images = ref([]) // each item: { url: string, width: number }
const localError = ref('')
const contentArea = ref(null)
const fileInput = ref(null)
const pendingImages = ref([]) // images uploaded but pending insertion (references to objects in images)
const showModal = ref(false)

watch(() => props.article, (newArticle) => {
  if (newArticle) {
    formData.value.title = newArticle.title || ''
    formData.value.content = newArticle.content || ''
    images.value = newArticle.images ? newArticle.images.map(u => ({ url: u, width: 100 })) : []
  }
}, { immediate: true })

onMounted(() => {
  // 调试输出，帮助定位为什么输入不可编辑
  console.debug('ArticleForm mounted', { articleProp: props.article, isSubmitting: props.isSubmitting, formData: formData.value })
  // 监听粘贴事件以支持粘贴图片上传
  if (contentArea.value) {
    contentArea.value.addEventListener('paste', handlePaste)
  }
})

// 清理监听
onBeforeUnmount(() => {
  if (contentArea.value) contentArea.value.removeEventListener('paste', handlePaste)
})

const handlePaste = async (e) => {
  const items = e.clipboardData && e.clipboardData.items
  if (!items) return
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.kind === 'file' && item.type.startsWith('image/')) {
      const blob = item.getAsFile()
      await uploadAndInsert(blob)
    }
  }
}

// 上传文件并在上传完成后打开选择宽度的模态框
const uploadFilesAndOpenModal = async (fileList) => {
  const uploaded = []
  for (let i = 0; i < fileList.length; i++) {
    const file = fileList[i]
    try {
      const fd = new FormData()
      fd.append('file', file)
      const res = await api.uploadImage(fd)
      const url = res.data.url
      const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'
      const full = url.startsWith('http') ? url : (base + url)
      const obj = { url: full, width: 100 }
      images.value.push(obj)
      uploaded.push(obj)
    } catch (err) {
      console.error('上传图片失败', err)
      localError.value = '上传图片失败'
    }
  }
  if (uploaded.length > 0) {
    pendingImages.value = uploaded
    showModal.value = true
  }
}

const submit = () => {
  localError.value = ''
  if (!formData.value.title || !formData.value.content) {
    localError.value = '标题和内容不能为空'
    return
  }
  // 仅触发提交事件，真正的网络请求和导航由父组件负责
  // 只传回图片 URL 列表（后端模型期望 []string）
  emit('submitted', { ...formData.value, images: images.value.map(i => i.url) })
}

const onFileChange = async (e) => {
  const files = e.target.files
  if (!files || files.length === 0) return
  await uploadFilesAndOpenModal(files)
  e.target.value = ''
}

const uploadAndInsert = async (file) => {
  try {
    const fd = new FormData()
    fd.append('file', file)
    const res = await api.uploadImage(fd)
    const url = res.data.url
    // 将上传返回的相对路径转换为绝对后缀（以后端为基准）
    const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'
    const full = url.startsWith('http') ? url : (base + url)
    images.value.push({ url: full, width: 100 })
    // 将图片加入待处理列表并打开模态由 onFileChange 统一处理
    pendingImages.value.push(images.value[images.value.length - 1])
    showModal.value = true
  } catch (err) {
    console.error('上传图片失败', err)
    localError.value = '上传图片失败'
  }
}

const insertImageAtCursor = (imgObj) => {
  const ta = contentArea.value
  const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'
  const full = imgObj.url.startsWith('http') ? imgObj.url : (base + imgObj.url)
  const style = 'style="width:' + (imgObj.width || 100) + '%;height:auto;"'
  if (!ta) {
    formData.value.content += '\n<img src="' + full + '" ' + style + ' />\n'
    return
  }
  const start = ta.selectionStart || 0
  const end = ta.selectionEnd || 0
  const before = formData.value.content.slice(0, start)
  const after = formData.value.content.slice(end)
  const inserted = '<img src="' + full + '" ' + style + ' />'
  formData.value.content = before + inserted + after
  // restore cursor after inserted text
  nextTick(() => {
    ta.selectionStart = ta.selectionEnd = start + inserted.length
    ta.focus()
  })
}

const removeImage = async (idx, imgObj) => {
  try {
    // 从 url 中提取文件名
    const parts = imgObj.url.split('/')
    const filename = parts[parts.length - 1]
    await api.deleteImage(filename)
  } catch (err) {
    console.warn('删除远程图片失败（可忽略）', err)
  }
  images.value.splice(idx, 1)
  // 同步从内容中移除所有该 url
  formData.value.content = formData.value.content.split(imgObj.url).join('')
}

const confirmInsertPending = async () => {
  if (!pendingImages.value || pendingImages.value.length === 0) {
    showModal.value = false
    return
  }
  if (!confirm('确定插入所选图片吗？')) return
  // 插入所有 pendingImages（保持顺序）
  const ta = contentArea.value
  // if cursor exists, insert first at cursor then append others after
  for (let i = 0; i < pendingImages.value.length; i++) {
    insertImageAtCursor(pendingImages.value[i])
  }
  // 清空 pending 并关闭模态
  pendingImages.value = []
  showModal.value = false
}

const cancelPending = async () => {
  // 删除已上传但未插入的文件以避免垃圾文件
  for (let i = 0; i < pendingImages.value.length; i++) {
    const parts = pendingImages.value[i].url.split('/')
    const filename = parts[parts.length - 1]
    try { await api.deleteImage(filename) } catch (e) { /* ignore */ }
    // 从 images 列表中移除该项
    const idx = images.value.findIndex(it => it.url === pendingImages.value[i].url)
    if (idx >= 0) images.value.splice(idx, 1)
  }
  pendingImages.value = []
  showModal.value = false
}
</script>

<style scoped>
.image-list{ display:flex; gap:8px; flex-wrap:wrap; margin-top:8px }
.image-tile{ width:160px; border:1px solid #eee; padding:8px; border-radius:6px; display:flex; flex-direction:column; align-items:center }
.image-tile img{ max-width:100%; max-height:120px; object-fit:cover }
.image-actions{ display:flex; gap:8px; margin-top:8px; align-items:center }
.image-actions .btn{ padding:8px 12px; min-width:88px }
.image-actions select{ min-width:84px; padding:6px 8px }
.label-strong{ font-weight:700 }
.form-group label{ font-weight:700 }
textarea[required]{ font-weight:400 }
.btn-plus{ background:#2b8aef; color:#fff; border-radius:6px; padding:8px 12px; border:none }
.btn-plus:hover{ background:#1f6fd6 }
.width-label{ font-weight:600; display:flex; align-items:center; gap:6px }
.form-group label{ font-weight:700 }
/* Ensure textarea content has normal weight */
.form-group textarea, textarea[required]{ font-weight:400 }

/* Modal */
.modal-backdrop{ position:fixed; inset:0; background:rgba(0,0,0,0.4); display:flex; align-items:center; justify-content:center; z-index:1200 }
.modal{ background:#fff; padding:16px; max-width:720px; width:90%; border-radius:8px; box-shadow:0 6px 24px rgba(0,0,0,0.2) }
.modal .modal-images{ display:flex; gap:12px; flex-wrap:wrap; max-height:60vh; overflow:auto }
.modal .modal-footer{ display:flex; justify-content:flex-end; gap:8px; margin-top:12px }
</style>