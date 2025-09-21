package middleware

import (
	"myFirstBlogWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut {
			// 确保表单已解析
			if c.Request.Form == nil {
				c.Request.ParseForm()
			}

			// 遍历所有POST/PUT参数
			for key, values := range c.Request.PostForm {
				// 为每个值执行净化
				for i, value := range values {
					sanitized := utils.SanitizeInput(value)
					c.Request.PostForm[key][i] = sanitized
				}
			}
		}
		c.Next()
	}
}
