package models

type Comment struct {
	PostBase

	Content    string     `json:"content"`
	Hashtags   []Tag      `json:"tags" gorm:"many2many:comment_tags"`
	Categories []Category `json:"categories" gorm:"many2many:comment_categories"`
	Reactions  []Reaction `json:"reactions"`
	ReplyID    *uint      `json:"reply_id"`
	ReplyTo    *Comment   `json:"reply_to" gorm:"foreignKey:ReplyID"`

	ArticleID *uint    `json:"article_id"`
	MomentID  *uint    `json:"moment_id"`
	Article   *Article `json:"article"`
	Moment    *Moment  `json:"moment"`
}

func (p *Comment) GetReplyTo() PostInterface {
	return p.ReplyTo
}

func (p *Comment) GetHashtags() []Tag {
	return p.Hashtags
}

func (p *Comment) GetCategories() []Category {
	return p.Categories
}

func (p *Comment) SetHashtags(tags []Tag) {
	p.Hashtags = tags
}

func (p *Comment) SetCategories(categories []Category) {
	p.Categories = categories
}
