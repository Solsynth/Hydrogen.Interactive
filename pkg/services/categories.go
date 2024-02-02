package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"errors"
	"gorm.io/gorm"
)

func GetCategory(alias, name string) (models.Category, error) {
	var category models.Category
	if err := database.C.Where(models.Category{Alias: alias}).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			category = models.Category{
				Alias: alias,
				Name:  name,
			}
			err := database.C.Save(&category).Error
			return category, err
		}
		return category, err
	}
	return category, nil
}

func GetTag(alias, name string) (models.Tag, error) {
	var tag models.Tag
	if err := database.C.Where(models.Category{Alias: alias}).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = models.Tag{
				Alias: alias,
				Name:  name,
			}
			err := database.C.Save(&tag).Error
			return tag, err
		}
		return tag, err
	}
	return tag, nil
}
