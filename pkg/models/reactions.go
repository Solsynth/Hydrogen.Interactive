package models

import "time"

type PostLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostID    uint      `json:"post_id"`
	AccountID uint      `json:"account_id"`
}

type PostDislike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostID    uint      `json:"post_id"`
	AccountID uint      `json:"account_id"`
}
