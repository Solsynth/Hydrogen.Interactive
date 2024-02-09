package models

type Realm struct {
	BaseModel

	Name        string        `json:"name"`
	Description string        `json:"description"`
	Posts       []Post        `json:"posts"`
	Members     []RealmMember `json:"members"`
	IsPublic    bool          `json:"is_public"`
	AccountID   uint          `json:"account_id"`
}

type RealmMember struct {
	BaseModel

	RealmID   uint `json:"realm_id"`
	AccountID uint `json:"account_id"`
}
