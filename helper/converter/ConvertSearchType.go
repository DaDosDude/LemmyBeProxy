package converter

import (
	"LemmyBeProxy/dto/model/lemmy"
	"LemmyBeProxy/dto/model/piefed"
)

func ConvertSearchType(in piefed.SearchType) lemmy.SearchType {
	return lemmy.SearchType(in)
}

// ReverseConvertSearchType: PieFed has no "All" content type. Lemmy's
// SearchTypeAll is mapped to Posts as the most useful single-type default —
// a Lemmy client requesting "All" won't get communities/users/comments back,
// only posts, since PieFed's /search endpoint requires picking one type_.
//
// An empty/unset type is mapped the same way. Unlike real Lemmy, which
// treats type_ as optional, PieFed rejects both an absent and an empty
// type_ with "Must be one of: Communities, Posts, Users, Url, Comments."
// Clients that omit type_ entirely (lemmyBB does) would otherwise get a
// 400 on every search.
func ReverseConvertSearchType(in lemmy.SearchType) piefed.SearchType {
	if in == lemmy.SearchTypeAll || in == "" {
		return piefed.SearchTypePosts
	}

	return piefed.SearchType(in)
}
