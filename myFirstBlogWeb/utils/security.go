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

func GenerateSalt() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func ValidateUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return re.MatchString(username)
}
