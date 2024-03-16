package models

type Article struct {
	PostBase

	Title       string       `json:"title"`
	Hashtags    []Tag        `json:"tags" gorm:"many2many:article_tags"`
	Categories  []Category   `json:"categories" gorm:"many2many:article_categories"`
	Reactions   []Reaction   `json:"reactions"`
	Attachments []Attachment `json:"attachments"`
	Description string       `json:"description"`
	Content     string       `json:"content"`
	RealmID     *uint        `json:"realm_id"`
	Realm       *Realm       `json:"realm"`

	Comments []Comment `json:"comments" gorm:"foreignKey:ArticleID"`
}

func (p *Article) GetReplyTo() PostInterface {
	return nil
}

func (p *Article) GetRepostTo() PostInterface {
	return nil
}

func (p *Article) GetHashtags() []Tag {
	return p.Hashtags
}

func (p *Article) GetCategories() []Category {
	return p.Categories
}

func (p *Article) SetHashtags(tags []Tag) {
	p.Hashtags = tags
}

func (p *Article) SetCategories(categories []Category) {
	p.Categories = categories
}
