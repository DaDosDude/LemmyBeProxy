package backend

import (
	lemmyRequest "LemmyBeProxy/dto/request/lemmy"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	"LemmyBeProxy/http"
)

// Backend is the contract every backend implementation (Piefed, real Lemmy,
// and any future addition) must satisfy. Every method works entirely in
// Lemmy-shaped request/response types — that's already this proxy's
// canonical internal shape, since every controller was built to speak it
// regardless of what's actually behind it. A Piefed implementation
// translates to/from Piefed's real shape internally; a real-Lemmy
// implementation is close to a passthrough, since the shapes already
// match.
//
// This interface now covers every endpoint this proxy implements — Post,
// Comment, Community, User, Search, Site, Upload, and notifications
// (mentions, replies, private messages).
type Backend interface {
	GetPosts(request *lemmyRequest.GetPostsRequest, headers http.Headers) (*lemmyResponse.GetPostsResponse, error)
	GetPost(request *lemmyRequest.GetPostRequest, headers http.Headers) (*lemmyResponse.GetPostResponse, error)
	CreatePost(request *lemmyRequest.CreatePostRequest, headers http.Headers) (*lemmyResponse.GetPostResponse, error)
	EditPost(request *lemmyRequest.EditPostRequest, headers http.Headers) (*lemmyResponse.GetPostResponse, error)
	LikePost(request *lemmyRequest.CreatePostLikeRequest, headers http.Headers) (*lemmyResponse.GetPostResponse, error)
	MarkPostAsRead(request *lemmyRequest.MarkPostAsReadRequest, headers http.Headers) (*lemmyResponse.SuccessResponse, error)

	GetComments(request *lemmyRequest.GetCommentsRequest, headers http.Headers) (*lemmyResponse.GetCommentsResponse, error)
	GetComment(request *lemmyRequest.GetCommentRequest, headers http.Headers) (*lemmyResponse.GetCommentResponse, error)
	CreateComment(request *lemmyRequest.CreateCommentRequest, headers http.Headers) (*lemmyResponse.CreateCommentResponse, error)
	LikeComment(request *lemmyRequest.CreateCommentLikeRequest, headers http.Headers) (*lemmyResponse.GetCommentResponse, error)

	GetCommunity(request *lemmyRequest.GetCommunityRequest, headers http.Headers) (*lemmyResponse.GetCommunityResponse, error)
	GetCommunities(request *lemmyRequest.GetCommunitiesRequest, headers http.Headers) (*lemmyResponse.GetCommunitiesResponse, error)
	FollowCommunity(request *lemmyRequest.FollowCommunityRequest, headers http.Headers) (*lemmyResponse.CommunityResponse, error)
	BlockCommunity(request *lemmyRequest.BlockCommunityRequest, headers http.Headers) (*lemmyResponse.BlockCommunityResponse, error)

	Login(request *lemmyRequest.LoginRequest, headers http.Headers) (*lemmyResponse.LoginResponse, error)
	Register(request *lemmyRequest.RegisterRequest, headers http.Headers) (*lemmyResponse.LoginResponse, error)
	GetUnreadCount(headers http.Headers) (*lemmyResponse.GetUnreadCountResponse, error)
	GetUser(request *lemmyRequest.GetUserRequest, headers http.Headers) (*lemmyResponse.GetUserResponse, error)
	BlockPerson(request *lemmyRequest.BlockPersonRequest, headers http.Headers) (*lemmyResponse.BlockPersonResponse, error)
	SaveUserSettings(request *lemmyRequest.SaveUserSettingsRequest, headers http.Headers) (*lemmyResponse.SaveUserSettingsResponse, error)

	Search(request *lemmyRequest.SearchRequest, headers http.Headers) (*lemmyResponse.SearchResponse, error)

	// Site takes no request DTO — Lemmy's GET /site has no parameters
	// beyond auth. PiefedBackend's implementation does real work here
	// (an ActivityPub actor fetch to supplement fields Piefed's own API
	// doesn't expose); LemmyBackend's is close to a passthrough, since
	// real Lemmy's own /site response is already complete on its own.
	Site(headers http.Headers) (*lemmyResponse.GetSiteResponse, error)

	// UploadImage returns the resulting image's real, browser-accessible
	// URL — nothing more. UploadController's hex-token-and-redirect
	// mechanism (built for Piefed, whose own upload API doesn't return a
	// directly servable URL) works identically for either backend from
	// that URL alone, so it stays in the controller rather than being
	// duplicated per backend. jwt is a raw token, not a Bearer/Cookie
	// header — each backend applies it however its own upload endpoint
	// actually expects auth.
	UploadImage(fileBytes []byte, filename string, jwt string) (string, error)

	// GetPersonMentions, GetReplies, and GetPrivateMessages are fetched
	// on every authenticated page load by real Lemmy clients (lemmyBB
	// included) to populate notification/message counts — a client
	// crashes outright if any of the three 404s. Piefed doesn't have
	// person mentions, comment replies, or private messages as concepts
	// in the same shape, so PiefedBackend's implementation of these is
	// necessarily limited — see its own comments.
	GetPersonMentions(request *lemmyRequest.GetPersonMentionsRequest, headers http.Headers) (*lemmyResponse.GetPersonMentionsResponse, error)
	GetReplies(request *lemmyRequest.GetRepliesRequest, headers http.Headers) (*lemmyResponse.GetRepliesResponse, error)
	GetPrivateMessages(request *lemmyRequest.GetPrivateMessagesRequest, headers http.Headers) (*lemmyResponse.GetPrivateMessagesResponse, error)

	// ResolveObject looks up a fediverse handle or ActivityPub URL
	// (e.g. "!community@instance.tld" or "https://instance.tld/c/name")
	// and returns whichever of comment/post/community/person it
	// resolved to. lemmyBB uses this to resolve every community listed
	// in an admin's lemmybb_categories.hjson on the frontpage.
	ResolveObject(request *lemmyRequest.ResolveObjectRequest, headers http.Headers) (*lemmyResponse.ResolveObjectResponse, error)

	// GetCaptcha is called unconditionally by real Lemmy clients
	// (lemmyBB included) loading the registration form, regardless of
	// whether captcha is actually enabled for the site — a missing
	// route here blocks registration entirely, not just captcha
	// display. See GetCaptchaResponse's own comment for why this
	// always returns Ok: nil.
	GetCaptcha(headers http.Headers) (*lemmyResponse.GetCaptchaResponse, error)
}
