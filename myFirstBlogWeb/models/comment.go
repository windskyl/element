package models

import "time"

type Comment struct {
	CommentID  string    `bson:"comment_id"`
	Content    string    `bson:"content"`
	AuthorID   string    `bson:"author_id"`
	ArticleID  string    `bson:"article_id"`
	CreateTime time.Time `bson:"create_time"`
	ModifyTime time.Time `bson:"modify_time"`
}

type CommentRequest struct {
	Content   string `json:"content" binding:"required"`
	ArticleID string `json:"article_id" binding:"required"`
}
