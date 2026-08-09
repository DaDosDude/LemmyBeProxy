package frontend

import (
	lemmyRequest "LemmyBeProxy/dto/request/lemmy"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	"LemmyBeProxy/helper"
	"LemmyBeProxy/http"
)

// Frontend019 is the current-generation Lemmy wire format — this proxy's
// canonical shape already is this format, so every method here is a
// direct parse/return with no translation at all. This is exactly what
// every controller did inline before the Frontend interface existed.
type Frontend019 struct{}

func NewFrontend019() *Frontend019 {
	return &Frontend019{}
}

func (receiver *Frontend019) ParseGetPostsRequest(request *http.Request) (*lemmyRequest.GetPostsRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetPostsRequest](request)
}

func (receiver *Frontend019) BuildGetPostsResponse(resp *lemmyResponse.GetPostsResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetPostRequest(request *http.Request) (*lemmyRequest.GetPostRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetPostRequest](request)
}

func (receiver *Frontend019) BuildGetPostResponse(resp *lemmyResponse.GetPostResponse) any {
	return resp
}

func (receiver *Frontend019) ParseCreatePostRequest(request *http.Request) (*lemmyRequest.CreatePostRequest, error) {
	return helper.ParseRequest[lemmyRequest.CreatePostRequest](request)
}

func (receiver *Frontend019) ParseEditPostRequest(request *http.Request) (*lemmyRequest.EditPostRequest, error) {
	return helper.ParseRequest[lemmyRequest.EditPostRequest](request)
}

func (receiver *Frontend019) ParseCreatePostLikeRequest(request *http.Request) (*lemmyRequest.CreatePostLikeRequest, error) {
	return helper.ParseRequest[lemmyRequest.CreatePostLikeRequest](request)
}

func (receiver *Frontend019) BuildPostMutationResponse(resp *lemmyResponse.GetPostResponse) any {
	return resp
}

func (receiver *Frontend019) ParseMarkPostAsReadRequest(request *http.Request) (*lemmyRequest.MarkPostAsReadRequest, error) {
	return helper.ParseRequest[lemmyRequest.MarkPostAsReadRequest](request)
}

func (receiver *Frontend019) BuildSuccessResponse(resp *lemmyResponse.SuccessResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetCommentsRequest(request *http.Request) (*lemmyRequest.GetCommentsRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetCommentsRequest](request)
}

func (receiver *Frontend019) BuildGetCommentsResponse(resp *lemmyResponse.GetCommentsResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetCommentRequest(request *http.Request) (*lemmyRequest.GetCommentRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetCommentRequest](request)
}

func (receiver *Frontend019) BuildGetCommentResponse(resp *lemmyResponse.GetCommentResponse) any {
	return resp
}

func (receiver *Frontend019) ParseCreateCommentRequest(request *http.Request) (*lemmyRequest.CreateCommentRequest, error) {
	return helper.ParseRequest[lemmyRequest.CreateCommentRequest](request)
}

func (receiver *Frontend019) BuildCreateCommentResponse(resp *lemmyResponse.CreateCommentResponse) any {
	return resp
}

func (receiver *Frontend019) ParseCreateCommentLikeRequest(request *http.Request) (*lemmyRequest.CreateCommentLikeRequest, error) {
	return helper.ParseRequest[lemmyRequest.CreateCommentLikeRequest](request)
}

func (receiver *Frontend019) BuildLikeCommentResponse(resp *lemmyResponse.GetCommentResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetCommunityRequest(request *http.Request) (*lemmyRequest.GetCommunityRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetCommunityRequest](request)
}

func (receiver *Frontend019) BuildGetCommunityResponse(resp *lemmyResponse.GetCommunityResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetCommunitiesRequest(request *http.Request) (*lemmyRequest.GetCommunitiesRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetCommunitiesRequest](request)
}

func (receiver *Frontend019) BuildGetCommunitiesResponse(resp *lemmyResponse.GetCommunitiesResponse) any {
	return resp
}

