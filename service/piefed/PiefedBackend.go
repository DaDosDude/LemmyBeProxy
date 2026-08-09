package piefed

import (
	lemmyModel "LemmyBeProxy/dto/model/lemmy"
	piefedModel "LemmyBeProxy/dto/model/piefed"
	lemmyRequest "LemmyBeProxy/dto/request/lemmy"
	piefedRequest "LemmyBeProxy/dto/request/piefed"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	piefedResponse "LemmyBeProxy/dto/response/piefed"
	"LemmyBeProxy/helper"
	"LemmyBeProxy/helper/converter"
	appHttp "LemmyBeProxy/http"
	"LemmyBeProxy/service"
)

// PiefedBackend implements the backend.Backend interface, translating
// between Lemmy-shaped requests/responses (this proxy's canonical shape)
// and Piefed's real API shape. It wraps the low-level *Piefed client,
// which just makes the raw HTTP calls. activityPub and simulateLemmy
// exist solely for Site — Piefed's own /site response needs supplementing
// with a direct ActivityPub actor fetch for fields it doesn't expose
// natively, a genuinely Piefed-specific need no other endpoint has.
type PiefedBackend struct {
	client        *Piefed
	activityPub   *service.ActivityPub
	simulateLemmy bool
}

func NewPiefedBackend(client *Piefed, activityPub *service.ActivityPub, simulateLemmy bool) *PiefedBackend {
	return &PiefedBackend{client: client, activityPub: activityPub, simulateLemmy: simulateLemmy}
}

// convertGetPostResponse is shared by every Post method that returns
// Piefed's GetPostResponse shape (get, create, edit, like all do).
func convertGetPostResponse(resp *piefedResponse.GetPostResponse) *lemmyResponse.GetPostResponse {
	return &lemmyResponse.GetPostResponse{
		CommunityView: converter.ConvertCommunityView(resp.CommunityView),
		CrossPosts:    helper.MapSlice(resp.CrossPosts, converter.ConvertPostView),
		Moderators:    helper.MapSlice(resp.Moderators, converter.ConvertCommunityModeratorView),
		PostView:      converter.ConvertPostView(resp.PostView),
	}
}

