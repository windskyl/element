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

func CreateArticle(c *gin.Context) {
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	article := models.Article{
		ArticleID:  utils.GenerateUUID(),
		Title:      utils.SanitizeInput(req.Title),
		Content:    utils.SanitizeInput(req.Content),
		AuthorID:   c.MustGet("userID").(string),
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
	}

	collection := db.GetCollection("articles")
	if _, err := collection.InsertOne(c, article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.MustGet("userID").(string)

	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	filter := bson.M{
		"article_id": articleID,
		"author_id":  userID,
	}
	update := bson.M{"$set": bson.M{
		"title":       utils.SanitizeInput(req.Title),
		"content":     utils.SanitizeInput(req.Content),
		"modify_time": time.Now(),
	}}

	collection := db.GetCollection("articles")
	result, err := collection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或无权修改"})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.MustGet("userID").(string)

	filter := bson.M{
		"article_id": articleID,
		"author_id":  userID,
	}

	collection := db.GetCollection("articles")
	result, err := collection.DeleteOne(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或无权删除"})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetArticles(c *gin.Context) {
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

	collection := db.GetCollection("articles")
	cursor, err := collection.Find(c, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}
	defer cursor.Close(context.Background())

	var articles []models.Article
	if err := cursor.All(c, &articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析文章失败"})
		return
	}

	total, _ := collection.CountDocuments(c, bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}