func (receiver *Frontend019) ParseFollowCommunityRequest(request *http.Request) (*lemmyRequest.FollowCommunityRequest, error) {
	return helper.ParseRequest[lemmyRequest.FollowCommunityRequest](request)
}

func (receiver *Frontend019) BuildCommunityResponse(resp *lemmyResponse.CommunityResponse) any {
	return resp
}

func (receiver *Frontend019) ParseBlockCommunityRequest(request *http.Request) (*lemmyRequest.BlockCommunityRequest, error) {
	return helper.ParseRequest[lemmyRequest.BlockCommunityRequest](request)
}

func (receiver *Frontend019) BuildBlockCommunityResponse(resp *lemmyResponse.BlockCommunityResponse) any {
	return resp
}

func (receiver *Frontend019) ParseLoginRequest(request *http.Request) (*lemmyRequest.LoginRequest, error) {
	return helper.ParseRequest[lemmyRequest.LoginRequest](request)
}

func (receiver *Frontend019) BuildLoginResponse(resp *lemmyResponse.LoginResponse) any {
	return resp
}

func (receiver *Frontend019) BuildGetUnreadCountResponse(resp *lemmyResponse.GetUnreadCountResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetUserRequest(request *http.Request) (*lemmyRequest.GetUserRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetUserRequest](request)
}

func (receiver *Frontend019) BuildGetUserResponse(resp *lemmyResponse.GetUserResponse) any {
	return resp
}

func (receiver *Frontend019) ParseBlockPersonRequest(request *http.Request) (*lemmyRequest.BlockPersonRequest, error) {
	return helper.ParseRequest[lemmyRequest.BlockPersonRequest](request)
}

func (receiver *Frontend019) BuildBlockPersonResponse(resp *lemmyResponse.BlockPersonResponse) any {
	return resp
}

func (receiver *Frontend019) ParseSearchRequest(request *http.Request) (*lemmyRequest.SearchRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.SearchRequest](request)
}

func (receiver *Frontend019) BuildSearchResponse(resp *lemmyResponse.SearchResponse) any {
	return resp
}

func (receiver *Frontend019) BuildGetSiteResponse(resp *lemmyResponse.GetSiteResponse) any {
	return resp
}

func (receiver *Frontend019) ParseSaveUserSettingsRequest(request *http.Request) (*lemmyRequest.SaveUserSettingsRequest, error) {
	return helper.ParseRequest[lemmyRequest.SaveUserSettingsRequest](request)
}

func (receiver *Frontend019) BuildSaveUserSettingsResponse(resp *lemmyResponse.SaveUserSettingsResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetPersonMentionsRequest(request *http.Request) (*lemmyRequest.GetPersonMentionsRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetPersonMentionsRequest](request)
}

func (receiver *Frontend019) BuildGetPersonMentionsResponse(resp *lemmyResponse.GetPersonMentionsResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetRepliesRequest(request *http.Request) (*lemmyRequest.GetRepliesRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetRepliesRequest](request)
}

func (receiver *Frontend019) BuildGetRepliesResponse(resp *lemmyResponse.GetRepliesResponse) any {
	return resp
}

func (receiver *Frontend019) ParseGetPrivateMessagesRequest(request *http.Request) (*lemmyRequest.GetPrivateMessagesRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.GetPrivateMessagesRequest](request)
}

func (receiver *Frontend019) BuildGetPrivateMessagesResponse(resp *lemmyResponse.GetPrivateMessagesResponse) any {
	return resp
}

func (receiver *Frontend019) ParseResolveObjectRequest(request *http.Request) (*lemmyRequest.ResolveObjectRequest, error) {
	return helper.ParseRequestQuery[lemmyRequest.ResolveObjectRequest](request)
}

func (receiver *Frontend019) BuildResolveObjectResponse(resp *lemmyResponse.ResolveObjectResponse) any {
	return resp
}

func (receiver *Frontend019) BuildGetCaptchaResponse(resp *lemmyResponse.GetCaptchaResponse) any {
	return resp
}
