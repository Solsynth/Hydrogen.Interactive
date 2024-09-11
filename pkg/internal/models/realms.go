package models

import "git.solsynth.dev/hydrogen/dealer/pkg/hyper"

type Realm struct {
	hyper.BaseRealm

	Posts []Post `json:"posts"`
}