func (receiver *PiefedBackend) GetPosts(request *lemmyRequest.GetPostsRequest, headers appHttp.Headers) (*lemmyResponse.GetPostsResponse, error) {
	// Piefed's /post/list ignores community_id/community_name entirely
	// when type_ is also set to Subscribed (confirmed directly: sending
	// both returns the full subscribed feed, not the requested
	// community's posts — a real Piefed backend quirk). A
	// community-scoped request is unambiguous on its own, so the listing
	// type is dropped whenever either community field is present rather
	// than risking Piefed silently discarding the more specific filter.
	isCommunityScoped := request.CommunityId != nil || request.CommunityName != nil

	resp, err := receiver.client.GetPosts(&piefedRequest.GetPostsRequest{
		Type: helper.SafeDereference(request.Type, func(in lemmyModel.ListingType) *piefedModel.ListingType {
			if isCommunityScoped {
				return nil
			}
			return helper.ToPointer(converter.ReverseConvertListingType(in))
		}),
		Sort: helper.SafeDereference(request.Sort, func(in lemmyModel.SortType) *piefedModel.SortType {
			return helper.ToPointer(converter.ReverseConvertSortType(in))
		}),
		Page:          request.Page,
		Limit:         request.Limit,
		CommunityId:   request.CommunityId,
		PersonId:      nil,
		CommunityName: request.CommunityName,
		LikedOnly:     request.LikedOnly,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetPostsResponse{
		NextPage: resp.NextPage,
		Posts:    helper.MapSlice(resp.Posts, converter.ConvertPostView),
	}, nil
}

func (receiver *PiefedBackend) GetPost(request *lemmyRequest.GetPostRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	resp, err := receiver.client.GetPost(&piefedRequest.GetPostRequest{
		CommentId: request.CommentId,
		Id:        request.Id,
	}, headers)
	if err != nil {
		return nil, err
	}

	return convertGetPostResponse(resp), nil
}

// CreatePost maps Lemmy's Name field to Piefed's Title — same content,
// different field name. Honeypot has no Piefed equivalent and is dropped.
func (receiver *PiefedBackend) CreatePost(request *lemmyRequest.CreatePostRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	resp, err := receiver.client.CreatePost(&piefedRequest.CreatePostRequest{
		Title:       request.Name,
		CommunityId: request.CommunityId,
		Body:        request.Body,
		Url:         request.Url,
		Nsfw:        request.Nsfw,
		LanguageId:  request.LanguageId,
	}, headers)
	if err != nil {
		return nil, err
	}

	return convertGetPostResponse(resp), nil
}

func (receiver *PiefedBackend) EditPost(request *lemmyRequest.EditPostRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	resp, err := receiver.client.EditPost(&piefedRequest.EditPostRequest{
		PostId:     request.PostId,
		Title:      request.Name,
		Body:       request.Body,
		Url:        request.Url,
		Nsfw:       request.Nsfw,
		LanguageId: request.LanguageId,
	}, headers)
	if err != nil {
		return nil, err
	}

	return convertGetPostResponse(resp), nil
}

func (receiver *PiefedBackend) LikePost(request *lemmyRequest.CreatePostLikeRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	resp, err := receiver.client.LikePost(&piefedRequest.LikePostRequest{
		PostId: request.PostId,
		Score:  request.Score,
	}, headers)
	if err != nil {
		return nil, err
	}

	return convertGetPostResponse(resp), nil
}

func (receiver *PiefedBackend) MarkPostAsRead(request *lemmyRequest.MarkPostAsReadRequest, headers appHttp.Headers) (*lemmyResponse.SuccessResponse, error) {
	resp, err := receiver.client.MarkPostAsRead(&piefedRequest.MarkPostAsReadRequest{
		PostId:  request.PostId,
		PostIds: request.PostIds,
		Read:    request.Read,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.SuccessResponse{Success: resp.Success}, nil
}

func (receiver *PiefedBackend) GetComments(request *lemmyRequest.GetCommentsRequest, headers appHttp.Headers) (*lemmyResponse.GetCommentsResponse, error) {
	resp, err := receiver.client.GetComments(&piefedRequest.GetCommentsRequest{
		Type: helper.SafeDereference(request.Type, func(in lemmyModel.ListingType) *piefedModel.ListingType {
			return helper.ToPointer(converter.ReverseConvertListingType(in))
		}),
		PersonId:    nil,
		MaxDepth:    request.MaxDepth,
		Page:        request.Page,
		ParentId:    request.ParentId,
		CommunityId: request.CommunityId,
		PostId:      request.PostId,
		Limit:       request.Limit,
		Sort: helper.SafeDereference(request.Sort, func(in lemmyModel.CommentSortType) *piefedModel.CommentSortType {
			return helper.ToPointer(converter.ReverseConvertCommentSortType(in))
		}),
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetCommentsResponse{
		Comments: helper.MapSlice(resp.Comments, converter.ConvertCommentView),
	}, nil
}

func (receiver *PiefedBackend) GetComment(request *lemmyRequest.GetCommentRequest, headers appHttp.Headers) (*lemmyResponse.GetCommentResponse, error) {
	resp, err := receiver.client.GetComment(&piefedRequest.GetCommentRequest{
		Id: request.Id,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetCommentResponse{
		CommentView:  converter.ConvertCommentView(resp.CommentView),
		RecipientIds: []uint{},
	}, nil
}

func (receiver *PiefedBackend) CreateComment(request *lemmyRequest.CreateCommentRequest, headers appHttp.Headers) (*lemmyResponse.CreateCommentResponse, error) {
	resp, err := receiver.client.CreateComment(&piefedRequest.CreateCommentRequest{
		Body:       request.Content,
		PostId:     request.PostId,
		ParentId:   request.ParentId,
		LanguageId: request.LanguageId,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.CreateCommentResponse{
		CommentView:  converter.ConvertCommentView(resp.CommentView),
		RecipientIds: []uint{},
	}, nil
}

func (receiver *PiefedBackend) LikeComment(request *lemmyRequest.CreateCommentLikeRequest, headers appHttp.Headers) (*lemmyResponse.GetCommentResponse, error) {
	resp, err := receiver.client.LikeComment(&piefedRequest.LikeCommentRequest{
		CommentId: request.CommentId,
		Score:     request.Score,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetCommentResponse{
		CommentView:  converter.ConvertCommentView(resp.CommentView),
		RecipientIds: []uint{},
	}, nil
}

func (receiver *PiefedBackend) GetCommunity(request *lemmyRequest.GetCommunityRequest, headers appHttp.Headers) (*lemmyResponse.GetCommunityResponse, error) {
	resp, err := receiver.client.GetCommunity(&piefedRequest.GetCommunityRequest{
		Id:   request.Id,
		Name: request.Name,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetCommunityResponse{
		CommunityView: converter.ConvertCommunityView(resp.CommunityView),
		Moderators:    helper.MapSlice(resp.Moderators, converter.ConvertCommunityModeratorView),
	}, nil
}

func (receiver *PiefedBackend) GetCommunities(request *lemmyRequest.GetCommunitiesRequest, headers appHttp.Headers) (*lemmyResponse.GetCommunitiesResponse, error) {
	resp, err := receiver.client.GetCommunities(&piefedRequest.GetCommunitiesRequest{
		Type: helper.SafeDereference(request.Type, func(in lemmyModel.ListingType) *piefedModel.ListingType {
			return helper.ToPointer(converter.ReverseConvertListingType(in))
		}),
		Sort: helper.SafeDereference(request.Sort, func(in lemmyModel.SortType) *piefedModel.SortType {
			return helper.ToPointer(converter.ReverseConvertCommunitySortType(in))
		}),
		ShowNsfw: request.ShowNsfw,
		Page:     request.Page,
		Limit:    request.Limit,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetCommunitiesResponse{
		Communities: helper.MapSlice(resp.Communities, converter.ConvertCommunityView),
	}, nil
}

func (receiver *PiefedBackend) FollowCommunity(request *lemmyRequest.FollowCommunityRequest, headers appHttp.Headers) (*lemmyResponse.CommunityResponse, error) {
	resp, err := receiver.client.FollowCommunity(&piefedRequest.FollowCommunityRequest{
		CommunityId: request.CommunityId,
		Follow:      request.Follow,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.CommunityResponse{
		CommunityView:       converter.ConvertCommunityView(resp.CommunityView),
		DiscussionLanguages: resp.DiscussionLanguages,
	}, nil
}

func (receiver *PiefedBackend) BlockCommunity(request *lemmyRequest.BlockCommunityRequest, headers appHttp.Headers) (*lemmyResponse.BlockCommunityResponse, error) {
	resp, err := receiver.client.BlockCommunity(&piefedRequest.BlockCommunityRequest{
		CommunityId: request.CommunityId,
		Block:       request.Block,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.BlockCommunityResponse{
		Blocked:       resp.Blocked,
		CommunityView: converter.ConvertCommunityView(resp.CommunityView),
	}, nil
}

func (receiver *PiefedBackend) Login(request *lemmyRequest.LoginRequest, headers appHttp.Headers) (*lemmyResponse.LoginResponse, error) {
	resp, err := receiver.client.Login(&piefedRequest.LoginRequest{
		Username: request.UsernameOrEmail,
		Password: request.Password,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.LoginResponse{
		Jwt: helper.ToPointer(resp.Jwt),
	}, nil
}

func (receiver *PiefedBackend) GetUnreadCount(headers appHttp.Headers) (*lemmyResponse.GetUnreadCountResponse, error) {
	resp, err := receiver.client.GetUnreadCount(headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetUnreadCountResponse{
		Mentions:        resp.Mentions,
		PrivateMessages: resp.PrivateMessages,
		Replies:         resp.Replies,
	}, nil
}

func (receiver *PiefedBackend) GetUser(request *lemmyRequest.GetUserRequest, headers appHttp.Headers) (*lemmyResponse.GetUserResponse, error) {
	resp, err := receiver.client.GetUser(&piefedRequest.GetUserRequest{
		PersonId: request.PersonId,
		Username: request.Username,
		Sort: helper.SafeDereference(request.Sort, func(in lemmyModel.SortType) *piefedModel.SortType {
			return helper.ToPointer(converter.ReverseConvertSortType(in))
		}),
		Page:      request.Page,
		Limit:     request.Limit,
		SavedOnly: request.SavedOnly,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.GetUserResponse{
		Comments:   helper.MapSlice(resp.Comments, converter.ConvertCommentView),
		Moderates:  helper.MapSlice(resp.Moderates, converter.ConvertCommunityModeratorView),
		PersonView: converter.ConvertPersonView(resp.PersonView),
		Posts:      helper.MapSlice(resp.Posts, converter.ConvertPostView),
	}, nil
}

func (receiver *PiefedBackend) BlockPerson(request *lemmyRequest.BlockPersonRequest, headers appHttp.Headers) (*lemmyResponse.BlockPersonResponse, error) {
	resp, err := receiver.client.BlockPerson(&piefedRequest.BlockPersonRequest{
		PersonId: request.PersonId,
		Block:    request.Block,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.BlockPersonResponse{
		Blocked:    resp.Blocked,
		PersonView: converter.ConvertPersonView(resp.PersonView),
	}, nil
}

// SaveUserSettings only forwards the four fields Piefed's own
// save_user_settings endpoint actually supports (ShowNsfw,
// DefaultSortType, DefaultCommentSortType, ShowReadPosts). Everything
// else is accepted (so the save doesn't fail outright) but has no
// Piefed field to go to and is silently dropped.
func (receiver *PiefedBackend) SaveUserSettings(request *lemmyRequest.SaveUserSettingsRequest, headers appHttp.Headers) (*lemmyResponse.SaveUserSettingsResponse, error) {
	_, err := receiver.client.SaveUserSettings(&piefedRequest.SaveUserSettingsRequest{
		ShowNsfw: request.ShowNsfw,
		DefaultSortType: helper.SafeDereference(request.DefaultSortType, func(in string) *string {
			return helper.ToPointer(converter.ClampDefaultSortType(in))
		}),
		DefaultCommentSortType: helper.SafeDereference(request.DefaultCommentSortType, func(in string) *string {
			return helper.ToPointer(converter.ClampDefaultCommentSortType(in))
		}),
		ShowReadPosts: request.ShowReadPosts,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.SaveUserSettingsResponse{}, nil
}

func (receiver *PiefedBackend) Search(request *lemmyRequest.SearchRequest, headers appHttp.Headers) (*lemmyResponse.SearchResponse, error) {
	resp, err := receiver.client.Search(&piefedRequest.SearchRequest{
		Q:    request.Q,
		Type: converter.ReverseConvertSearchType(request.Type),
		ListingType: helper.SafeDereference(request.ListingType, func(in lemmyModel.ListingType) *piefedModel.ListingType {
			return helper.ToPointer(converter.ReverseConvertListingType(in))
		}),
		Sort: helper.SafeDereference(request.Sort, func(in lemmyModel.SortType) *piefedModel.SortType {
			return helper.ToPointer(converter.ReverseConvertSortType(in))
		}),
		Page:          request.Page,
		Limit:         request.Limit,
		CommunityName: request.CommunityName,
		CommunityId:   request.CommunityId,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.SearchResponse{
		Type:        converter.ConvertSearchType(resp.Type),
		Communities: helper.MapSlice(resp.Communities, converter.ConvertCommunityView),
		Posts:       helper.MapSlice(resp.Posts, converter.ConvertPostView),
		Users:       helper.MapSlice(resp.Users, converter.ConvertPersonView),
		Comments:    helper.MapSlice(resp.Comments, converter.ConvertCommentView),
	}, nil
}

// Site does real, Piefed-specific work no other endpoint needs: an
// ActivityPub actor fetch to supplement fields Piefed's own /site
// response doesn't expose natively (moved here unchanged from the old
// SiteController implementation).
func (receiver *PiefedBackend) Site(headers appHttp.Headers) (*lemmyResponse.GetSiteResponse, error) {
	resp, err := receiver.client.Site(headers)
	if err != nil {
		return nil, err
	}

	apActor, err := receiver.activityPub.FetchActor(resp.Site.ActorId)
	if err != nil {
		return nil, err
	}

	var version string
	if receiver.simulateLemmy {
		version = "0.19.11"
	} else {
		version = resp.Version
	}

	return &lemmyResponse.GetSiteResponse{
		Admins:       helper.MapSlice(resp.Admins, converter.ConvertPersonView),
		AllLanguages: helper.MapSlice(resp.Site.AllLanguages, converter.ConvertLanguageView),
		BlockedUrls:  []lemmyModel.LocalSiteUrlBlocklist{},
		CustomEmojis: []lemmyModel.CustomEmojiView{},
		DiscussionLanguages: helper.MapSlice(resp.Site.AllLanguages, func(in piefedModel.LanguageView) uint {
			return in.Id
		}),
		MyUser:   converter.ConvertMyUserInfo(resp.MyUser, apActor),
		SiteView: converter.ConvertSiteToView(&resp.Site, apActor),
		Taglines: []lemmyModel.Tagline{},
		Version:  version,
	}, nil
}

// UploadImage forwards to Piefed's own /upload/image endpoint, which
// expects Bearer auth (unlike real Lemmy's pict-rs, which expects a
// Cookie — see LemmyBackend's implementation for the difference).
func (receiver *PiefedBackend) UploadImage(fileBytes []byte, filename string, jwt string) (string, error) {
	resp, err := receiver.client.UploadImage(fileBytes, filename, jwt)
	if err != nil {
		return "", err
	}

	return resp.Url, nil
}

// GetPersonMentions, GetReplies, and GetPrivateMessages all return an
// empty (not missing, not erroring) list against Piefed. Piefed exposes
// unread counts for these (see GetUnreadCount above) but has no
// equivalent endpoint to fetch the actual list of mentions, comment
// replies, or private messages we've found — person mentions, comment
// replies, and private messages aren't native Piefed concepts in the
// same shape Lemmy has them. Returning an empty list rather than a 404
// matters concretely: a 404 here is exactly what crashes a real Lemmy
// client (lemmyBB included), since these are fetched on every
// authenticated page load. This is a real, known gap — Piefed users
// won't see their actual mentions/replies/messages through this proxy —
// not a claim that Piefed's own notifications are being surfaced.
func (receiver *PiefedBackend) GetPersonMentions(request *lemmyRequest.GetPersonMentionsRequest, headers appHttp.Headers) (*lemmyResponse.GetPersonMentionsResponse, error) {
	return &lemmyResponse.GetPersonMentionsResponse{Mentions: []lemmyModel.PersonMentionView{}}, nil
}

func (receiver *PiefedBackend) GetReplies(request *lemmyRequest.GetRepliesRequest, headers appHttp.Headers) (*lemmyResponse.GetRepliesResponse, error) {
	return &lemmyResponse.GetRepliesResponse{Replies: []lemmyModel.CommentReplyView{}}, nil
}

func (receiver *PiefedBackend) GetPrivateMessages(request *lemmyRequest.GetPrivateMessagesRequest, headers appHttp.Headers) (*lemmyResponse.GetPrivateMessagesResponse, error) {
	return &lemmyResponse.GetPrivateMessagesResponse{PrivateMessages: []lemmyModel.PrivateMessageView{}}, nil
}

// ResolveObject forwards to Piefed's own /resolve_object endpoint — see
// dto/request/piefed/ResolveObjectRequest.go for a note on how that
// path was determined (Piefed's release notes confirm the feature
// exists; the exact path is a well-justified assumption based on every
// other endpoint's established naming convention, not directly
// source-verified).
func (receiver *PiefedBackend) ResolveObject(request *lemmyRequest.ResolveObjectRequest, headers appHttp.Headers) (*lemmyResponse.ResolveObjectResponse, error) {
	resp, err := receiver.client.ResolveObject(&piefedRequest.ResolveObjectRequest{
		Q: request.Q,
	}, headers)
	if err != nil {
		return nil, err
	}

	canonical := &lemmyResponse.ResolveObjectResponse{}
	if resp.Comment != nil {
		canonical.Comment = helper.ToPointer(converter.ConvertCommentView(*resp.Comment))
	}
	if resp.Post != nil {
		canonical.Post = helper.ToPointer(converter.ConvertPostView(*resp.Post))
	}
	if resp.Community != nil {
		canonical.Community = helper.ToPointer(converter.ConvertCommunityView(*resp.Community))
	}
	if resp.Person != nil {
		canonical.Person = helper.ToPointer(converter.ConvertPersonView(*resp.Person))
	}

	return canonical, nil
}

func (receiver *PiefedBackend) GetCaptcha(headers appHttp.Headers) (*lemmyResponse.GetCaptchaResponse, error) {
	return &lemmyResponse.GetCaptchaResponse{Ok: nil}, nil
}

// Register forwards to Piefed's own /user/register endpoint — see
// PiefedUser.go's Register and dto/request/piefed/RegisterRequest.go
// for a note on that path being an unverified assumption, not directly
// source-verified against Piefed.
func (receiver *PiefedBackend) Register(request *lemmyRequest.RegisterRequest, headers appHttp.Headers) (*lemmyResponse.LoginResponse, error) {
	resp, err := receiver.client.Register(&piefedRequest.RegisterRequest{
		Username:       request.Username,
		Password:       request.Password,
		PasswordVerify: request.PasswordVerify,
		ShowNsfw:       request.ShowNsfw,
		Email:          request.Email,
		Answer:         request.ApplicationQuestionAnswer,
	}, headers)
	if err != nil {
		return nil, err
	}

	return &lemmyResponse.LoginResponse{
		Jwt: resp.Jwt,
		// A successful call means a registration was created — separate
		// from whether Piefed grants an immediate JWT (nil here means
		// the registration is pending admin approval, matching real
		// Lemmy's own semantics for requireapplication mode).
		RegistrationCreated: true,
		VerifyEmailSent:     false,
	}, nil
}

// MarkAllAsRead has no Piefed equivalent — Piefed's own unread-count
// tracking has no "mark everything read in one call" action to map
// this to. Returning an empty GetRepliesResponse rather than a 404
// matters the same way it does for GetPersonMentions/GetReplies: a 404
// here is what would break the client, not what's actually missing
// from the user's perspective (there was nothing to mark as read
// against Piefed in the first place, since GetPersonMentions/GetReplies
// already return empty for the same reason).
func (receiver *PiefedBackend) MarkAllAsRead(headers appHttp.Headers) (*lemmyResponse.GetRepliesResponse, error) {
	return &lemmyResponse.GetRepliesResponse{Replies: []lemmyModel.CommentReplyView{}}, nil
}
