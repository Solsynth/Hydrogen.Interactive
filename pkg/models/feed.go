package models

type Feed struct {
	BaseModel

	Alias       string `json:"alias"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	ModelType   string `json:"model_type"`

	CommentCount  int64 `json:"comment_count"`
	ReactionCount int64 `json:"reaction_count"`

	AuthorID uint  `json:"author_id"`
	RealmID  *uint `json:"realm_id"`

	Author Account `json:"author" gorm:"embedded"`

	Attachments  []Attachment     `json:"attachments" gorm:"-"`
	ReactionList map[string]int64 `json:"reaction_list" gorm:"-"`
}
