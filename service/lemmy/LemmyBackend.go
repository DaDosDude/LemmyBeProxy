package lemmy

import (
	lemmyRequest "LemmyBeProxy/dto/request/lemmy"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	appHttp "LemmyBeProxy/http"
	"LemmyBeProxy/router"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	goHttp "net/http"
)

// LemmyBackend implements the backend.Backend interface against a real
// Lemmy instance. Unlike PiefedBackend, this needs almost no translation:
// this proxy's canonical DTOs are already modeled on Lemmy's real field
// names, so requests and responses are forwarded close to as-is. The one
// exception is the Piefed-specific community_id + type_=Subscribed
// workaround in PiefedBackend, which doesn't apply here — that was a
// real Piefed backend bug, not a general Lemmy API quirk, so real Lemmy
// gets the request exactly as the client sent it.
type LemmyBackend struct {
	client *Lemmy
}

func NewLemmyBackend(client *Lemmy) *LemmyBackend {
	return &LemmyBackend{client: client}
}

func (receiver *LemmyBackend) GetPosts(request *lemmyRequest.GetPostsRequest, headers appHttp.Headers) (*lemmyResponse.GetPostsResponse, error) {
	return defaultHandler[lemmyResponse.GetPostsResponse](receiver.client, "/post/list", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) GetPost(request *lemmyRequest.GetPostRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	return defaultHandler[lemmyResponse.GetPostResponse](receiver.client, "/post", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) CreatePost(request *lemmyRequest.CreatePostRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	return defaultHandler[lemmyResponse.GetPostResponse](receiver.client, "/post", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) EditPost(request *lemmyRequest.EditPostRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	return defaultHandler[lemmyResponse.GetPostResponse](receiver.client, "/post", router.HttpMethodPut, request, headers)
}

func (receiver *LemmyBackend) LikePost(request *lemmyRequest.CreatePostLikeRequest, headers appHttp.Headers) (*lemmyResponse.GetPostResponse, error) {
	return defaultHandler[lemmyResponse.GetPostResponse](receiver.client, "/post/like", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) MarkPostAsRead(request *lemmyRequest.MarkPostAsReadRequest, headers appHttp.Headers) (*lemmyResponse.SuccessResponse, error) {
	return defaultHandler[lemmyResponse.SuccessResponse](receiver.client, "/post/mark_as_read", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) GetComments(request *lemmyRequest.GetCommentsRequest, headers appHttp.Headers) (*lemmyResponse.GetCommentsResponse, error) {
	return defaultHandler[lemmyResponse.GetCommentsResponse](receiver.client, "/comment/list", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) GetComment(request *lemmyRequest.GetCommentRequest, headers appHttp.Headers) (*lemmyResponse.GetCommentResponse, error) {
	return defaultHandler[lemmyResponse.GetCommentResponse](receiver.client, "/comment", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) CreateComment(request *lemmyRequest.CreateCommentRequest, headers appHttp.Headers) (*lemmyResponse.CreateCommentResponse, error) {
	return defaultHandler[lemmyResponse.CreateCommentResponse](receiver.client, "/comment", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) LikeComment(request *lemmyRequest.CreateCommentLikeRequest, headers appHttp.Headers) (*lemmyResponse.GetCommentResponse, error) {
	return defaultHandler[lemmyResponse.GetCommentResponse](receiver.client, "/comment/like", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) GetCommunity(request *lemmyRequest.GetCommunityRequest, headers appHttp.Headers) (*lemmyResponse.GetCommunityResponse, error) {
	return defaultHandler[lemmyResponse.GetCommunityResponse](receiver.client, "/community", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) GetCommunities(request *lemmyRequest.GetCommunitiesRequest, headers appHttp.Headers) (*lemmyResponse.GetCommunitiesResponse, error) {
	return defaultHandler[lemmyResponse.GetCommunitiesResponse](receiver.client, "/community/list", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) FollowCommunity(request *lemmyRequest.FollowCommunityRequest, headers appHttp.Headers) (*lemmyResponse.CommunityResponse, error) {
	return defaultHandler[lemmyResponse.CommunityResponse](receiver.client, "/community/follow", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) BlockCommunity(request *lemmyRequest.BlockCommunityRequest, headers appHttp.Headers) (*lemmyResponse.BlockCommunityResponse, error) {
	return defaultHandler[lemmyResponse.BlockCommunityResponse](receiver.client, "/community/block", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) Login(request *lemmyRequest.LoginRequest, headers appHttp.Headers) (*lemmyResponse.LoginResponse, error) {
	return defaultHandler[lemmyResponse.LoginResponse](receiver.client, "/user/login", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) Register(request *lemmyRequest.RegisterRequest, headers appHttp.Headers) (*lemmyResponse.LoginResponse, error) {
	return defaultHandler[lemmyResponse.LoginResponse](receiver.client, "/user/register", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) GetUnreadCount(headers appHttp.Headers) (*lemmyResponse.GetUnreadCountResponse, error) {
	return defaultHandler[lemmyResponse.GetUnreadCountResponse](receiver.client, "/user/unread_count", router.HttpMethodGet, nil, headers)
}

func (receiver *LemmyBackend) GetUser(request *lemmyRequest.GetUserRequest, headers appHttp.Headers) (*lemmyResponse.GetUserResponse, error) {
	return defaultHandler[lemmyResponse.GetUserResponse](receiver.client, "/user", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) BlockPerson(request *lemmyRequest.BlockPersonRequest, headers appHttp.Headers) (*lemmyResponse.BlockPersonResponse, error) {
	return defaultHandler[lemmyResponse.BlockPersonResponse](receiver.client, "/user/block", router.HttpMethodPost, request, headers)
}

func (receiver *LemmyBackend) SaveUserSettings(request *lemmyRequest.SaveUserSettingsRequest, headers appHttp.Headers) (*lemmyResponse.SaveUserSettingsResponse, error) {
	return defaultHandler[lemmyResponse.SaveUserSettingsResponse](receiver.client, "/user/save_user_settings", router.HttpMethodPut, request, headers)
}

func (receiver *LemmyBackend) Search(request *lemmyRequest.SearchRequest, headers appHttp.Headers) (*lemmyResponse.SearchResponse, error) {
	return defaultHandler[lemmyResponse.SearchResponse](receiver.client, "/search", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) Site(headers appHttp.Headers) (*lemmyResponse.GetSiteResponse, error) {
	return defaultHandler[lemmyResponse.GetSiteResponse](receiver.client, "/site", router.HttpMethodGet, nil, headers)
}

// UploadImage forwards directly to real Lemmy's own pict-rs endpoint at
// /pictrs/image — outside /api/v3 entirely, the same path a real mlmym
// client talking to real Lemmy would use, with the same Cookie-based
// auth (unlike Piefed's /upload/image, which expects a Bearer header —
// see PiefedBackend's implementation for that difference). Real Lemmy's
// pict-rs already returns a servable file token directly in its
// response, so the resulting URL is built from that token rather than
// needing any translation.
func (receiver *LemmyBackend) UploadImage(fileBytes []byte, filename string, jwt string) (string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("images[]", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(fileBytes); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := goHttp.NewRequest(
		"POST",
		fmt.Sprintf("%s/pictrs/image", receiver.client.baseURL()),
		body,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", "jwt="+jwt)

	res, err := goHttp.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// Real Lemmy's pict-rs returns 201 Created for a successful upload,
	// not 200 — confirmed directly against retrolemmy.com. Piefed's own
	// upload endpoint returns 200, which is why this distinction matters
	// here specifically and not in PiefedBackend.
	if res.StatusCode != goHttp.StatusOK && res.StatusCode != goHttp.StatusCreated {
		return "", fmt.Errorf("lemmy image upload failed with status %d: %s", res.StatusCode, string(respBody))
	}

	var response struct {
		Files []struct {
			File string `json:"file"`
		} `json:"files"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}
	if len(response.Files) == 0 {
		return "", fmt.Errorf("lemmy image upload succeeded but returned no files")
	}

	return fmt.Sprintf("%s/pictrs/image/%s", receiver.client.baseURL(), response.Files[0].File), nil
}

func (receiver *LemmyBackend) GetPersonMentions(request *lemmyRequest.GetPersonMentionsRequest, headers appHttp.Headers) (*lemmyResponse.GetPersonMentionsResponse, error) {
	return defaultHandler[lemmyResponse.GetPersonMentionsResponse](receiver.client, "/user/mention", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) GetReplies(request *lemmyRequest.GetRepliesRequest, headers appHttp.Headers) (*lemmyResponse.GetRepliesResponse, error) {
	return defaultHandler[lemmyResponse.GetRepliesResponse](receiver.client, "/user/replies", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) GetPrivateMessages(request *lemmyRequest.GetPrivateMessagesRequest, headers appHttp.Headers) (*lemmyResponse.GetPrivateMessagesResponse, error) {
	return defaultHandler[lemmyResponse.GetPrivateMessagesResponse](receiver.client, "/private_message/list", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) ResolveObject(request *lemmyRequest.ResolveObjectRequest, headers appHttp.Headers) (*lemmyResponse.ResolveObjectResponse, error) {
	return defaultHandler[lemmyResponse.ResolveObjectResponse](receiver.client, "/resolve_object", router.HttpMethodGet, request, headers)
}

func (receiver *LemmyBackend) GetCaptcha(headers appHttp.Headers) (*lemmyResponse.GetCaptchaResponse, error) {
	return defaultHandler[lemmyResponse.GetCaptchaResponse](receiver.client, "/user/get_captcha", router.HttpMethodGet, nil, headers)
}
