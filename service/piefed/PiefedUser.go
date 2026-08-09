package piefed

import (
	piefedRequest "LemmyBeProxy/dto/request/piefed"
	piefedResponse "LemmyBeProxy/dto/response/piefed"
	appHttp "LemmyBeProxy/http"
	"LemmyBeProxy/router"
)

func (receiver *Piefed) Login(request *piefedRequest.LoginRequest, headers appHttp.Headers) (*piefedResponse.LoginResponse, error) {
	return defaultHandler[piefedResponse.LoginResponse](
		receiver,
		"/user/login",
		router.HttpMethodPost,
		request,
		headers,
	)
}

// Register — path assumed to match Lemmy's own "/user/register" naming
// 1:1, per every other endpoint's established convention. Not directly
// source-verified against Piefed.
func (receiver *Piefed) Register(request *piefedRequest.RegisterRequest, headers appHttp.Headers) (*piefedResponse.RegisterResponse, error) {
	return defaultHandler[piefedResponse.RegisterResponse](
		receiver,
		"/user/register",
		router.HttpMethodPost,
		request,
		headers,
	)
}

func (receiver *Piefed) GetUnreadCount(headers appHttp.Headers) (*piefedResponse.GetUnreadCountResponse, error) {
	return defaultHandler[piefedResponse.GetUnreadCountResponse](
		receiver,
		"/user/unread_count",
		router.HttpMethodGet,
		nil,
		headers,
	)
}

func (receiver *Piefed) GetUser(request *piefedRequest.GetUserRequest, headers appHttp.Headers) (*piefedResponse.GetUserResponse, error) {
	return defaultHandler[piefedResponse.GetUserResponse](
		receiver,
		"/user",
		router.HttpMethodGet,
		request,
		headers,
	)
}

func (receiver *Piefed) BlockPerson(request *piefedRequest.BlockPersonRequest, headers appHttp.Headers) (*piefedResponse.BlockPersonResponse, error) {
	return defaultHandler[piefedResponse.BlockPersonResponse](
		receiver,
		"/user/block",
		router.HttpMethodPost,
		request,
		headers,
	)
}

func (receiver *Piefed) SaveUserSettings(request *piefedRequest.SaveUserSettingsRequest, headers appHttp.Headers) (*piefedResponse.SaveUserSettingsResponse, error) {
	return defaultHandler[piefedResponse.SaveUserSettingsResponse](
		receiver,
		"/user/save_user_settings",
		router.HttpMethodPut,
		request,
		headers,
	)
}
