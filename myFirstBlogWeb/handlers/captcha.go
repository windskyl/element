package handlers

import (
	"myFirstBlogWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifyCaptchaHandler 校验验证码，前端可以调用此接口来验证给定 captcha_id 与 captcha 是否匹配
func VerifyCaptchaHandler(c *gin.Context) {
	var req struct {
		CaptchaID string `json:"captcha_id" binding:"required"`
		Captcha   string `json:"captcha" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	ok := utils.CheckCaptcha(req.CaptchaID, req.Captcha)
	c.JSON(http.StatusOK, gin.H{"ok": ok})
}
