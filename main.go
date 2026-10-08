package main

import (
	"LemmyBeProxy/config"
	lemmyModel "LemmyBeProxy/dto/model/lemmy"
	lemmyResponse "LemmyBeProxy/dto/response/lemmy"
	"LemmyBeProxy/dto/response/piefed"
	"LemmyBeProxy/helper"
	appHttp "LemmyBeProxy/http"
	"LemmyBeProxy/router"
	lemmyService "LemmyBeProxy/service/lemmy"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM, syscall.SIGINT)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if config.CorsRegex != nil && config.CorsRegex.MatchString(request.Header.Get("Origin")) {
			writer.Header().Add("Access-Control-Allow-Origin", request.Header.Get("Origin"))
			writer.Header().Add("Access-Control-Allow-Headers", "Content-Type, User-Agent, Authorization")
		}

		if request.Method == "OPTIONS" {
			appHttp.WriteHttpResponse(appHttp.NoContent(), writer)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			appHttp.WriteHttpResponse(appHttp.InternalProxyError(), writer)
			log.Println(err)
			return
		}
		defer request.Body.Close()

		appRequest := appHttp.NewRequest(
			body,
			request.Method,
			helper.SingularizeMapValue(request.Header),
		)

		query := request.URL.Query()
		if query != nil {
			appRequest.QueryParams = helper.SingularizeMapValue(query)
		}

		// Lemmy 0.18.x and earlier clients send the auth token as an
		// "auth" field in the JSON body (POST/PUT) or as an "auth" query
		// parameter (GET), rather than an Authorization: Bearer header —
		// the convention Lemmy adopted from 0.19 onward and the only one
		// this proxy otherwise recognizes. If no Authorization header is
		// already present, check both older-style locations and
		// synthesize the header Piefed (and the rest of this proxy)
		// expects, so older Lemmy-API clients aren't silently treated as
		// unauthenticated.
		if _, hasAuth := appRequest.Headers["Authorization"]; !hasAuth {
			if token, ok := appRequest.QueryParams["auth"]; ok && token != "" {
				appRequest.Headers["Authorization"] = "Bearer " + token
			} else if len(body) > 0 {
				var bodyFields map[string]any
				if jsonErr := json.Unmarshal(body, &bodyFields); jsonErr == nil {
					if token, ok := bodyFields["auth"].(string); ok && token != "" {
						appRequest.Headers["Authorization"] = "Bearer " + token
					}
				}
			}
		}

		httpMethod, methodErr := router.HttpMethodFromString(request.Method)
		if methodErr != nil {
			// HEAD, TRACE, etc. — a client error, not a proxy failure.
			appHttp.WriteHttpResponse(&appHttp.Response{
				StatusCode: http.StatusMethodNotAllowed,
				Body:       map[string]string{"error": "method not allowed"},
			}, writer)
			return
		}

		var result *appHttp.Response
		for _, route := range config.AppRouter.Routes {
			matches, params, errRoute := router.RouteMatches(route, httpMethod, request.URL.Path)
			if errRoute != nil {
				err = errRoute
				break
			}
			if !matches {
				continue
			}

			appRequest.RouteParams = params
			result, err = route.ControllerMethod(appRequest)
			break
		}

		var responseError *piefed.ErrorResponse
		if errors.As(err, &responseError) {
			appHttp.WriteHttpResponse(&appHttp.Response{
				StatusCode: responseError.StatusCode,
				Body:       helper.ConvertPiefedErrorToLemmyError(responseError),
			}, writer)
			return
		}

		// Same idea as the Piefed branch above, for when BACKEND_TYPE=lemmy.
		// Real Lemmy's error body is already Lemmy-shaped, so no
		// translation is needed here — just forwarding its real status
		// code and error string instead of falling through to a generic
		// opaque 500.
		var lemmyApiError *lemmyService.LemmyApiError
		if errors.As(err, &lemmyApiError) {
			appHttp.WriteHttpResponse(&appHttp.Response{
				StatusCode: lemmyApiError.StatusCode,
				Body:       &lemmyResponse.ErrorResponse{Error: lemmyModel.ErrorCode(lemmyApiError.ErrorCode)},
			}, writer)
			return
		}

		if err != nil {
			appHttp.WriteHttpResponse(appHttp.InternalProxyError(), writer)
			log.Println(err)
			return
		}

		if result == nil {
			appHttp.WriteHttpResponse(appHttp.NotFoundProxyError(), writer)
			return
		}

		// See helper.StripTimezoneSuffixes for why this is needed: real
		// Lemmy 0.17.x rejects the timezone suffix our canonical model's
		// timestamps carry. This runs once here, after every controller
		// and Frontend017 conversion has already produced its normal
		// result, rather than being threaded through every individual
		// converter — see that function's own comment for the reasoning.
		//
		// Only JSON-object bodies are rewritten. A string or []byte body
		// is already final (a redirect's empty body, raw image bytes) —
		// re-marshaling it turned "" into the two-character JSON string
		// "\"\"" and would corrupt binary content.
		_, isString := result.Body.(string)
		_, isBytes := result.Body.([]byte)
		if config.FrontendVersion == "0.17" && result.Body != nil && !isString && !isBytes {
			marshaled, jsonErr := json.Marshal(result.Body)
			if jsonErr == nil {
				result.Body = string(helper.StripTimezoneSuffixes(marshaled))
			}
		}

		appHttp.WriteHttpResponse(result, writer)
	})

	go func() {
		var err error
		serveAddr := fmt.Sprintf(":%d", config.AppHttpPort)
		if config.CertificateFile != "" && config.CertificateKey != "" {
			err = http.ListenAndServeTLS(serveAddr, config.CertificateFile, config.CertificateKey, nil)
		} else {
			err = http.ListenAndServe(serveAddr, nil)
		}
		if err != nil {
			panic(err)
		}
	}()

	log.Println("Server started...")
	<-gracefulShutdown
	log.Println("Server shutting down...")
}
