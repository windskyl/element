package handlers

import (
	"context"
	"log"
	"myFirstBlogWeb/db"
	"myFirstBlogWeb/models"
	"myFirstBlogWeb/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateArticle(c *gin.Context) {
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求", "reason": "Invalid request body"})
		return
	}

	article := models.Article{
		ArticleID:  utils.GenerateUUID(),
		Title:      utils.SanitizeInput(req.Title),
		Content:    utils.SanitizeContent(req.Content),
		AuthorID:   c.MustGet("userID").(string),
		CreateTime: time.Now(),
		ModifyTime: time.Now(),
		Images:     req.Images,
	}

	collection := db.GetCollection("articles")
	if _, err := collection.InsertOne(c, article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败", "reason": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.MustGet("userID").(string)

	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求", "reason": "Invalid request body"})
		return
	}

	filter := bson.M{
		"article_id": articleID,
		"author_id":  userID,
	}
	update := bson.M{"$set": bson.M{
		"title":       utils.SanitizeInput(req.Title),
		"content":     utils.SanitizeContent(req.Content),
		"modify_time": time.Now(),
		"images":      req.Images,
	}}

	collection := db.GetCollection("articles")
	result, err := collection.UpdateOne(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败", "reason": "Failed to update article"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或无权修改", "reason": "Article not found or no permission"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败", "reason": "Failed to delete article"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或无权删除", "reason": "Article not found or no permission"})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetArticle(c *gin.Context) {
	articleID := c.Param("id")

	collection := db.GetCollection("articles")
	var article models.Article
	if err := collection.FindOne(c, bson.M{"article_id": articleID}).Decode(&article); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在", "reason": "Article not found"})
		return
	}

	// 查询作者用户名（如果存在）
	usersColl := db.GetCollection("users")
	var user models.User
	authorUsername := ""
	if err := usersColl.FindOne(c, bson.M{"user_id": article.AuthorID}).Decode(&user); err == nil {
		authorUsername = user.Username
	}

	c.JSON(http.StatusOK, gin.H{"article": article, "author_username": authorUsername})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析文章失败", "reason": "Failed to decode articles"})
		return
	}
	// 确保返回空切片而不是 null，便于前端渲染
	if articles == nil {
		articles = make([]models.Article, 0)
	}

	total, err := collection.CountDocuments(c, bson.M{})
	if err != nil {
		log.Printf("GetArticles: CountDocuments error: %v", err)
		total = 0
	}

	// 临时调试日志：打印从 DB 查询到的条目数与 total，帮助排查前端/后端不一致问题
	log.Printf("GetArticles: returned %d articles slice, total in collection=%d, page=%d, limit=%d", len(articles), total, page, limit)

	// 为每篇文章附上作者用户名（如果可用）
	usersColl := db.GetCollection("users")
	for i := range articles {
		var u models.User
		if err := usersColl.FindOne(c, bson.M{"user_id": articles[i].AuthorID}).Decode(&u); err == nil {
			articles[i].AuthorUsername = u.Username
		}
	}

	// 使用单次聚合批量获取所有文章的评论计数，避免 N+1 查询
	commentsColl := db.GetCollection("comments")
	// 收集所有 article_id
	ids := make([]interface{}, 0, len(articles))
	for _, a := range articles {
		ids = append(ids, a.ArticleID)
	}
	if len(ids) > 0 {
		pipeline := mongo.Pipeline{
			bson.D{{"$match", bson.D{{"article_id", bson.D{{"$in", ids}}}}}},
			bson.D{{"$group", bson.D{{"_id", "$article_id"}, {"count", bson.D{{"$sum", 1}}}}}},
		}
		cur, err := commentsColl.Aggregate(c, pipeline)
		if err != nil {
			log.Printf("GetArticles: Aggregate comments error: %v", err)
		} else {
			var groups []struct {
				ID    string `bson:"_id"`
				Count int64  `bson:"count"`
			}
			if err := cur.All(c, &groups); err != nil {
				log.Printf("GetArticles: decode aggregate result error: %v", err)
			} else {
				// map counts by article_id
				cntMap := make(map[string]int64, len(groups))
				for _, g := range groups {
					cntMap[g.ID] = g.Count
				}
				for i := range articles {
					if v, ok := cntMap[articles[i].ArticleID]; ok {
						articles[i].CommentCount = int(v)
					} else {
						articles[i].CommentCount = 0
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}
