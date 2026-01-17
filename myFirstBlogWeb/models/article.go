package models

import "time"

type Article struct {
	ArticleID      string    `bson:"article_id" json:"article_id"`
	Title          string    `bson:"title" json:"title"`
	Content        string    `bson:"content" json:"content"`
	AuthorID       string    `bson:"author_id" json:"author_id"`
	Images         []string  `bson:"images,omitempty" json:"images,omitempty"`
	AuthorUsername string    `bson:"-" json:"author_username,omitempty"`
	CreateTime     time.Time `bson:"create_time" json:"create_time"`
	ModifyTime     time.Time `bson:"modify_time" json:"modify_time"`
	CommentCount   int       `bson:"-" json:"comment_count,omitempty"`
}

type ArticleRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images,omitempty"`
}
