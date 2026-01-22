package main

import (
	"log"
	"myFirstBlogWeb/db"
	"myFirstBlogWeb/handlers"
	"myFirstBlogWeb/middleware"
	"myFirstBlogWeb/utils"

	"github.com/gin-gonic/gin"
)

// 此项目要实现一个简单的 Blog 系统
// 需要支持如下功能：
// 用户登录和注册（不需密码找回）。
// 用户发贴（不需要支持富文本，只需要支持纯文本）。
// 用户评论（不需要支持富文本，只需要支持纯文本）。
// 用户登录时的密码不应该保存为明文，应该用 MD5+Salt 来保存。
// 用户登录后，对于用户自己的贴子可以有“重新编辑”或 “删除”的功能，但是无权编辑或删除其它用户的贴子。
// 登录成功后，展示所有用户的帖子，按时间倒序排列，最新的在前面，按页显示，一页最多20个
// 点开一个帖子时，除了显示帖子内容，也显示评论区前20个，评论按页显示在帖子下方

// 图片验证码
// 阻止用户在发文章或评论时输入带 HTML 或 JavaScript 的内容。
// 用户注册，输入用户名和密码，注册成功后返回用户 ID。

func main() {
	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatal("数据库连接失败:", err)
		return
	}
	defer db.CloseDB()

	// 创建Gin实例
	r := gin.Default()
	// 注册中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SanitizeInput())

	// 图片验证码路由
	r.POST("/captcha", func(c *gin.Context) {
		id, b64 := utils.GenerateCaptcha()
		// 返回兼容字段名，前端多个地方可能期待不同字段名
		c.JSON(200, gin.H{"captcha_id": id, "image_base64": b64, "id": id, "image": b64})
	})
	// 验证验证码（前端可直接调用以判断输入是否匹配）
	r.POST("/captcha/verify", handlers.VerifyCaptchaHandler)

	// 静态文件：公开上传的图片
	r.Static("/uploads", "public/uploads")

	// 注册路由
	r.POST("/login", handlers.Login)       // 用户登录
	r.POST("/register", handlers.Register) // 用户注册

	// 临时调试路由 - 返回 articles 集合计数，便于前端快速验证 DB 状态
	r.GET("/debug/articles-count", handlers.DebugArticlesCount)

	authGroup := r.Group("/ctx")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/articles", handlers.CreateArticle)                // 创建文章
		authGroup.PUT("/articles/:id", handlers.UpdateArticle)             // 更新文章
		authGroup.DELETE("/articles/:id", handlers.DeleteArticle)          // 删除文章
		authGroup.GET("/articles", handlers.GetArticles)                   // 获取文章列表
		authGroup.GET("/user/articles-count", handlers.GetMyArticlesCount) // 获取当前用户文章计数
		authGroup.POST("/comments/:articleId", handlers.CreateComment)     // 创建评论
		authGroup.GET("/comments/:articleId", handlers.GetComments)        // 获取文章评论
		// 上传/删除图片（登录用户）
		authGroup.POST("/upload", handlers.UploadImage)
		authGroup.DELETE("/images/:name", handlers.DeleteImage)
		// 单个文章详情（公开访问）
	}

	// 公开文章详情路由
	r.GET("/articles/:id", handlers.GetArticle)
	// 公开文章列表（无需登录）
	r.GET("/articles", handlers.GetArticles)

	// 启动服务（改用 3000 以避开系统保留端口）
	log.Println("服务运行在 :3000")
	r.Run(":3000")
}
