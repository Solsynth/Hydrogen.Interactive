package models

import "git.solsynth.dev/hydrogen/dealer/pkg/hyper"

type Tag struct {
	hyper.BaseModel

	Alias       string `json:"alias" gorm:"uniqueIndex" validate:"lowercase"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `json:"posts" gorm:"many2many:post_tags"`
}

type Category struct {
	hyper.BaseModel

	Alias       string `json:"alias" gorm:"uniqueIndex" validate:"lowercase,alphanum"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `json:"posts" gorm:"many2many:post_categories"`
}
