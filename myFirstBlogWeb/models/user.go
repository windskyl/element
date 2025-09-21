package models

import "time"

type User struct {
	UserID       string    `bson:"user_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"password_hash"`
	Salt         string    `bson:"salt"`
	RegTime      time.Time `bson:"reg_time"`
	LastLogin    time.Time `bson:"last_login"`
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
