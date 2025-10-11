package handlers

import (
	"myFirstBlogWeb/config"
	"myFirstBlogWeb/db"
	"myFirstBlogWeb/models"
	"myFirstBlogWeb/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

// 新增：验证码校验函数变量，默认指向 utils.VerifyCaptcha
var VerifyCaptchaFunc = utils.VerifyCaptcha

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	if !VerifyCaptchaFunc(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	collection := db.GetCollection("users")
	var user models.User
	filter := bson.M{"username": req.Username}
	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 比较服务端存储的密码哈希（加盐）和客户端发送的密码哈希
	if user.PasswordHash != req.PasswordHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.UserID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.JWTSecret))

	c.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"userID":   user.UserID,
		"username": user.Username,
	})
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	if !utils.ValidateUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名格式错误"})
		return
	}

	// 使用可注入的校验函数
	if !VerifyCaptchaFunc(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	collection := db.GetCollection("users")
	count, _ := collection.CountDocuments(c, bson.M{"username": req.Username})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	user := models.User{
		UserID:       utils.GenerateUUID(),
		Username:     req.Username,
		PasswordHash: req.PasswordHash,
		RegTime:      time.Now(),
		LastLogin:    time.Now(),
	}

	if _, err := collection.InsertOne(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"userID":  user.UserID,
	})
}
