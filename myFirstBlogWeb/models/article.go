package models

import "time"

type Article struct {
	ArticleID  string    `bson:"article_id"`
	Title      string    `bson:"title"`
	Content    string    `bson:"content"`
	AuthorID   string    `bson:"author_id"`
	CreateTime time.Time `bson:"create_time"`
	ModifyTime time.Time `bson:"modify_time"`
}

type ArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
