package services

import (
	"context"
	"errors"
	"fmt"
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/grpc"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func GetRealm(id uint) (models.Realm, error) {
	var realm models.Realm
	response, err := grpc.Realms.GetRealm(context.Background(), &proto.RealmLookupRequest{
		Id: lo.ToPtr(uint64(id)),
	})
	if err != nil {
		return realm, err
	}
	return LinkRealm(response)
}

func GetRealmWithAlias(alias string) (models.Realm, error) {
	var realm models.Realm
	response, err := grpc.Realms.GetRealm(context.Background(), &proto.RealmLookupRequest{
		Alias: &alias,
	})
	if err != nil {
		return realm, err
	}
	return LinkRealm(response)
}

func GetRealmMember(realmId uint, userId uint) (*proto.RealmMemberResponse, error) {
	response, err := grpc.Realms.GetRealmMember(context.Background(), &proto.RealmMemberLookupRequest{
		RealmId: uint64(realmId),
		UserId:  lo.ToPtr(uint64(userId)),
	})
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func ListRealmMember(realmId uint) ([]*proto.RealmMemberResponse, error) {
	response, err := grpc.Realms.ListRealmMember(context.Background(), &proto.RealmMemberLookupRequest{
		RealmId: uint64(realmId),
	})
	if err != nil {
		return nil, err
	} else {
		return response.Data, nil
	}
}

func LinkRealm(info *proto.RealmResponse) (models.Realm, error) {
	var realm models.Realm
	if info == nil {
		return realm, fmt.Errorf("remote realm info was not found")
	}
	if err := database.C.Where(&models.Realm{
		ExternalID: uint(info.Id),
	}).First(&realm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			realm = models.Realm{
				Alias:       info.Alias,
				Name:        info.Name,
				Description: info.Description,
				IsPublic:    info.IsPublic,
				IsCommunity: info.IsCommunity,
				ExternalID:  uint(info.Id),
			}
			return realm, database.C.Save(&realm).Error
		}
		return realm, err
	}
	return realm, nil
}
