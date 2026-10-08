package http

import (
	"LemmyBeProxy/json"
	"log"
	"net/http"
	"strconv"
)

// WriteHttpResponse sets Content-Length explicitly rather than letting
// Go's server fall back to Transfer-Encoding: chunked (its default when
// no length is set). Confirmed directly: chunked responses from this
// server were read as truncated by Rust's reqwest client (lemmyBB),
// while curl read the exact same response correctly — an interop edge
// case in how reqwest's reader handles chunked encoding from this
// specific server, not a content bug. Setting Content-Length sidesteps
// the whole chunked-encoding code path rather than chasing the
// underlying interop cause.
func WriteHttpResponse(response *Response, writer http.ResponseWriter) {
	body := response.Body
	headers := response.Headers
	statusCode := response.StatusCode

	if headers == nil {
		headers = make(map[string]string)
	}
	if body == nil {
		body = make(map[string]string)
	}
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	_, ok := headers["Content-Type"]
	if !ok {
		headers["Content-Type"] = "application/json"
	}

	// Every response this proxy sends represents live, constantly-changing
	// forum data — nothing here should ever be cached by a downstream
	// client. Confirmed this matters concretely, not just in principle:
	// with no caching headers at all, a plain 200 OK is "cacheable by
	// default" per the HTTP spec (see http-cache-semantics's
	// is_storable()), so lemmyBB's own HTTP cache was storing every
	// distinct response it ever received from this proxy, unbounded,
	// eventually consuming gigabytes of memory and disk. Setting
	// no-store here means is_storable() returns false outright — the
	// entry is never written at all, fixing the actual root cause
	// rather than just periodically clearing lemmyBB's cache after the
	// fact.
	_, ok = headers["Cache-Control"]
	if !ok {
		headers["Cache-Control"] = "no-store"
	}

	var err error
	if _, isBytes := body.([]byte); isBytes {
		// Already-final raw content (e.g. proxied image bytes) — sent
		// as-is rather than JSON-encoded into a base64 string.
	} else if _, ok = body.(string); !ok {
		body, err = json.ToJson(body)
		if err != nil {
			body, _ = json.ToJson(map[string]string{
				"error": "Internal request error",
			})
			statusCode = http.StatusInternalServerError
			log.Println(err)
		}
	} else {
		body = []byte(body.(string))
	}

	bodyBytes := body.([]byte)
	headers["Content-Length"] = strconv.Itoa(len(bodyBytes))

	for key, value := range headers {
		writer.Header().Set(key, value)
	}
	writer.WriteHeader(statusCode)

	_, err = writer.Write(bodyBytes)
	if err != nil {
		log.Println(err)
	}
}
