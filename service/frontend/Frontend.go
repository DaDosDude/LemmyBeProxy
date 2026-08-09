package frontend

import (
	lemmyRequest "LemmyBeProxy/dto/request/lemmy"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	"LemmyBeProxy/http"
)

// Frontend is the counterpart to backend.Backend, but for the other end
// of the pipeline: instead of choosing what this proxy talks to, it
// chooses what wire format this proxy accepts from and replies to
// clients with. Every method parses/builds using this proxy's canonical
// (current Lemmy 0.19.x-shaped) request/response types — the same ones
// backend.Backend works in — so a controller using both interfaces never
// needs to know or care which frontend or backend version is actually in
// play on either side of it.
//
// This interface covers Post, Comment, Community, User (fully, including
// SaveUserSettings), Search, and Site (including MyUser) methods so
// far, mirroring how backend.Backend started with Post and Comment.
//
// Known gap: Frontend017's BuildSuccessResponse (used for mark_as_read)
// returns the canonical {success: bool} shape rather than the full
// PostResponse{post_view} real Lemmy 0.17.x actually returns there,
// since the canonical MarkPostAsRead call doesn't carry post_view data
// and building it accurately would need an extra backend round-trip —
// a real design decision deferred rather than rushed. Every other method
// in this slice is fully accurate to 0.17.2, confirmed against Lemmy's
// own source at that tag.
type Frontend interface {
	ParseGetPostsRequest(request *http.Request) (*lemmyRequest.GetPostsRequest, error)
	BuildGetPostsResponse(resp *lemmyResponse.GetPostsResponse) any

	ParseGetPostRequest(request *http.Request) (*lemmyRequest.GetPostRequest, error)
	BuildGetPostResponse(resp *lemmyResponse.GetPostResponse) any

	ParseCreatePostRequest(request *http.Request) (*lemmyRequest.CreatePostRequest, error)
	ParseEditPostRequest(request *http.Request) (*lemmyRequest.EditPostRequest, error)
	ParseCreatePostLikeRequest(request *http.Request) (*lemmyRequest.CreatePostLikeRequest, error)
	// BuildPostMutationResponse is shared by create/edit/like — the
	// canonical backend unifies all three into GetPostResponse, matching
	// how modern Lemmy actually behaves.
	BuildPostMutationResponse(resp *lemmyResponse.GetPostResponse) any

	ParseMarkPostAsReadRequest(request *http.Request) (*lemmyRequest.MarkPostAsReadRequest, error)
	BuildSuccessResponse(resp *lemmyResponse.SuccessResponse) any

	ParseGetCommentsRequest(request *http.Request) (*lemmyRequest.GetCommentsRequest, error)
	BuildGetCommentsResponse(resp *lemmyResponse.GetCommentsResponse) any

	ParseGetCommentRequest(request *http.Request) (*lemmyRequest.GetCommentRequest, error)
	BuildGetCommentResponse(resp *lemmyResponse.GetCommentResponse) any

	ParseCreateCommentRequest(request *http.Request) (*lemmyRequest.CreateCommentRequest, error)
	BuildCreateCommentResponse(resp *lemmyResponse.CreateCommentResponse) any

	ParseCreateCommentLikeRequest(request *http.Request) (*lemmyRequest.CreateCommentLikeRequest, error)
	BuildLikeCommentResponse(resp *lemmyResponse.GetCommentResponse) any

	ParseGetCommunityRequest(request *http.Request) (*lemmyRequest.GetCommunityRequest, error)
	BuildGetCommunityResponse(resp *lemmyResponse.GetCommunityResponse) any

	ParseGetCommunitiesRequest(request *http.Request) (*lemmyRequest.GetCommunitiesRequest, error)
	BuildGetCommunitiesResponse(resp *lemmyResponse.GetCommunitiesResponse) any

	ParseFollowCommunityRequest(request *http.Request) (*lemmyRequest.FollowCommunityRequest, error)
	BuildCommunityResponse(resp *lemmyResponse.CommunityResponse) any

	ParseBlockCommunityRequest(request *http.Request) (*lemmyRequest.BlockCommunityRequest, error)
	BuildBlockCommunityResponse(resp *lemmyResponse.BlockCommunityResponse) any

	ParseLoginRequest(request *http.Request) (*lemmyRequest.LoginRequest, error)
	BuildLoginResponse(resp *lemmyResponse.LoginResponse) any

	BuildGetUnreadCountResponse(resp *lemmyResponse.GetUnreadCountResponse) any

	ParseGetUserRequest(request *http.Request) (*lemmyRequest.GetUserRequest, error)
	BuildGetUserResponse(resp *lemmyResponse.GetUserResponse) any

	ParseBlockPersonRequest(request *http.Request) (*lemmyRequest.BlockPersonRequest, error)
	BuildBlockPersonResponse(resp *lemmyResponse.BlockPersonResponse) any

	ParseSearchRequest(request *http.Request) (*lemmyRequest.SearchRequest, error)
	BuildSearchResponse(resp *lemmyResponse.SearchResponse) any

	// BuildGetSiteResponse now covers MyUser too — the logged-in user's
	// own subscriptions, blocks, and moderated communities.
	BuildGetSiteResponse(resp *lemmyResponse.GetSiteResponse) any

	ParseSaveUserSettingsRequest(request *http.Request) (*lemmyRequest.SaveUserSettingsRequest, error)
	BuildSaveUserSettingsResponse(resp *lemmyResponse.SaveUserSettingsResponse) any

	ParseGetPersonMentionsRequest(request *http.Request) (*lemmyRequest.GetPersonMentionsRequest, error)
	BuildGetPersonMentionsResponse(resp *lemmyResponse.GetPersonMentionsResponse) any

	ParseGetRepliesRequest(request *http.Request) (*lemmyRequest.GetRepliesRequest, error)
	BuildGetRepliesResponse(resp *lemmyResponse.GetRepliesResponse) any

	ParseGetPrivateMessagesRequest(request *http.Request) (*lemmyRequest.GetPrivateMessagesRequest, error)
	BuildGetPrivateMessagesResponse(resp *lemmyResponse.GetPrivateMessagesResponse) any

	ParseResolveObjectRequest(request *http.Request) (*lemmyRequest.ResolveObjectRequest, error)
	BuildResolveObjectResponse(resp *lemmyResponse.ResolveObjectResponse) any

	BuildGetCaptchaResponse(resp *lemmyResponse.GetCaptchaResponse) any
}
