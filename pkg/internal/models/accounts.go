package models

import "git.solsynth.dev/hydrogen/dealer/pkg/hyper"

type Account struct {
	hyper.BaseUser

	Posts         []Post         `json:"posts" gorm:"foreignKey:AuthorID"`
	Reactions     []Reaction     `json:"reactions"`
	Subscriptions []Subscription `json:"subscriptions" gorm:"foreginKey:FollowerID"`

	TotalUpvote   int `json:"total_upvote"`
	TotalDownvote int `json:"total_downvote"`
}
