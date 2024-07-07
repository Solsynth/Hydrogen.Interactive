package models

type Tag struct {
	BaseModel

	Alias       string    `json:"alias" gorm:"uniqueIndex" validate:"lowercase,alphanum"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Posts       []Post    `json:"posts" gorm:"many2many:post_tags"`
	Articles    []Article `json:"articles" gorm:"many2many:article_tags"`
}

type Category struct {
	BaseModel

	Alias       string    `json:"alias" gorm:"uniqueIndex" validate:"lowercase,alphanum"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Posts       []Post    `json:"posts" gorm:"many2many:post_categories"`
	Articles    []Article `json:"articles" gorm:"many2many:article_categories"`
}
