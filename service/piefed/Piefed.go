package piefed

import (
	piefedResponse "LemmyBeProxy/dto/response/piefed"
	"LemmyBeProxy/helper"
	"LemmyBeProxy/http"
	"LemmyBeProxy/router"
	"encoding/json"
	"fmt"
	"io"
	"log"
	goHttp "net/http"
	"strings"
)

const applicationJson = "application/json"

// truncateForLog caps how much of a response body gets logged, since some
// error bodies (an HTML error page, for instance) could be large — we only
// need enough to diagnose what went wrong, not the whole thing.
func truncateForLog(body []byte) string {
	const maxLen = 500
	if len(body) > maxLen {
		return string(body[:maxLen]) + "... (truncated)"
	}
	return string(body)
}

func defaultHandler[TResponse any](
	piefed *Piefed,
	urlPath string,
	httpMethod router.HttpMethod,
	request any,
	headers http.Headers,
) (*TResponse, error) {
	resp, err := piefed.sendRequest(
		urlPath,
		httpMethod,
		request,
		headers,
	)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != goHttp.StatusOK {
		var errorResponse *piefedResponse.ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil || errorResponse == nil {
			// Piefed's error body wasn't a JSON object we could parse
			// (empty body, plain text, an HTML error page, etc.) — build
			// a fallback instead of dereferencing a nil pointer, which
			// previously panicked here and killed the whole request
			// mid-flight (the client would see the connection die rather
			// than get any error response at all).
			errorResponse = &piefedResponse.ErrorResponse{
				ErrorCode: fmt.Sprintf("unparseable_error_response (http %d)", resp.StatusCode),
			}
		}
		errorResponse.StatusCode = resp.StatusCode

		// Log every non-200 with enough detail to diagnose it after the
		// fact — previously these passed through silently (the client just
		// saw a generic error), which made intermittent failures like this
		// unreproducible from server-side logs alone. resp.Request.URL
		// includes the actual query string that was sent, which the bare
		// urlPath argument doesn't — that's the difference between
		// guessing which parameter was invalid and knowing immediately.
		log.Printf(
			"piefed error response: %s %s -> HTTP %d, error_code=%q, body=%s",
			httpMethod, resp.Request.URL.String(), resp.StatusCode, errorResponse.ErrorCode, truncateForLog(body),
		)

		return nil, errorResponse
	}

	var response TResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// A 200 from Piefed whose body doesn't match the Go struct we
		// expected — this is exactly the class of bug the negative
		// downvotes-count issue was (a uint field fed a value Piefed's
		// data can actually produce). Logging the raw body here means the
		// next occurrence of a shape mismatch is diagnosable immediately
		// instead of requiring another round of manual reproduction.
		log.Printf(
			"piefed response unmarshal failed: %s %s -> %v, body=%s",
			httpMethod, urlPath, err, truncateForLog(body),
		)
		return nil, err
	}

	return &response, nil
}

type Piefed struct {
	instance string
}

func NewPiefed(instance string) *Piefed {
	return &Piefed{
		instance: instance,
	}
}

// baseURL defaults to https:// but, like the Lemmy client, honours an
// explicit scheme in BACKEND_INSTANCE (e.g. "http://piefed-app:5000" to
// reach the backend over an internal network).
func (receiver *Piefed) baseURL() string {
	if strings.HasPrefix(receiver.instance, "http://") || strings.HasPrefix(receiver.instance, "https://") {
		return strings.TrimSuffix(receiver.instance, "/")
	}
	return "https://" + receiver.instance
}

func (receiver *Piefed) url() string {
	return receiver.baseURL() + "/api/alpha"
}

func (receiver *Piefed) sendRequest(
	path string,
	method router.HttpMethod,
	request any,
	headers http.Headers,
) (*goHttp.Response, error) {
	var body io.Reader
	queryString := ""
	if request != nil {
		if method == router.HttpMethodGet {
			marshalled, err := helper.MarshalToQueryString(request)
			if err != nil {
				return nil, err
			}
			queryString = "?" + marshalled
		} else {
			body = helper.MarshalToReader(request)
		}
	}

	req, err := goHttp.NewRequest(
		string(method),
		receiver.url()+path+queryString,
		body,
	)
	if err != nil {
		return nil, err
	}

	// Log the exact outbound request we're about to send. This matters
	// specifically for GET requests — the query string is built by
	// helper.MarshalToQueryString from a Go struct, and if that produces
	// something Piefed doesn't recognize (a malformed or unexpected value
	// for a field), the error response alone doesn't show what we
	// actually sent, only what Piefed's validator complained about.
	log.Printf("piefed request: %s %s", method, receiver.url()+path+queryString)

	req.Header = helper.MapMap(helper.MapFilter(headers, func(value, key string) bool {
		lowerKey := strings.ToLower(key)
		// Content-Length is recalculated by the Go HTTP client for this
		// request's own body, so a stale client-supplied value must not
		// be forwarded. Accept-Encoding must not be forwarded either —
		// if a client-set value is present, Go's transport disables its
		// own automatic gzip request/decompress handling, and a gzip
		// response body then gets handed to json.Unmarshal as raw bytes
		// (visible as "invalid character '\x1f'..." — 0x1f is gzip's
		// magic byte). Leaving this header unset lets Go manage
		// compression transparently, which is what defaultHandler's
		// plain json.Unmarshal on the response body assumes.
		return lowerKey != "content-length" && lowerKey != "accept-encoding"
	}), func(value, key string) []string {
		return []string{value}
	})
	if body != nil {
		req.Header.Set("Content-Type", applicationJson)
	}

	return goHttp.DefaultClient.Do(req)
}
