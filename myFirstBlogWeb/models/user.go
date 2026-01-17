package models

import "time"

type User struct {
	UserID       string    `bson:"user_id" json:"user_id"`
	Username     string    `bson:"username" json:"username"`
	PasswordHash string    `bson:"password_hash" json:"password_hash"`
	Salt         string    `bson:"salt" json:"salt"`
	RegTime      time.Time `bson:"reg_time" json:"reg_time"`
	LastLogin    time.Time `bson:"last_login" json:"last_login"`
}

type LoginRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Captcha      string `json:"captcha"`
	CaptchaID    string `json:"captcha_id"`
}

type RegisterRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Captcha      string `json:"captcha"`
	CaptchaID    string `json:"captcha_id"`
}
