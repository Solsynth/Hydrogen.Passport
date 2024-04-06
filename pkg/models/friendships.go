package models

type FriendshipStatus = int8

const (
	FriendshipPending = FriendshipStatus(iota)
	FriendshipActive
	FriendshipBlocked
)

type AccountFriendship struct {
	BaseModel

	AccountID uint             `json:"account_id"`
	RelatedID uint             `json:"related_id"`
	BlockedBy *uint            `json:"blocked_by"`
	Account   Account          `json:"account"`
	Related   Account          `json:"related"`
	Status    FriendshipStatus `json:"status"`
}
