package frontend

import (
	lemmyRequest "LemmyBeProxy/dto/request/lemmy"
	lemmyRequest017 "LemmyBeProxy/dto/request/lemmy017"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	lemmyResponse017 "LemmyBeProxy/dto/response/lemmy017"
	"LemmyBeProxy/helper"
	"LemmyBeProxy/helper/converter"
	"LemmyBeProxy/http"
)

// Frontend017 implements the Frontend interface for Lemmy 0.17.x clients
// (built for lemmyBB, which is pinned to lemmy_api_common 0.17.2). Every
// shape difference here is confirmed against Lemmy's own source at tag
// 0.17.2, not assumed — see dto/request/lemmy017 and
// dto/response/lemmy017 for the specific differences each type documents.
type Frontend017 struct{}

func NewFrontend017() *Frontend017 {
	return &Frontend017{}
}

func (receiver *Frontend017) ParseGetPostsRequest(request *http.Request) (*lemmyRequest.GetPostsRequest, error) {
	reqDto, err := helper.ParseRequestQuery[lemmyRequest017.GetPostsRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.GetPostsRequest{
		Type:          reqDto.Type,
		Sort:          reqDto.Sort,
		Page:          reqDto.Page,
		Limit:         reqDto.Limit,
		CommunityId:   reqDto.CommunityId,
		CommunityName: reqDto.CommunityName,
		SavedOnly:     reqDto.SavedOnly,
		// LikedOnly, DislikedOnly, PageCursor, ShowHidden, ShowNsfw,
		// ShowRead don't exist in 0.17.x — left unset.
	}, nil
}

func (receiver *Frontend017) BuildGetPostsResponse(resp *lemmyResponse.GetPostsResponse) any {
	// 0.17.x has no next_page concept at all — dropped, not translated.
	return &lemmyResponse017.GetPostsResponse{
		Posts: helper.MapSlice(resp.Posts, converter.ConvertPostViewTo017),
	}
}

func (receiver *Frontend017) ParseGetPostRequest(request *http.Request) (*lemmyRequest.GetPostRequest, error) {
	reqDto, err := helper.ParseRequestQuery[lemmyRequest017.GetPostRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.GetPostRequest{
		Id:        reqDto.Id,
		CommentId: reqDto.CommentId,
	}, nil
}

func (receiver *Frontend017) BuildGetPostResponse(resp *lemmyResponse.GetPostResponse) any {
	// 0.17.x has no cross_posts field and expects an "online" viewer
	// count our canonical model doesn't track at all — 0 is the honest
	// Online is 0, an honest not-guessed value here rather than a guess
	// standing in for real data.
	return &lemmyResponse017.GetPostResponse{
		PostView:      converter.ConvertPostViewTo017(resp.PostView),
		CommunityView: converter.ConvertCommunityViewTo017(resp.CommunityView),
		Moderators:    helper.MapSlice(resp.Moderators, converter.ConvertCommunityModeratorViewTo017),
		Online:        0,
	}
}

func (receiver *Frontend017) ParseCreatePostRequest(request *http.Request) (*lemmyRequest.CreatePostRequest, error) {
	reqDto, err := helper.ParseRequest[lemmyRequest017.CreatePostRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.CreatePostRequest{
		Name:        reqDto.Name,
		CommunityId: reqDto.CommunityId,
		Body:        reqDto.Body,
		Url:         reqDto.Url,
		Nsfw:        reqDto.Nsfw,
		LanguageId:  reqDto.LanguageId,
		Honeypot:    reqDto.Honeypot,
	}, nil
}

func (receiver *Frontend017) ParseEditPostRequest(request *http.Request) (*lemmyRequest.EditPostRequest, error) {
	reqDto, err := helper.ParseRequest[lemmyRequest017.EditPostRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.EditPostRequest{
		PostId:     reqDto.PostId,
		Name:       reqDto.Name,
		Body:       reqDto.Body,
		Url:        reqDto.Url,
		Nsfw:       reqDto.Nsfw,
		LanguageId: reqDto.LanguageId,
	}, nil
}

func (receiver *Frontend017) ParseCreatePostLikeRequest(request *http.Request) (*lemmyRequest.CreatePostLikeRequest, error) {
	reqDto, err := helper.ParseRequest[lemmyRequest017.CreatePostLikeRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.CreatePostLikeRequest{
		PostId: reqDto.PostId,
		Score:  reqDto.Score,
	}, nil
}

func (receiver *Frontend017) BuildPostMutationResponse(resp *lemmyResponse.GetPostResponse) any {
	// 0.17.x's create/edit/like all return the lean PostResponse{post_view}
	// shape — confirmed against each handler in Lemmy's own source, not
	// the fuller GetPostResponse used for fetching.
	return &lemmyResponse017.PostResponse{
		PostView: converter.ConvertPostViewTo017(resp.PostView),
	}
}

func (receiver *Frontend017) ParseMarkPostAsReadRequest(request *http.Request) (*lemmyRequest.MarkPostAsReadRequest, error) {
	reqDto, err := helper.ParseRequest[lemmyRequest017.MarkPostAsReadRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.MarkPostAsReadRequest{
		PostId: helper.ToPointer(reqDto.PostId),
		Read:   reqDto.Read,
	}, nil
}

// BuildSuccessResponse is the one known gap in this slice — see the
// comment on the Frontend interface for why. Real Lemmy 0.17.x expects
// PostResponse{post_view} here; this returns the canonical
// {success: bool} shape instead, since building an accurate post_view
// would need an extra backend round-trip that's a real design decision,
// not something to bolt on silently.
func (receiver *Frontend017) BuildSuccessResponse(resp *lemmyResponse.SuccessResponse) any {
	return resp
}

func (receiver *Frontend017) ParseGetCommentsRequest(request *http.Request) (*lemmyRequest.GetCommentsRequest, error) {
	reqDto, err := helper.ParseRequestQuery[lemmyRequest017.GetCommentsRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.GetCommentsRequest{
		Type:          reqDto.Type,
		Sort:          reqDto.Sort,
		MaxDepth:      reqDto.MaxDepth,
		Page:          reqDto.Page,
		Limit:         reqDto.Limit,
		CommunityId:   reqDto.CommunityId,
		CommunityName: reqDto.CommunityName,
		PostId:        reqDto.PostId,
		ParentId:      reqDto.ParentId,
		SavedOnly:     reqDto.SavedOnly,
		// LikedOnly, DislikedOnly don't exist in 0.17.x — left unset.
	}, nil
}

func (receiver *Frontend017) BuildGetCommentsResponse(resp *lemmyResponse.GetCommentsResponse) any {
	return &lemmyResponse017.GetCommentsResponse{
		Comments: helper.MapSlice(resp.Comments, converter.ConvertCommentViewTo017),
	}
}

func (receiver *Frontend017) ParseGetCommentRequest(request *http.Request) (*lemmyRequest.GetCommentRequest, error) {
	reqDto, err := helper.ParseRequestQuery[lemmyRequest017.GetCommentRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.GetCommentRequest{
		Id: reqDto.Id,
	}, nil
}

func (receiver *Frontend017) BuildGetCommentResponse(resp *lemmyResponse.GetCommentResponse) any {
	return &lemmyResponse017.CommentResponse{
		CommentView:  converter.ConvertCommentViewTo017(resp.CommentView),
		RecipientIds: resp.RecipientIds,
		// FormId: known gap, see CommentResponse017's own comment.
	}
}

func (receiver *Frontend017) ParseCreateCommentRequest(request *http.Request) (*lemmyRequest.CreateCommentRequest, error) {
	reqDto, err := helper.ParseRequest[lemmyRequest017.CreateCommentRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.CreateCommentRequest{
		Content:    reqDto.Content,
		PostId:     reqDto.PostId,
		ParentId:   reqDto.ParentId,
		LanguageId: reqDto.LanguageId,
	}, nil
}

func (receiver *Frontend017) BuildCreateCommentResponse(resp *lemmyResponse.CreateCommentResponse) any {
	return &lemmyResponse017.CommentResponse{
		CommentView:  converter.ConvertCommentViewTo017(resp.CommentView),
		RecipientIds: resp.RecipientIds,
	}
}

func (receiver *Frontend017) ParseCreateCommentLikeRequest(request *http.Request) (*lemmyRequest.CreateCommentLikeRequest, error) {
	reqDto, err := helper.ParseRequest[lemmyRequest017.CreateCommentLikeRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.CreateCommentLikeRequest{
		CommentId: reqDto.CommentId,
		Score:     reqDto.Score,
	}, nil
}

func (receiver *Frontend017) BuildLikeCommentResponse(resp *lemmyResponse.GetCommentResponse) any {
	return &lemmyResponse017.CommentResponse{
		CommentView:  converter.ConvertCommentViewTo017(resp.CommentView),
		RecipientIds: resp.RecipientIds,
	}
}

// ParseGetCommunityRequest reuses the canonical parse directly — 0.17.2's
// real GetCommunity (id, name, auth) matches our canonical model exactly.
func (receiver *Frontend017) ParseGetCommunityRequest(request *http.Request) (*lemmyRequest.GetCommunityRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetCommunityRequest](request)
}

func (receiver *Frontend017) BuildGetCommunityResponse(resp *lemmyResponse.GetCommunityResponse) any {
	return &lemmyResponse017.GetCommunityResponse{
		CommunityView:       converter.ConvertCommunityViewTo017(resp.CommunityView),
		Site:                nil,
		Moderators:          helper.MapSlice(resp.Moderators, converter.ConvertCommunityModeratorViewTo017),
		Online:              0,
		DiscussionLanguages: []uint{},
		DefaultPostLanguage: nil,
	}
}

func (receiver *Frontend017) ParseGetCommunitiesRequest(request *http.Request) (*lemmyRequest.GetCommunitiesRequest, error) {
	reqDto, err := helper.ParseRequestQuery[lemmyRequest017.GetCommunitiesRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.GetCommunitiesRequest{
		Type:  reqDto.Type,
		Sort:  reqDto.Sort,
		Page:  reqDto.Page,
		Limit: reqDto.Limit,
		// ShowNsfw doesn't exist in 0.17.x — left unset.
	}, nil
}

func (receiver *Frontend017) BuildGetCommunitiesResponse(resp *lemmyResponse.GetCommunitiesResponse) any {
	return &lemmyResponse017.GetCommunitiesResponse{
		Communities: helper.MapSlice(resp.Communities, converter.ConvertCommunityViewTo017),
	}
}

// ParseFollowCommunityRequest reuses the canonical parse directly —
// 0.17.2's real FollowCommunity (community_id, follow, auth) matches
// exactly.
func (receiver *Frontend017) ParseFollowCommunityRequest(request *http.Request) (*lemmyRequest.FollowCommunityRequest, error) {
	return helper.ParseRequest[lemmyRequest.FollowCommunityRequest](request)
}

func (receiver *Frontend017) BuildCommunityResponse(resp *lemmyResponse.CommunityResponse) any {
	return &lemmyResponse017.CommunityResponse{
		CommunityView:       converter.ConvertCommunityViewTo017(resp.CommunityView),
		DiscussionLanguages: resp.DiscussionLanguages,
	}
}

// ParseBlockCommunityRequest reuses the canonical parse directly —
// 0.17.2's real BlockCommunity (community_id, block, auth) matches
// exactly.
func (receiver *Frontend017) ParseBlockCommunityRequest(request *http.Request) (*lemmyRequest.BlockCommunityRequest, error) {
	return helper.ParseRequest[lemmyRequest.BlockCommunityRequest](request)
}

func (receiver *Frontend017) BuildBlockCommunityResponse(resp *lemmyResponse.BlockCommunityResponse) any {
	return &lemmyResponse017.BlockCommunityResponse{
		Blocked:       resp.Blocked,
		CommunityView: converter.ConvertCommunityViewTo017(resp.CommunityView),
	}
}

// ParseLoginRequest reuses the canonical parse directly — 0.17.2's real
// Login (username_or_email, password) matches exactly.
func (receiver *Frontend017) ParseLoginRequest(request *http.Request) (*lemmyRequest.LoginRequest, error) {
	return helper.ParseRequest[lemmyRequest.LoginRequest](request)
}

func (receiver *Frontend017) BuildLoginResponse(resp *lemmyResponse.LoginResponse) any {
	return resp
}

func (receiver *Frontend017) BuildGetUnreadCountResponse(resp *lemmyResponse.GetUnreadCountResponse) any {
	return resp
}

// ParseGetUserRequest reuses the canonical parse directly — 0.17.2's real
// GetPersonDetails matches our canonical GetUserRequest's fields exactly
// (community_id isn't wired through on either side currently).
func (receiver *Frontend017) ParseGetUserRequest(request *http.Request) (*lemmyRequest.GetUserRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetUserRequest](request)
}

func (receiver *Frontend017) BuildGetUserResponse(resp *lemmyResponse.GetUserResponse) any {
	return &lemmyResponse017.GetUserResponse{
		Comments:   helper.MapSlice(resp.Comments, converter.ConvertCommentViewTo017),
		Moderates:  helper.MapSlice(resp.Moderates, converter.ConvertCommunityModeratorViewTo017),
		PersonView: converter.ConvertPersonViewTo017(resp.PersonView),
		Posts:      helper.MapSlice(resp.Posts, converter.ConvertPostViewTo017),
	}
}

// ParseBlockPersonRequest reuses the canonical parse directly — 0.17.2's
// real BlockPerson (person_id, block, auth) matches exactly.
func (receiver *Frontend017) ParseBlockPersonRequest(request *http.Request) (*lemmyRequest.BlockPersonRequest, error) {
	return helper.ParseRequest[lemmyRequest.BlockPersonRequest](request)
}

func (receiver *Frontend017) BuildBlockPersonResponse(resp *lemmyResponse.BlockPersonResponse) any {
	return &lemmyResponse017.BlockPersonResponse{
		Blocked:    resp.Blocked,
		PersonView: converter.ConvertPersonViewTo017(resp.PersonView),
	}
}

// ParseSearchRequest reuses the canonical parse directly — 0.17.2's
// real Search matches our canonical model's fields exactly (creator_id
// isn't supported on either side, a pre-existing Piefed limitation, not
// something specific to 0.17.x).
func (receiver *Frontend017) ParseSearchRequest(request *http.Request) (*lemmyRequest.SearchRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.SearchRequest](request)
}

func (receiver *Frontend017) BuildSearchResponse(resp *lemmyResponse.SearchResponse) any {
	return &lemmyResponse017.SearchResponse{
		Type:        string(resp.Type),
		Comments:    helper.MapSlice(resp.Comments, converter.ConvertCommentViewTo017),
		Posts:       helper.MapSlice(resp.Posts, converter.ConvertPostViewTo017),
		Communities: helper.MapSlice(resp.Communities, converter.ConvertCommunityViewTo017),
		Users:       helper.MapSlice(resp.Users, converter.ConvertPersonViewTo017),
	}
}

// BuildGetSiteResponse now covers MyUser too, using the full MyUserInfo
// translation built for this — see the interface's own comment for the
// remaining known gap (FederatedInstances).
func (receiver *Frontend017) BuildGetSiteResponse(resp *lemmyResponse.GetSiteResponse) any {
	return &lemmyResponse017.GetSiteResponse{
		SiteView:            converter.ConvertSiteViewTo017(resp.SiteView),
		Admins:              helper.MapSlice(resp.Admins, converter.ConvertPersonViewTo017),
		Online:              0,
		Version:             resp.Version,
		MyUser:              converter.ConvertMyUserInfoTo017(resp.MyUser),
		AllLanguages:        resp.AllLanguages,
		DiscussionLanguages: resp.DiscussionLanguages,
		Taglines:            resp.Taglines,
		FederatedInstances:  nil,
	}
}

// ParseSaveUserSettingsRequest converts 0.17.2's real, simpler field set
// (no default_comment_sort_type at all — a 0.17.x client genuinely can't
// set that preference) into canonical's shape. DefaultSortType/
// DefaultListingType use the same numeric-index mapping built for
// MyUserInfo. Fields with no canonical equivalent (avatar, banner,
// display_name, email, bio, matrix_user_id, bot_account,
// show_new_post_notifs) are simply not forwarded — canonical has no
// field for them, same "accepted but nothing to forward to" pattern as
// everywhere else in this endpoint.
func (receiver *Frontend017) ParseSaveUserSettingsRequest(request *http.Request) (*lemmyRequest.SaveUserSettingsRequest, error) {
	reqDto, err := helper.ParseRequest[lemmyRequest017.SaveUserSettingsRequest](request)
	if err != nil {
		return nil, err
	}

	return &lemmyRequest.SaveUserSettingsRequest{
		ShowNsfw: reqDto.ShowNsfw,
		DefaultSortType: helper.SafeDereference(reqDto.DefaultSortType, func(in int16) *string {
			return helper.ToPointer(string(converter.ConvertIndex017ToSortType(in)))
		}),
		DefaultListingType: helper.SafeDereference(reqDto.DefaultListingType, func(in int16) *string {
			return helper.ToPointer(string(converter.ConvertIndex017ToListingType(in)))
		}),
		InterfaceLanguage:        reqDto.InterfaceLanguage,
		ShowAvatars:              reqDto.ShowAvatars,
		ShowScores:               reqDto.ShowScores,
		SendNotificationsToEmail: reqDto.SendNotificationsToEmail,
		ShowBotAccounts:          reqDto.ShowBotAccounts,
		ShowReadPosts:            reqDto.ShowReadPosts,
		DiscussionLanguages:      reqDto.DiscussionLanguages,
	}, nil
}

func (receiver *Frontend017) BuildSaveUserSettingsResponse(resp *lemmyResponse.SaveUserSettingsResponse) any {
	return resp
}

// ParseGetPersonMentionsRequest reuses the canonical parse directly —
// 0.17.2's real GetPersonMentions matches exactly.
func (receiver *Frontend017) ParseGetPersonMentionsRequest(request *http.Request) (*lemmyRequest.GetPersonMentionsRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetPersonMentionsRequest](request)
}

func (receiver *Frontend017) BuildGetPersonMentionsResponse(resp *lemmyResponse.GetPersonMentionsResponse) any {
	return &lemmyResponse017.GetPersonMentionsResponse{
		Mentions: helper.MapSlice(resp.Mentions, converter.ConvertPersonMentionViewTo017),
	}
}

// ParseGetRepliesRequest reuses the canonical parse directly — 0.17.2's
// real GetReplies matches exactly.
func (receiver *Frontend017) ParseGetRepliesRequest(request *http.Request) (*lemmyRequest.GetRepliesRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetRepliesRequest](request)
}

func (receiver *Frontend017) BuildGetRepliesResponse(resp *lemmyResponse.GetRepliesResponse) any {
	return &lemmyResponse017.GetRepliesResponse{
		Replies: helper.MapSlice(resp.Replies, converter.ConvertCommentReplyViewTo017),
	}
}

// ParseGetPrivateMessagesRequest reuses the canonical parse directly —
// 0.17.2's real GetPrivateMessages matches exactly.
func (receiver *Frontend017) ParseGetPrivateMessagesRequest(request *http.Request) (*lemmyRequest.GetPrivateMessagesRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetPrivateMessagesRequest](request)
}

func (receiver *Frontend017) BuildGetPrivateMessagesResponse(resp *lemmyResponse.GetPrivateMessagesResponse) any {
	return &lemmyResponse017.GetPrivateMessagesResponse{
		PrivateMessages: helper.MapSlice(resp.PrivateMessages, converter.ConvertPrivateMessageViewTo017),
	}
}

// ParseResolveObjectRequest reuses the canonical parse directly —
// 0.17.2's real ResolveObject (q, auth) matches exactly.
func (receiver *Frontend017) ParseResolveObjectRequest(request *http.Request) (*lemmyRequest.ResolveObjectRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.ResolveObjectRequest](request)
}

func (receiver *Frontend017) BuildResolveObjectResponse(resp *lemmyResponse.ResolveObjectResponse) any {
	result := &lemmyResponse017.ResolveObjectResponse{}
	if resp.Comment != nil {
		result.Comment = helper.ToPointer(converter.ConvertCommentViewTo017(*resp.Comment))
	}
	if resp.Post != nil {
		result.Post = helper.ToPointer(converter.ConvertPostViewTo017(*resp.Post))
	}
	if resp.Community != nil {
		result.Community = helper.ToPointer(converter.ConvertCommunityViewTo017(*resp.Community))
	}
	if resp.Person != nil {
		result.Person = helper.ToPointer(converter.ConvertPersonViewTo017(*resp.Person))
	}
	return result
}

func (receiver *Frontend017) BuildGetCaptchaResponse(resp *lemmyResponse.GetCaptchaResponse) any {
	result := &lemmyResponse017.GetCaptchaResponse{}
	if resp.Ok != nil {
		result.Ok = &lemmyResponse017.CaptchaResponse{
			Png:  resp.Ok.Png,
			Wav:  resp.Ok.Wav,
			Uuid: resp.Ok.Uuid,
		}
	}
	return result
}
