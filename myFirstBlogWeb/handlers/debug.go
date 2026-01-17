package handlers

import (
	"log"
	"net/http"

	"myFirstBlogWeb/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// DebugArticlesCount 返回 articles 集合的当前文档数（临时调试用）
func DebugArticlesCount(c *gin.Context) {
	collection := db.GetCollection("articles")
	total, err := collection.CountDocuments(c, bson.M{})
	if err != nil {
		log.Printf("DebugArticlesCount: CountDocuments error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败", "reason": err.Error()})
		return
	}

	log.Printf("DebugArticlesCount: total=%d", total)
	c.JSON(http.StatusOK, gin.H{"count": total})
}

// GetMyArticlesCount 返回当前登录用户的文章数量（需要 AuthMiddleware）
func GetMyArticlesCount(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	collection := db.GetCollection("articles")
	filter := bson.M{"author_id": userID.(string)}
	total, err := collection.CountDocuments(c, filter)
	if err != nil {
		log.Printf("GetMyArticlesCount: CountDocuments error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败", "reason": err.Error()})
		return
	}

	log.Printf("GetMyArticlesCount: user=%s count=%d", userID.(string), total)
	c.JSON(http.StatusOK, gin.H{"count": total})
}
