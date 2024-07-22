package models

type Tag struct {
	BaseModel

	Alias       string `json:"alias" gorm:"uniqueIndex" validate:"lowercase"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `json:"posts" gorm:"many2many:post_tags"`
}

type Category struct {
	BaseModel

	Alias       string `json:"alias" gorm:"uniqueIndex" validate:"lowercase,alphanum"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `json:"posts" gorm:"many2many:post_categories"`
}
