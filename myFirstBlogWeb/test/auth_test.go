package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"myFirstBlogWeb/db"
	"myFirstBlogWeb/handlers"
	"myFirstBlogWeb/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// 初始化数据库
	if err := db.InitDB(); err != nil {
		panic(err)
	}
	// 清空 users 集合，避免用户名冲突
	_ = db.GetCollection("users").Drop(context.Background())
	code := m.Run()
	db.CloseDB()
	os.Exit(code)
}

func mockVerifyCaptcha(id, answer string) bool { return true }

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	return r
}

func TestRegisterAndLogin(t *testing.T) {
	// 注入 mock 校验函数
	handlers.VerifyCaptchaFunc = mockVerifyCaptcha

	r := setupRouter()

	// 注册
	registerReq := models.RegisterRequest{
		Username:     "testuser1",
		PasswordHash: "testhash",
		Captcha:      "1234",
		CaptchaID:    "cid",
	}
	body, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var regResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &regResp)
	assert.Equal(t, "注册成功", regResp["message"])

	// 登录
	loginReq := models.LoginRequest{
		Username:     "testuser1",
		PasswordHash: "testhash",
		Captcha:      "1234",
		CaptchaID:    "cid",
	}
	body, _ = json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NotEmpty(t, loginResp["token"])
	assert.Equal(t, "testuser1", loginResp["username"])
}
