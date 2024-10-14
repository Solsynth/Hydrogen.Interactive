package models

import "git.solsynth.dev/hydrogen/dealer/pkg/hyper"

type Subscription struct {
	hyper.BaseModel

	FollowerID uint     `json:"follower_id"`
	Follower   Account  `json:"follower"`
	AccountID  *uint    `json:"account_id,omitempty"`
	Account    *Account `json:"account,omitempty"`
	TagID      *uint    `json:"tag_id,omitempty"`
	Tag        Tag      `json:"tag,omitempty"`
	CategoryID *uint    `json:"category_id,omitempty"`
	Category   Category `json:"category,omitempty"`
	RealmID    *uint    `json:"realm_id,omitempty"`
	Realm      Realm    `json:"realm,omitempty"`
}
