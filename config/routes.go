package config

import (
	"LemmyBeProxy/controller"
	"LemmyBeProxy/router"
)

func newRoute(path string, method router.HttpMethod, controller router.ControllerMethod) *router.Route {
	return router.NewRoute("/api/v3"+path, method, controller)
}

func init() {
	userController := controller.NewUserController(activeBackend, activeFrontend)
	siteController := controller.NewSiteController(activeBackend, activeFrontend)
	postController := controller.NewPostController(activeBackend, activeFrontend)
	commentController := controller.NewCommentController(activeBackend, activeFrontend)
	communityController := controller.NewCommunityController(activeBackend, activeFrontend)
	searchController := controller.NewSearchController(activeBackend, activeFrontend)
	uploadController := controller.NewUploadController(activeBackend, imageHosts)

	// implemented
	AppRouter.AddRoute(newRoute("/user/login", router.HttpMethodPost, userController.Login))
	AppRouter.AddRoute(newRoute("/user/unread_count", router.HttpMethodGet, userController.GetUnreadCount))
	AppRouter.AddRoute(newRoute("/user", router.HttpMethodGet, userController.GetUser))
	AppRouter.AddRoute(newRoute("/user/block", router.HttpMethodPost, userController.BlockPerson))
	AppRouter.AddRoute(newRoute("/user/save_user_settings", router.HttpMethodPut, userController.SaveUserSettings))
	AppRouter.AddRoute(newRoute("/user/mention", router.HttpMethodGet, userController.GetPersonMentions))
	AppRouter.AddRoute(newRoute("/user/replies", router.HttpMethodGet, userController.GetReplies))
	AppRouter.AddRoute(newRoute("/private_message/list", router.HttpMethodGet, userController.GetPrivateMessages))
	AppRouter.AddRoute(newRoute("/user/get_captcha", router.HttpMethodGet, userController.GetCaptcha))
	AppRouter.AddRoute(newRoute("/user/mark_all_as_read", router.HttpMethodPost, userController.MarkAllAsRead))
	AppRouter.AddRoute(newRoute("/site", router.HttpMethodGet, siteController.Site))
	AppRouter.AddRoute(newRoute("/resolve_object", router.HttpMethodGet, siteController.ResolveObject))
	AppRouter.AddRoute(newRoute("/post/list", router.HttpMethodGet, postController.GetPosts))
	AppRouter.AddRoute(newRoute("/post", router.HttpMethodGet, postController.GetPost))
	AppRouter.AddRoute(newRoute("/post/like", router.HttpMethodPost, postController.LikePost))
	AppRouter.AddRoute(newRoute("/post/mark_as_read", router.HttpMethodPost, postController.MarkPostAsRead))
	AppRouter.AddRoute(newRoute("/post", router.HttpMethodPost, postController.CreatePost))
	AppRouter.AddRoute(newRoute("/post", router.HttpMethodPut, postController.EditPost))
	AppRouter.AddRoute(newRoute("/comment/list", router.HttpMethodGet, commentController.GetComments))
	AppRouter.AddRoute(newRoute("/comment", router.HttpMethodGet, commentController.GetComment))
	AppRouter.AddRoute(newRoute("/comment", router.HttpMethodPost, commentController.CreateComment))
	AppRouter.AddRoute(newRoute("/comment/like", router.HttpMethodPost, commentController.LikeComment))
	AppRouter.AddRoute(newRoute("/community", router.HttpMethodGet, communityController.GetCommunity))
	AppRouter.AddRoute(newRoute("/community/list", router.HttpMethodGet, communityController.GetCommunities))
	AppRouter.AddRoute(newRoute("/community/follow", router.HttpMethodPost, communityController.FollowCommunity))
	AppRouter.AddRoute(newRoute("/community/block", router.HttpMethodPost, communityController.BlockCommunity))
	AppRouter.AddRoute(newRoute("/search", router.HttpMethodGet, searchController.Search))
	AppRouter.AddRoute(newRoute("/user/register", router.HttpMethodPost, userController.Register))

	// These two intentionally bypass the /api/v3 prefix — mlmym (and real
	// Lemmy pict-rs) upload/serve images at the site root, not under the API.
	AppRouter.AddRoute(router.NewRoute("/pictrs/image", router.HttpMethodPost, uploadController.UploadImage))
	AppRouter.AddRoute(router.NewRoute("/pictrs/image/{token}", router.HttpMethodGet, uploadController.ServeImage))

	// stub: always reports zero reports, no backend call
	AppRouter.AddRoute(newRoute("/user/report_count", router.HttpMethodGet, userController.GetReportCount))
}
