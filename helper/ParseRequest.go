package helper

import (
	"LemmyBeProxy/dto/model"
	"LemmyBeProxy/http"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"strings"
)

var validate = validator.New()

func validateDto(result any) (err error) {
	err = validate.Struct(result)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			err = &model.ValidationError{
				Err:   validationErrs,
				Input: result,
			}
		}
	}

	return
}

func ParseRequest[T any](request *http.Request) (*T, error) {
	var result T
	err := json.Unmarshal(request.Body, &result)
	if err != nil {
		return nil, err
	}

	err = validateDto(result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// queryFieldKinds maps each json tag name of struct type t to the
// underlying reflect.Kind of that field (pointers dereferenced; named
// string types like SortType reduce to reflect.String).
func queryFieldKinds(t reflect.Type) map[string]reflect.Kind {
	kinds := make(map[string]reflect.Kind)
	if t == nil {
		return kinds
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return kinds
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			continue
		}
		ft := field.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		kinds[name] = ft.Kind()
	}
	return kinds
}

// ParseRequestQuery converts query parameters into T, choosing each
// value's JSON type from the target field's Go type rather than from the
// shape of the value. Guessing from the value broke every string field
// whose content looked like a number or bool: a search for "1986", "486"
// or "true", or a community named "2600", reached json.Unmarshal as a
// number/bool going into a string field and the request failed with a
// 500.
func ParseRequestQuery[T any](request *http.Request) (*T, error) {
	var result T
	kinds := queryFieldKinds(reflect.TypeOf(result))
	invalid := func() error {
		return &model.ValidationError{Input: result}
	}

	normalized := make(map[string]any)
	for key, val := range request.QueryParams {
		kind, known := kinds[key]
		if !known {
			// Not a field of T (e.g. the legacy "auth" param) — ignore.
			continue
		}

		switch kind {
		case reflect.String:
			normalized[key] = val
		case reflect.Bool:
			if val == "" {
				continue
			}
			b, err := strconv.ParseBool(strings.ToLower(val))
			if err != nil {
				return nil, invalid()
			}
			normalized[key] = b
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			if val == "" {
				continue
			}
			// json.Number keeps the exact digits; json.Unmarshal then
			// enforces the field's own range and sign.
			if _, err := strconv.ParseFloat(val, 64); err != nil {
				return nil, invalid()
			}
			normalized[key] = json.Number(val)
		default:
			normalized[key] = val
		}
	}

	jsonBytes, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	err = validateDto(result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
