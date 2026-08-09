package lemmy017

import "LemmyBeProxy/dto/model/lemmy"

// PostView matches Lemmy 0.17.2's real PostView field set exactly —
// confirmed against Lemmy's own source. Community and Post are reused
// directly from the canonical model since both were checked field-by-field
// against 0.17.2 and found to already contain every field 0.17.x
// requires. Creator and Counts use the 0.17-specific types in this
// package instead, since those two genuinely differ.
type PostView struct {
	Post                       lemmy.Post           `json:"post" validate:"required"`
	Creator                    Person               `json:"creator" validate:"required"`
	Community                  lemmy.Community      `json:"community" validate:"required"`
	CreatorBannedFromCommunity bool                 `json:"creator_banned_from_community" validate:"required"`
	Counts                     PostAggregates       `json:"counts" validate:"required"`
	Subscribed                 lemmy.SubscribedType `json:"subscribed" validate:"required"`
	Saved                      bool                 `json:"saved" validate:"required"`
	Read                       bool                 `json:"read" validate:"required"`
	CreatorBlocked             bool                 `json:"creator_blocked" validate:"required"`
	MyVote                     *int                 `json:"my_vote"`
	UnreadComments             int64                `json:"unread_comments" validate:"required"`
}
