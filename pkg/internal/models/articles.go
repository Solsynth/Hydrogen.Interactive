package models

import (
	"time"

	"gorm.io/datatypes"
)

type Article struct {
	BaseModel

	Alias       string                    `json:"alias" gorm:"uniqueIndex"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Content     string                    `json:"content"`
	Language    string                    `json:"language"`
	Tags        []Tag                     `json:"tags" gorm:"many2many:article_tags"`
	Categories  []Category                `json:"categories" gorm:"many2many:article_categories"`
	Reactions   []Reaction                `json:"reactions"`
	Attachments datatypes.JSONSlice[uint] `json:"attachments"`
	RealmID     *uint                     `json:"realm_id"`
	Realm       *Realm                    `json:"realm"`

	IsDraft     bool       `json:"is_draft"`
	PublishedAt *time.Time `json:"published_at"`

	AuthorID uint    `json:"author_id"`
	Author   Account `json:"author"`

	Metric PostMetric `json:"metric" gorm:"-"`
}
