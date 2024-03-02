package models

import "time"

// Account profiles basically fetched from Hydrogen.Identity
// But cache at here for better usage
// At the same time this model can make relations between local models
type Account struct {
	BaseModel

	Name             string           `json:"name"`
	Nick             string           `json:"nick"`
	Avatar           string           `json:"avatar"`
	Description      string           `json:"description"`
	EmailAddress     string           `json:"email_address"`
	PowerLevel       int              `json:"power_level"`
	Moments          []Moment         `json:"moments" gorm:"foreignKey:AuthorID"`
	Articles         []Article        `json:"articles" gorm:"foreignKey:AuthorID"`
	Attachments      []Attachment     `json:"attachments" gorm:"foreignKey:AuthorID"`
	LikedMoments     []MomentLike     `json:"liked_moments"`
	DislikedMoments  []MomentDislike  `json:"disliked_moments"`
	LikedArticles    []ArticleLike    `json:"liked_articles"`
	DislikedArticles []ArticleDislike `json:"disliked_articles"`
	LikedComments    []CommentLike    `json:"liked_comments"`
	DislikedComments []CommentDislike `json:"disliked_comments"`
	RealmIdentities  []RealmMember    `json:"identities"`
	Realms           []Realm          `json:"realms"`
	ExternalID       uint             `json:"external_id"`
}

type AccountMembership struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Follower    Account   `json:"follower"`
	Following   Account   `json:"following"`
	FollowerID  uint
	FollowingID uint
}
