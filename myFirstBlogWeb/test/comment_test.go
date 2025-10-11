package test

import (
	"context"
	"encoding/json"
	"myFirstBlogWeb/db"
	"myFirstBlogWeb/handlers"
	"myFirstBlogWeb/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var cleanComments = true // 控制是否清理 comments 表和测试数据

func cleanCommentsTable() {
	if cleanComments {
		_, _ = db.GetCollection("comments").DeleteMany(context.Background(), bson.M{})
	}
}

func TestGetComments_Normal(t *testing.T) {
	cleanCommentsTable()
	// 插入测试数据
	articleID := "test-article-1"
	comments := []interface{}{
		models.Comment{
			CommentID:  "c1",
			Content:    "hello",
			AuthorID:   "u1",
			ArticleID:  articleID,
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
		},
		models.Comment{
			CommentID:  "c2",
			Content:    "world",
			AuthorID:   "u2",
			ArticleID:  articleID,
			CreateTime: time.Now(),
			ModifyTime: time.Now(),
		},
	}
	_, err := db.GetCollection("comments").InsertMany(context.Background(), comments)
	if err != nil {
		t.Fatalf("插入测试评论失败: %v", err)
	}

	// 构造请求
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/comments/:articleId", handlers.GetComments)

	req, _ := http.NewRequest("GET", "/comments/"+articleID+"?page=1&limit=2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码200，实际%d", w.Code)
	}

	var resp struct {
		Comments []models.Comment `json:"comments"`
		Total    int64            `json:"total"`
		Page     int              `json:"page"`
		Limit    int              `json:"limit"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应解析失败: %v", err)
	}
	if resp.Total != 2 {
		t.Errorf("期望total=2，实际%d", resp.Total)
	}
	if len(resp.Comments) != 2 {
		t.Errorf("期望2条评论，实际%d", len(resp.Comments))
	}
	if resp.Page != 1 || resp.Limit != 2 {
		t.Errorf("分页参数错误: page=%d, limit=%d", resp.Page, resp.Limit)
	}

	cleanCommentsTable()
}

func TestGetComments_Empty(t *testing.T) {
	cleanCommentsTable()
	articleID := "no-comment-article"
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/comments/:articleId", handlers.GetComments)

	req, _ := http.NewRequest("GET", "/comments/"+articleID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码200，实际%d", w.Code)
	}
	var resp struct {
		Comments []models.Comment `json:"comments"`
		Total    int64            `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应解析失败: %v", err)
	}
	if resp.Total != 0 || len(resp.Comments) != 0 {
		t.Errorf("期望无评论，实际total=%d, len=%d", resp.Total, len(resp.Comments))
	}
	cleanCommentsTable()
}

func TestGetComments_DBError(t *testing.T) {
	// 关闭数据库，模拟错误
	db.CloseDB()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/comments/:articleId", handlers.GetComments)

	req, _ := http.NewRequest("GET", "/comments/someid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 由于没有连接数据库，应该返回500
	if w.Code != http.StatusInternalServerError {
		t.Errorf("期望500，实际%d", w.Code)
	}
}
