package models

import (
	"time"

	"gorm.io/datatypes"
)

const (
	PostTypeStory   = "story"
	PostTypeArticle = "article"
)

type PostVisibilityLevel = int8

const (
	PostVisibilityAll = PostVisibilityLevel(iota)
	PostVisibilityFriends
	PostVisibilityFiltered
	PostVisibilitySelected
	PostVisibilityNone
)

type Post struct {
	BaseModel

	Type       string            `json:"type"`
	Body       datatypes.JSONMap `json:"body"`
	Language   string            `json:"language"`
	Alias      *string           `json:"alias"`
	AreaAlias  *string           `json:"area_alias"`
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

	VisibleUsers   datatypes.JSONSlice[uint] `json:"visible_users_list"`
	InvisibleUsers datatypes.JSONSlice[uint] `json:"invisible_users_list"`
	Visibility     PostVisibilityLevel       `json:"visibility"`

	EditedAt *time.Time `json:"edited_at"`
	PinnedAt *time.Time `json:"pinned_at"`
	LockedAt *time.Time `json:"locked_at"`

	IsDraft        bool       `json:"is_draft"`
	PublishedAt    *time.Time `json:"published_at"`
	PublishedUntil *time.Time `json:"published_until"`

	TotalUpvote   int `json:"total_upvote"`
	TotalDownvote int `json:"total_downvote"`

	AuthorID uint    `json:"author_id"`
	Author   Account `json:"author"`

	Metric PostMetric `json:"metric" gorm:"-"`
}

type PostStoryBody struct {
	Thumbnail   *uint   `json:"thumbnail"`
	Title       *string `json:"title"`
	Content     string  `json:"content"`
	Location    *string `json:"location"`
	Attachments []uint  `json:"attachments"`
}

type PostArticleBody struct {
	Thumbnail   *uint   `json:"thumbnail"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Content     string  `json:"content"`
	Attachments []uint  `json:"attachments"`
}
