// 日期格式化工具
export const formatDate = (dateString, format = 'YYYY-MM-DD HH:mm') => {
  const date = new Date(dateString)
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
}

// 截取字符串并添加省略号
export const truncate = (text, maxLength = 100) => {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

// 生成模拟数据
export const generateMockArticles = (count = 5, page = 1) => {
  return Array.from({ length: count }, (_, i) => ({
    article_id: `art-${page}-${i+1}`,
    title: `文章标题 ${(page-1)*count + i+1}`,
    content: '这是一篇示例文章内容，展示博客系统的文章显示功能。在实际应用中，这里将显示真实的文章内容。',
    author_id: 'admin',
    create_time: new Date(Date.now() - (i * 24 * 60 * 60 * 1000)).toISOString(),
    modify_time: new Date(Date.now() - (i * 12 * 60 * 60 * 1000)).toISOString()
  }))
}

export const generateMockComments = (articleId, count = 3) => {
  return Array.from({ length: count }, (_, i) => ({
    comment_id: `com-${articleId}-${i+1}`,
    content: `这是一条示例评论 ${i+1}，展示博客系统的评论功能。`,
    author_id: `user${i+1}`,
    article_id: articleId,
    create_time: new Date(Date.now() - (i * 2 * 60 * 60 * 1000)).toISOString(),
    modify_time: new Date(Date.now() - (i * 60 * 60 * 1000)).toISOString()
  }))
}