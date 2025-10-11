package handlers

import (
	"context"
	"myFirstBlogWeb/db"
	"myFirstBlogWeb/models"
	"myFirstBlogWeb/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateComment(c *gin.Context) {
	articleID := c.Param("articleId")

	var req models.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求", "reason": "Invalid request body"})
		return
	}

	if req.ArticleID != articleID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章ID不匹配", "reason": "Article ID mismatch"})
		return
	}

	comment := models.Comment{
		CommentID:  utils.GenerateUUID(),
		Content:    utils.SanitizeInput(req.Content),
		AuthorID:   c.MustGet("userID").(string),
		ArticleID:  articleID,
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
	}

	collection := db.GetCollection("comments")
	if _, err := collection.InsertOne(c, comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败", "reason": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	articleID := c.Param("articleId")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"create_time": -1})
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	filter := bson.M{"article_id": articleID}
	collection := db.GetCollection("comments")
	cursor, err := collection.Find(c, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败", "reason": "Failed to get comments"})
		return
	}
	defer cursor.Close(context.Background())

	var comments []models.Comment
	if err := cursor.All(c, &comments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析评论失败", "reason": "Failed to parse comments"})
		return
	}

	total, _ := collection.CountDocuments(c, filter)

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}
