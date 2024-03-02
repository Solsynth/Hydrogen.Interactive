package models

import "time"

type CommentLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ArticleID *uint     `json:"article_id"`
	MomentID  *uint     `json:"moment_id"`
	CommentID *uint     `json:"comment_id"`
	AccountID uint      `json:"account_id"`
}

type CommentDislike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ArticleID *uint     `json:"article_id"`
	MomentID  *uint     `json:"moment_id"`
	CommentID *uint     `json:"comment_id"`
	AccountID uint      `json:"account_id"`
}

type MomentLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ArticleID *uint     `json:"article_id"`
	MomentID  *uint     `json:"moment_id"`
	CommentID *uint     `json:"comment_id"`
	AccountID uint      `json:"account_id"`
}

type MomentDislike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ArticleID *uint     `json:"article_id"`
	MomentID  *uint     `json:"moment_id"`
	CommentID *uint     `json:"comment_id"`
	AccountID uint      `json:"account_id"`
}

type ArticleLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ArticleID *uint     `json:"article_id"`
	MomentID  *uint     `json:"moment_id"`
	CommentID *uint     `json:"comment_id"`
	AccountID uint      `json:"account_id"`
}

type ArticleDislike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ArticleID *uint     `json:"article_id"`
	MomentID  *uint     `json:"moment_id"`
	CommentID *uint     `json:"comment_id"`
	AccountID uint      `json:"account_id"`
}
