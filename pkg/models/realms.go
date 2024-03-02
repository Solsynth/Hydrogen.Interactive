package models

type Realm struct {
	BaseModel

	Name        string        `json:"name"`
	Description string        `json:"description"`
	Articles    []Article     `json:"article"`
	Moments     []Moment      `json:"moments"`
	Members     []RealmMember `json:"members"`
	IsPublic    bool          `json:"is_public"`
	AccountID   uint          `json:"account_id"`
}

type RealmMember struct {
	BaseModel

	RealmID   uint `json:"realm_id"`
	AccountID uint `json:"account_id"`
}
