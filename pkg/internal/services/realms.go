package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func GetRealmWithExtID(id uint) (models.Realm, error) {
	var realm models.Realm
	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return realm, err
	}
	response, err := proto.NewRealmClient(pc).GetRealm(context.Background(), &proto.LookupRealmRequest{
		Id: lo.ToPtr(uint64(id)),
	})
	if err != nil {
		return realm, err
	}
	return LinkRealm(response)
}

func GetRealmWithAlias(alias string) (models.Realm, error) {
	var realm models.Realm
	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return realm, err
	}
	response, err := proto.NewRealmClient(pc).GetRealm(context.Background(), &proto.LookupRealmRequest{
		Alias: &alias,
	})
	if err != nil {
		return realm, err
	}
	return LinkRealm(response)
}

func GetRealmMember(realmId uint, userId uint) (*proto.RealmMemberInfo, error) {
	var realm models.Realm
	if err := database.C.Where("id = ?", realmId).First(&realm).Error; err != nil {
		return nil, err
	}
	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return nil, err
	}
	response, err := proto.NewRealmClient(pc).GetRealmMember(context.Background(), &proto.RealmMemberLookupRequest{
		RealmId: uint64(realm.ID),
		UserId:  lo.ToPtr(uint64(userId)),
	})
	if err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func LinkRealm(info *proto.RealmInfo) (models.Realm, error) {
	var realm models.Realm
	if info == nil {
		return realm, fmt.Errorf("remote realm info was not found")
	}
	if err := database.C.Where("id = ?", info.Id).First(&realm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			realm = models.Realm{
				BaseRealm: hyper.LinkRealm(info),
			}
			return realm, database.C.Save(&realm).Error
		}
		return realm, err
	}

	prev := realm
	realm.Alias = info.Alias
	realm.Name = info.Name
	realm.Description = info.Description
	realm.Avatar = info.Avatar
	realm.Banner = info.Banner
	realm.IsPublic = info.IsPublic
	realm.IsCommunity = info.IsCommunity

	var err error
	if !reflect.DeepEqual(prev, realm) {
		err = database.C.Save(&realm).Error
	}

	return realm, err
}
