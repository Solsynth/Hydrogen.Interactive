package models

type PostMetric struct {
	ReplyCount    int64            `json:"reply_count"`
	ReactionCount int64            `json:"reaction_count"`
	ReactionList  map[string]int64 `json:"reaction_list,omitempty"`
}
