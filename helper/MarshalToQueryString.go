package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
)

func MarshalToQueryString(in any) (result string, err error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return
	}

	normalized := make(map[string]any)
	err = json.Unmarshal(bytes, &normalized)
	if err != nil {
		return
	}

	queryParts := make([]string, 0, len(normalized))
	for key, value := range normalized {
		if value == nil {
			continue
		}

		var valueStr string
		if _, ok := value.(bool); ok {
			valueStr = fmt.Sprintf("%t", value)
		} else if _, ok := value.(string); ok {
			valueStr = value.(string)
		} else if _, ok := value.(int); ok {
			valueStr = fmt.Sprintf("%d", value)
		} else if _, ok := value.(uint); ok {
			valueStr = fmt.Sprintf("%d", value)
		} else if _, ok := value.(float64); ok {
			if value.(float64) == math.Trunc(value.(float64)) {
				valueStr = fmt.Sprintf("%d", int64(value.(float64)))
			} else {
				valueStr = fmt.Sprintf("%f", value)
			}
		} else {
			err = errors.New(fmt.Sprintf("Unknown type: %T", value))
			return
		}

		// Every value must be URL-encoded — a raw, unencoded space
		// (or &, #, %, etc.) corrupts the outgoing request at the HTTP
		// protocol level, not just the query semantics: real Lemmy
		// never even saw a well-formed request for a multi-word search
		// query, since a bare space in a URL breaks the request line
		// itself. Confirmed directly: search with spaces returned an
		// empty, unparseable response body until this was added.
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(valueStr)))
	}

	return strings.Join(queryParts, "&"), nil
}
