package models

type RealmType = int

const (
	RealmTypePublic = RealmType(iota)
	RealmTypeRestricted
	RealmTypePrivate
)

type Realm struct {
	BaseModel

	Name        string        `json:"name"`
	Description string        `json:"description"`
	Articles    []Article     `json:"article"`
	Moments     []Moment      `json:"moments"`
	Members     []RealmMember `json:"members"`
	RealmType   RealmType     `json:"realm_type"`
	AccountID   uint          `json:"account_id"`
}

type RealmMember struct {
	BaseModel

	RealmID   uint    `json:"realm_id"`
	AccountID uint    `json:"account_id"`
	Realm     Realm   `json:"realm"`
	Account   Account `json:"account"`
}
