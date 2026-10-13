package controller

import (
	"LemmyBeProxy/cache"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	"LemmyBeProxy/helper"
	"LemmyBeProxy/http"
	"LemmyBeProxy/service/backend"
	"LemmyBeProxy/service/frontend"
	goHttp "net/http"
)

// SiteController is now thin — all the Piefed-specific assembly logic
// (including the ActivityPub actor fetch) moved into PiefedBackend,
// since that's genuinely backend-specific work a real Lemmy backend
// doesn't need at all.
type SiteController struct {
	backend  backend.Backend
	frontend frontend.Frontend

	// resolveCache holds anonymous resolve_object results (nil = caching
	// disabled). lemmyBB resolves every community on its index page on
	// every view — 60-80 backend calls per page, all in one backend
	// rate-limit bucket since lemmyBB doesn't forward visitor IPs.
	resolveCache *cache.TTLCache[*lemmyResponse.ResolveObjectResponse]
	// cacheKeyPrefix scopes keys to the configured backend.
	cacheKeyPrefix string
}

func NewSiteController(
	backend backend.Backend,
	frontend frontend.Frontend,
	resolveCache *cache.TTLCache[*lemmyResponse.ResolveObjectResponse],
	cacheKeyPrefix string,
) *SiteController {
	return &SiteController{
		backend:        backend,
		frontend:       frontend,
		resolveCache:   resolveCache,
		cacheKeyPrefix: cacheKeyPrefix,
	}
}

// isAnonymous reports whether a request carries no credentials at all.
// Authenticated resolve_object responses can contain per-user state
// (subscribed, blocked, my_vote, saved), so only fully anonymous requests
// may share cached results. Any Cookie counts as possibly authenticated
// too: real Lemmy accepts a jwt cookie, and PieFed has session cookies.
// The legacy auth query/body param is already turned into an
// Authorization header in main.go before this runs.
func isAnonymous(headers http.Headers) bool {
	return headers["Authorization"] == "" && headers["Cookie"] == ""
}

func (receiver *SiteController) Site(request *http.Request) (*http.Response, error) {
	resp, err := receiver.backend.Site(request.Headers)
	if err != nil {
		return nil, err
	}

	return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildGetSiteResponse(resp)}, nil
}

func (receiver *SiteController) ResolveObject(request *http.Request) (*http.Response, error) {
	reqDto, err := receiver.frontend.ParseResolveObjectRequest(request)
	if err != nil {
		return helper.ConvertValidationErrorsToResponse(err), nil
	}

	load := func() (*lemmyResponse.ResolveObjectResponse, error) {
		return receiver.backend.ResolveObject(reqDto, request.Headers)
	}

	if !receiver.resolveCache.Enabled() || !isAnonymous(request.Headers) {
		resp, err := load()
		if err != nil {
			return nil, err
		}
		return &http.Response{StatusCode: goHttp.StatusOK, Body: receiver.frontend.BuildResolveObjectResponse(resp)}, nil
	}

	// Errors are returned but never cached (see TTLCache.GetOrLoad).
	resp, hit, err := receiver.resolveCache.GetOrLoad(receiver.cacheKeyPrefix+reqDto.Q, load)
	if err != nil {
		return nil, err
	}

	cacheStatus := "MISS"
	if hit {
		cacheStatus = "HIT"
	}
	return &http.Response{
		StatusCode: goHttp.StatusOK,
		Headers:    map[string]string{"X-Proxy-Cache": cacheStatus},
		Body:       receiver.frontend.BuildResolveObjectResponse(resp),
	}, nil
}
