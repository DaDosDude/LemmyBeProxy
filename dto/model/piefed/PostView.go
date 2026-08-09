package piefed

type PostView struct {
	Post                       Post           `json:"post" validate:"required"`
	Creator                    Person         `json:"creator" validate:"required"`
	Community                  Community      `json:"community" validate:"required"`
	CreatorBannedFromCommunity bool           `json:"creator_banned_from_community" validate:"required"`
	BannedFromCommunity        bool           `json:"banned_from_community" validate:"required"`
	CreatorIsModerator         bool           `json:"creator_is_moderator" validate:"required"`
	CreatorIsAdmin             bool           `json:"creator_is_admin" validate:"required"`
	Counts                     PostAggregates `json:"counts" validate:"required"`
	Subscribed                 SubscribedType `json:"subscribed" validate:"required"`
	Saved                      bool           `json:"saved" validate:"required"`
	ActivityAlert              *bool          `json:"activity_alert,omitempty"`
	Read                       bool           `json:"read" validate:"required"`
	Hidden                     bool           `json:"hidden" validate:"required"`
	CreatorBlocked             *bool          `json:"creator_blocked,omitempty"`
	MyVote                     *int           `json:"my_vote,omitempty"`
	UnreadComments             int64          `json:"unread_comments" validate:"required"`
}
