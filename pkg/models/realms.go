package models

type Realm struct {
	BaseModel

	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `json:"posts"`
	AccountID   uint   `json:"account_id"`
}
