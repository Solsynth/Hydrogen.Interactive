package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
)

func ListRealm() ([]models.Realm, error) {
	var realms []models.Realm
	if err := database.C.Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func ListRealmWithUser(user models.Account) ([]models.Realm, error) {
	var realms []models.Realm
	if err := database.C.Where(&models.Realm{AccountID: user.ID}).Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func NewRealm(user models.Account, name, description string, isPublic bool) (models.Realm, error) {
	realm := models.Realm{
		Name:        name,
		Description: description,
		AccountID:   user.ID,
		IsPublic:    isPublic,
		Members: []models.RealmMember{
			{AccountID: user.ID},
		},
	}

	err := database.C.Save(&realm).Error

	return realm, err
}

func InviteRealmMember(user models.Account, target models.Realm) error {
	member := models.RealmMember{
		RealmID:   target.ID,
		AccountID: user.ID,
	}

	err := database.C.Save(&member).Error

	return err
}

func KickRealmMember(user models.Account, target models.Realm) error {
	var member models.RealmMember

	if err := database.C.Where(&models.RealmMember{
		RealmID:   target.ID,
		AccountID: user.ID,
	}).First(&member).Error; err != nil {
		return err
	}

	return database.C.Delete(&member).Error
}

func EditRealm(realm models.Realm, name, description string, isPublic bool) (models.Realm, error) {
	realm.Name = name
	realm.Description = description
	realm.IsPublic = isPublic

	err := database.C.Save(&realm).Error

	return realm, err
}

func DeleteRealm(realm models.Realm) error {
	return database.C.Delete(&realm).Error
}
