package models

// Account profiles basically fetched from Hydrogen.Passport
// But cache at here for better usage
// At the same time this model can make relations between local models
type Account struct {
	BaseModel

	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	EmailAddress string `json:"email_address"`
	PowerLevel   int    `json:"power_level"`
	ExternalID   uint   `json:"external_id"`
}
