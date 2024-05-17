package models

import (
	"gorm.io/datatypes"
	"time"
)

type PostReactInfo struct {
	PostID       uint  `json:"post_id"`
	LikeCount    int64 `json:"like_count"`
	DislikeCount int64 `json:"dislike_count"`
	ReplyCount   int64 `json:"reply_count"`
	RepostCount  int64 `json:"repost_count"`
}

type Post struct {
	BaseModel

	Alias       string                      `json:"alias" gorm:"uniqueIndex"`
	Content     string                      `json:"content"`
	Tags        []Tag                       `json:"tags" gorm:"many2many:post_tags"`
	Categories  []Category                  `json:"categories" gorm:"many2many:post_categories"`
	Reactions   []Reaction                  `json:"reactions"`
	Replies     []Post                      `json:"replies" gorm:"foreignKey:ReplyID"`
	Attachments datatypes.JSONSlice[string] `json:"attachments"`
	ReplyID     *uint                       `json:"reply_id"`
	RepostID    *uint                       `json:"repost_id"`
	RealmID     *uint                       `json:"realm_id"`
	ReplyTo     *Post                       `json:"reply_to" gorm:"foreignKey:ReplyID"`
	RepostTo    *Post                       `json:"repost_to" gorm:"foreignKey:RepostID"`
	Realm       *Realm                      `json:"realm"`

	PublishedAt *time.Time `json:"published_at"`

	AuthorID uint    `json:"author_id"`
	Author   Account `json:"author"`

	// Dynamic Calculated Values
	ReplyCount    int64            `json:"comment_count"`
	ReactionCount int64            `json:"reaction_count"`
	ReactionList  map[string]int64 `json:"reaction_list" gorm:"-"`
}
