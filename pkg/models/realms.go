package models

// Realm profiles basically fetched from Hydrogen.Passport
// But cache at here for better usage and database relations
type Realm struct {
	BaseModel

	Alias       string    `json:"alias"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Articles    []Article `json:"article"`
	Moments     []Moment  `json:"moments"`
	IsPublic    bool      `json:"is_public"`
	IsCommunity bool      `json:"is_community"`
	ExternalID  uint      `json:"external_id"`
}
