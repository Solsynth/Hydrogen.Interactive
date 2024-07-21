package models

import (
	"time"

	"gorm.io/datatypes"
)

type Post struct {
	BaseModel

	Body       datatypes.JSONMap `json:"body"`
	Language   string            `json:"language"`
	Tags       []Tag             `json:"tags" gorm:"many2many:post_tags"`
	Categories []Category        `json:"categories" gorm:"many2many:post_categories"`
	Reactions  []Reaction        `json:"reactions"`
	Replies    []Post            `json:"replies" gorm:"foreignKey:ReplyID"`
	ReplyID    *uint             `json:"reply_id"`
	RepostID   *uint             `json:"repost_id"`
	RealmID    *uint             `json:"realm_id"`
	ReplyTo    *Post             `json:"reply_to" gorm:"foreignKey:ReplyID"`
	RepostTo   *Post             `json:"repost_to" gorm:"foreignKey:RepostID"`
	Realm      *Realm            `json:"realm"`

	IsDraft     bool       `json:"is_draft"`
	PublishedAt *time.Time `json:"published_at"`

	AuthorID uint    `json:"author_id"`
	Author   Account `json:"author"`

	Metric PostMetric `json:"metric" gorm:"-"`
}

type PostStoryBody struct {
	Title       *string `json:"title"`
	Content     string  `json:"content"`
	Location    *string `json:"location"`
	Attachments []uint  `json:"attachments"`
}

type PostArticleBody struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Content     string  `json:"content"`
	Attachments []uint  `json:"attachments"`
}
