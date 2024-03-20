package services

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"github.com/samber/lo"
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

func ListRealmIsAvailable(user models.Account) ([]models.Realm, error) {
	var realms []models.Realm
	var members []models.RealmMember
	if err := database.C.Where(&models.RealmMember{
		AccountID: user.ID,
	}).Find(&members).Error; err != nil {
		return realms, err
	}

	idx := lo.Map(members, func(item models.RealmMember, index int) uint {
		return item.RealmID
	})

	if err := database.C.Where(&models.Realm{
		RealmType: models.RealmTypePublic,
	}).Or("id IN ?", idx).Find(&realms).Error; err != nil {
		return realms, err
	}

	return realms, nil
}

func NewRealm(user models.Account, name, description string, realmType int) (models.Realm, error) {
	realm := models.Realm{
		Name:        name,
		Description: description,
		AccountID:   user.ID,
		RealmType:   realmType,
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

func EditRealm(realm models.Realm, name, description string, realmType int) (models.Realm, error) {
	realm.Name = name
	realm.Description = description
	realm.RealmType = realmType

	err := database.C.Save(&realm).Error

	return realm, err
}

func DeleteRealm(realm models.Realm) error {
	return database.C.Delete(&realm).Error
}
