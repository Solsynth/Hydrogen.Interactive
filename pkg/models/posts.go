package models

import "time"

type Post struct {
	BaseModel

	Content          string        `json:"content"`
	Hashtags         []Tag         `json:"tags" gorm:"many2many:post_tags"`
	Categories       []Category    `json:"categories" gorm:"many2many:post_categories"`
	Attachments      []Attachment  `json:"attachments"`
	LikedAccounts    []PostLike    `json:"liked_accounts"`
	DislikedAccounts []PostDislike `json:"disliked_accounts"`
	RepostTo         *Post         `json:"repost_to" gorm:"foreignKey:RepostID"`
	ReplyTo          *Post         `json:"reply_to" gorm:"foreignKey:ReplyID"`
	PinnedAt         *time.Time    `json:"pinned_at"`
	EditedAt         *time.Time    `json:"edited_at"`
	PublishedAt      time.Time     `json:"published_at"`
	RepostID         *uint         `json:"repost_id"`
	ReplyID          *uint         `json:"reply_id"`
	RealmID          *uint         `json:"realm_id"`
	AuthorID         uint          `json:"author_id"`
	Author           Account       `json:"author"`

	// Dynamic Calculating Values
	LikeCount    int64 `json:"like_count" gorm:"-"`
	DislikeCount int64 `json:"dislike_count" gorm:"-"`
	ReplyCount   int64 `json:"reply_count" gorm:"-"`
	RepostCount  int64 `json:"repost_count" gorm:"-"`
}
