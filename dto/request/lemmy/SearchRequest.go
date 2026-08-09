package lemmy

import "LemmyBeProxy/dto/model/lemmy"

// CreatorId is intentionally not included — PieFed's /search endpoint has
// no equivalent parameter, so searching by a specific user isn't supported
// through this proxy.
type SearchRequest struct {
	Q             string             `json:"q"`
	Type          lemmy.SearchType   `json:"type_,omitempty"`
	Limit         *uint              `json:"limit,omitempty"`
	ListingType   *lemmy.ListingType `json:"listing_type,omitempty"`
	Page          *uint              `json:"page,omitempty"`
	Sort          *lemmy.SortType    `json:"sort,omitempty"`
	CommunityName *string            `json:"community_name,omitempty"`
	CommunityId   *uint              `json:"community_id,omitempty"`
}
