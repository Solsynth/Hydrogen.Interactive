package models

type PostMetric struct {
	ReplyCount    int64            `json:"reply_count,omitempty"`
	ReactionCount int64            `json:"reaction_count,omitempty"`
	ReactionList  map[string]int64 `json:"reaction_list,omitempty"`
}
