package models

import (
	"time"
)

type PostReactInfo struct {
	PostID       uint  `json:"post_id"`
	LikeCount    int64 `json:"like_count"`
	DislikeCount int64 `json:"dislike_count"`
	ReplyCount   int64 `json:"reply_count"`
	RepostCount  int64 `json:"repost_count"`
}

type PostBase struct {
	BaseModel

	Alias       string       `json:"alias" gorm:"uniqueIndex"`
	Attachments []Attachment `json:"attachments"`
	PublishedAt *time.Time   `json:"published_at"`

	AuthorID uint    `json:"author_id"`
	Author   Account `json:"author"`

	// Dynamic Calculated Values
	ReactionList map[string]int64 `json:"reaction_list" gorm:"-"`
}

func (p *PostBase) GetID() uint {
	return p.ID
}

func (p *PostBase) GetReplyTo() PostInterface {
	return nil
}

func (p *PostBase) GetRepostTo() PostInterface {
	return nil
}

func (p *PostBase) GetAuthor() Account {
	return p.Author
}

func (p *PostBase) GetRealm() *Realm {
	return nil
}

func (p *PostBase) SetReactionList(list map[string]int64) {
	p.ReactionList = list
}

type PostInterface interface {
	GetID() uint
	GetHashtags() []Tag
	GetCategories() []Category
	GetReplyTo() PostInterface
	GetRepostTo() PostInterface
	GetAuthor() Account
	GetRealm() *Realm

	SetHashtags([]Tag)
	SetCategories([]Category)
	SetReactionList(map[string]int64)
}
