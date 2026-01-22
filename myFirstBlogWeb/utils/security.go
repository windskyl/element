package utils

import (
	"crypto/rand"
	"encoding/hex"
	"html"
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// 1. 移除所有HTML标签
	clean := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")

	// 2. HTML转义特殊字符
	clean = html.EscapeString(clean)

	// 3. 移除JavaScript事件属性
	clean = regexp.MustCompile(`\s+on\w+\s*=\s*["'][^"']*["']`).ReplaceAllString(clean, "")

	// 4. 移除多余空格
	return strings.TrimSpace(clean)
}

// SanitizeContent 清理用户提交的文章内容：移除 script 标签和事件属性(onXXX)，保留 img 等基本标签
func SanitizeContent(input string) string {
	// 移除 script 标签及其内容（忽略大小写和换行）
	reScript := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	clean := reScript.ReplaceAllString(input, "")

	// 移除 onXXX 事件属性
	reOn := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*("[^"]*"|'[^']*'|[^>\s]+)`)
	clean = reOn.ReplaceAllString(clean, "")

	// 仍然去除危险的 javascript: 协议在 href/src 等属性中
	reJs := regexp.MustCompile(`(?i)javascript:\s*`)
	clean = reJs.ReplaceAllString(clean, "")

	return strings.TrimSpace(clean)
}

func GenerateSalt() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func ValidateUsername(username string) bool {
	// 允许 2-20 个字母、数字或下划线（支持短用户名，如 "11"）
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{2,20}$`)
	return re.MatchString(username)
}
