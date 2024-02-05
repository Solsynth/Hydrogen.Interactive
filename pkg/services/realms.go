package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
)

func ListRealms(user models.Account) ([]models.Realm, error) {
	var realms []models.Realm
	if err := database.C.Where(&models.Realm{AccountID: user.ID}).Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func NewRealm(user models.Account, name, description string) (models.Realm, error) {
	realm := models.Realm{
		Name:        name,
		Description: description,
		AccountID:   user.ID,
	}

	err := database.C.Save(&realm).Error

	return realm, err
}

func EditRealm(realm models.Realm, name, description string) (models.Realm, error) {
	realm.Name = name
	realm.Description = description

	err := database.C.Save(&realm).Error

	return realm, err
}

func DeleteRealm(realm models.Realm) error {
	return database.C.Delete(&realm).Error
}
