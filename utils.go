package azdevops

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
)

func UnmarshalBody(response *http.Response, v any, unmarshalFunc func([]byte, any) error) (err error) {
	if response != nil && response.Body != nil {
		var err error
		defer func() {
			if closeError := response.Body.Close(); closeError != nil {
				err = closeError
			}
		}()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return unmarshalFunc(body, v)
	}
	return nil
}

func UnmarshalCollection(jsonValue []byte, v any) (err error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	} else {
		return errors.New("value type must be a pointer")
	}
	sType := reflect.StructOf([]reflect.StructField{
		{Name: "Count", Type: reflect.TypeOf(0)},
		{Name: "Value", Type: t},
	})
	sv := reflect.New(sType)
	if err = json.Unmarshal(jsonValue, sv.Interface()); err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	rv.Elem().Set(sv.Elem().FieldByName("Value"))
	return nil
}
