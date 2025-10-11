package test

import (
	"bytes"
	"context"
	"encoding/json"
	"myFirstBlogWeb/db"
	"myFirstBlogWeb/handlers"
	"myFirstBlogWeb/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var cleanArticles = true // 控制是否清理 articles 表和测试数据

func cleanArticlesTable() {
	if cleanArticles {
		_, _ = db.GetCollection("articles").DeleteMany(context.Background(), bson.M{})
	}
}

func setupArticleRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// mock userID
	r.Use(func(c *gin.Context) {
		c.Set("userID", "testuser1")
		c.Next()
	})
	r.POST("/articles", handlers.CreateArticle)
	r.GET("/articles", handlers.GetArticles)
	r.PUT("/articles/:id", handlers.UpdateArticle)
	r.DELETE("/articles/:id", handlers.DeleteArticle)
	return r
}

func TestArticleCRUD(t *testing.T) {
	cleanArticlesTable()
	r := setupArticleRouter()

	// 创建文章
	articleReq := models.ArticleRequest{
		Title:   "Test Title",
		Content: "Test Content",
	}
	body, _ := json.Marshal(articleReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/articles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var article models.Article
	json.Unmarshal(w.Body.Bytes(), &article)
	assert.Equal(t, "Test Title", article.Title)
	assert.Equal(t, "Test Content", article.Content)
	assert.Equal(t, "testuser1", article.AuthorID)

	articleID := article.ArticleID

	// 修改文章
	updateReq := models.ArticleRequest{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	body, _ = json.Marshal(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/articles/"+articleID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 获取文章列表，检查修改
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/articles", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var listResp struct {
		Articles []models.Article `json:"articles"`
		Total    int64            `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	found := false
	for _, a := range listResp.Articles {
		if a.ArticleID == articleID {
			found = true
			assert.Equal(t, "Updated Title", a.Title)
			assert.Equal(t, "Updated Content", a.Content)
		}
	}
	assert.True(t, found, "未找到刚刚修改的文章")

	// 删除文章
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/articles/"+articleID, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 再次获取文章列表，确认已删除
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/articles", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &listResp)
	for _, a := range listResp.Articles {
		assert.NotEqual(t, articleID, a.ArticleID, "文章未被删除")
	}

	cleanArticlesTable()
}
