package handlers

import (
	"fmt"
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
		fmt.Printf("Login payload: %+v\n", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求", "reason": "参数格式错误"})
		return
	}
	fmt.Printf("Login payload: %+v\n", req)

	// 使用不消耗验证码的校验，避免用户先验证图片然后因为用户名/密码错误导致验证码被意外消费
	if !utils.CheckCaptcha(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误", "reason": "验证码不正确"})
		return
	}

	collection := db.GetCollection("users")
	var user models.User
	filter := bson.M{"username": req.Username}
	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		fmt.Printf("Login user lookup error: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在", "reason": "数据库无此用户"})
		return
	}

	// 比较服务端存储的密码哈希（加盐）和客户端发送的密码哈希
	if user.PasswordHash != req.PasswordHash {
		fmt.Printf("Login password mismatch: stored='%s' provided='%s'\n", user.PasswordHash, req.PasswordHash)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误", "reason": "密码哈希不匹配"})
		return
	}

	// 登录成功后再消费验证码，确保用户验证与登录动作一致
	if ok := utils.VerifyCaptcha(req.CaptchaID, req.Captcha); !ok {
		fmt.Printf("Warning: VerifyCaptcha failed when consuming after successful login for id='%s' answer='%s'\n", req.CaptchaID, req.Captcha)
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
		fmt.Printf("format error\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求", "reason": "参数格式错误"})
		return
	}
	fmt.Printf("Register payload: %+v\n", req)

	if !utils.ValidateUsername(req.Username) {
		fmt.Printf("username rule error\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名格式错误", "reason": "用户名不符合规则"})
		return
	}

	// 使用可注入的校验函数（注册时也先做非消费校验）
	if !utils.CheckCaptcha(req.CaptchaID, req.Captcha) {
		fmt.Printf("captcha error\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误", "reason": "验证码不正确"})
		return
	}

	collection := db.GetCollection("users")
	count, _ := collection.CountDocuments(c, bson.M{"username": req.Username})
	if count > 0 {
		fmt.Printf("username existed error\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在", "reason": "数据库已存在该用户名"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败", "reason": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"userID":  user.UserID,
	})
}
