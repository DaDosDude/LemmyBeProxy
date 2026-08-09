package controller

import (
	"LemmyBeProxy/dto/request/lemmy"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	"LemmyBeProxy/helper"
	"LemmyBeProxy/http"
	"LemmyBeProxy/service/backend"
	"LemmyBeProxy/service/frontend"
	goHttp "net/http"
)

// UserController is now thin on both axes for every endpoint except
// GetReportCount, which stays as a pure stub response with no
// backend call at all — there's nothing for either interface to
// actually do for it.
type UserController struct {
	backend  backend.Backend
	frontend frontend.Frontend
}

func NewUserController(backend backend.Backend, frontend frontend.Frontend) *UserController {
	return &UserController{
		backend:  backend,
		frontend: frontend,
	}
}

func (receiver *UserController) Login(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseLoginRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.Login(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildLoginResponse(resp)}, nil
}

func (receiver *UserController) Register(request *http.Request) (*http.Response, error) {
	reqDto, err := helper.ParseRequest[lemmy.RegisterRequest](request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.Register(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildLoginResponse(resp)}, nil
}

func (receiver *UserController) GetUnreadCount(request *http.Request) (*http.Response, error) {
	resp, err := receiver.backend.GetUnreadCount(request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildGetUnreadCountResponse(resp)}, nil
}

func (receiver *UserController) GetCaptcha(request *http.Request) (*http.Response, error) {
	resp, err := receiver.backend.GetCaptcha(request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildGetCaptchaResponse(resp)}, nil
}

func (receiver *UserController) GetReportCount(request *http.Request) (*http.Response, error) {
	reqDto, err := helper.ParseRequestQuery[lemmy.GetReportCountRequest](request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	return &http.Response{
		StatusCode: goHttp.StatusOK,
		Body: &lemmyResponse.GetReportCountResponse{
			CommentReports:        0,
			CommunityId:           reqDto.CommunityId,
			PostReports:           0,
			PrivateMessageReports: nil,
		},
	}, nil
}

func (receiver *UserController) GetUser(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseGetUserRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.GetUser(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildGetUserResponse(resp)}, nil
}

func (receiver *UserController) BlockPerson(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseBlockPersonRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.BlockPerson(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildBlockPersonResponse(resp)}, nil
}

func (receiver *UserController) SaveUserSettings(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseSaveUserSettingsRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.SaveUserSettings(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildSaveUserSettingsResponse(resp)}, nil
}

func (receiver *UserController) GetPersonMentions(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseGetPersonMentionsRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.GetPersonMentions(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildGetPersonMentionsResponse(resp)}, nil
}

func (receiver *UserController) GetReplies(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseGetRepliesRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.GetReplies(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildGetRepliesResponse(resp)}, nil
}

func (receiver *UserController) GetPrivateMessages(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseGetPrivateMessagesRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	resp, err := receiver.backend.GetPrivateMessages(reqDto, request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildGetPrivateMessagesResponse(resp)}, nil
}
