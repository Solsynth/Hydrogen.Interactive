package models

type Tag struct {
	BaseModel

	Alias       string    `json:"alias" gorm:"uniqueIndex" validate:"lowercase,alphanum,min=4,max=24"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Articles    []Article `json:"articles" gorm:"many2many:article_tags"`
	Moments     []Moment  `json:"moments" gorm:"many2many:moment_tags"`
	Comments    []Comment `json:"comments" gorm:"many2many:comment_tags"`
}

type Category struct {
	BaseModel

	Alias       string    `json:"alias" gorm:"uniqueIndex" validate:"lowercase,alphanum,min=4,max=24"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Articles    []Article `json:"articles" gorm:"many2many:article_categories"`
	Moments     []Moment  `json:"moments" gorm:"many2many:moment_categories"`
	Comments    []Comment `json:"comments" gorm:"many2many:comment_categories"`
}
