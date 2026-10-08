package router

import (
	"LemmyBeProxy/helper"
	"LemmyBeProxy/http"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type HttpMethod string

const (
	HttpMethodAll    HttpMethod = "*"
	HttpMethodGet    HttpMethod = "GET"
	HttpMethodPost   HttpMethod = "POST"
	HttpMethodPut    HttpMethod = "PUT"
	HttpMethodDelete HttpMethod = "DELETE"
	HttpMethodPatch  HttpMethod = "PATCH"
)

func HttpMethodFromString(method string) (HttpMethod, error) {
	allowed := helper.MapSlice(
		[]HttpMethod{HttpMethodGet, HttpMethodPost, HttpMethodPut, HttpMethodDelete, HttpMethodPatch},
		func(in HttpMethod) string {
			// Was `return method` — comparing the input against itself,
			// so every verb passed and lowercase verbs failed.
			return string(in)
		},
	)
	if !slices.Contains(allowed, strings.ToUpper(method)) {
		return HttpMethodAll, fmt.Errorf("invalid http method: %s", method)
	}

	return HttpMethod(strings.ToUpper(method)), nil
}

type ControllerMethod func(request *http.Request) (*http.Response, error)

type Route struct {
	Path             string
	HttpMethod       HttpMethod
	ControllerMethod ControllerMethod
	// pattern is compiled once here instead of on every request for
	// every route (RouteMatches used to recompile each route's regex per
	// request — ~30 compiles per incoming call).
	pattern *regexp.Regexp
}

func NewRoute(path string, httpMethod HttpMethod, controller ControllerMethod) *Route {
	pattern, err := RegexifyRoute(path)
	if err != nil {
		panic(fmt.Sprintf("invalid route %q: %v", path, err))
	}
	return &Route{Path: path, ControllerMethod: controller, HttpMethod: httpMethod, pattern: pattern}
}

type Router struct {
	Routes []*Route
}

func NewRouter() *Router {
	return &Router{
		Routes: make([]*Route, 0),
	}
}

func (receiver *Router) AddRoute(route *Route) {
	receiver.Routes = append(receiver.Routes, route)
}
