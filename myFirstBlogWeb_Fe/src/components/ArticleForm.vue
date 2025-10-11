<template>
  <div>
    <form @submit.prevent="submit">
      <div class="form-group">
        <label>标题</label>
        <input type="text" v-model="formData.title" placeholder="请输入文章标题" required>
      </div>
      <div class="form-group">
        <label>内容</label>
        <textarea v-model="formData.content" placeholder="请输入文章内容" rows="8" required></textarea>
      </div>
      <div v-if="errorMessage" class="alert alert-error">{{ errorMessage }}</div>
      <div v-if="successMessage" class="alert alert-success">{{ successMessage }}</div>
      <button type="submit" class="btn btn-success btn-block">
        <i :class="isEditMode ? 'fas fa-save' : 'fas fa-plus'"></i>
        {{ isEditMode ? '保存修改' : '发布文章' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits } from 'vue'

const props = defineProps({
  article: Object,
  isEditMode: Boolean
})

const emit = defineEmits(['submitted'])

const formData = ref({
  title: '',
  content: ''
})
const errorMessage = ref('')
const successMessage = ref('')

watch(() => props.article, (newArticle) => {
  if (newArticle) {
    formData.value.title = newArticle.title || ''
    formData.value.content = newArticle.content || ''
  }
}, { immediate: true })

const submit = async () => {
  if (!formData.value.title || !formData.value.content) {
    errorMessage.value = '标题和内容不能为空'
    successMessage.value = ''
    return
  }
  emit('submitted', { ...formData.value })
  successMessage.value = props.isEditMode ? '修改成功' : '发布成功'
  errorMessage.value = ''
}
</script>