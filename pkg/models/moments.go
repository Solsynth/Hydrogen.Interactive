package models

type Moment struct {
	PostBase

	Content     string       `json:"content"`
	Hashtags    []Tag        `json:"tags" gorm:"many2many:moment_tags"`
	Categories  []Category   `json:"categories" gorm:"many2many:moment_categories"`
	Reactions   []Reaction   `json:"reactions"`
	Attachments []Attachment `json:"attachments"`
	RealmID     *uint        `json:"realm_id"`
	RepostID    *uint        `json:"repost_id"`
	Realm       *Realm       `json:"realm"`
	RepostTo    *Moment      `json:"repost_to" gorm:"foreignKey:RepostID"`

	Comments []Comment `json:"comments" gorm:"foreignKey:MomentID"`
}

func (p *Moment) GetRepostTo() PostInterface {
	return p.RepostTo
}

func (p *Moment) GetRealm() *Realm {
	return p.Realm
}

func (p *Moment) GetHashtags() []Tag {
	return p.Hashtags
}

func (p *Moment) GetCategories() []Category {
	return p.Categories
}

func (p *Moment) SetHashtags(tags []Tag) {
	p.Hashtags = tags
}

func (p *Moment) SetCategories(categories []Category) {
	p.Categories = categories
}
