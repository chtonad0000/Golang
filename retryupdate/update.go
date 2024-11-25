//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

var lastSetVersion uuid.UUID

func ErrorHandler(getResp *kvapi.GetResponse, err error, c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	if err != nil {
		var apiErr *kvapi.APIError
		if errors.As(err, &apiErr) {
			if errors.Is(apiErr.Unwrap(), kvapi.ErrKeyNotFound) {
				val, updErr := updateFn(nil)
				if updErr != nil {
					return updErr
				}
				_, err = c.Set(&kvapi.SetRequest{Key: key, Value: val, OldVersion: uuid.UUID{}, NewVersion: uuid.Must(uuid.NewV4())})
				if err == nil {
					return nil
				}
				return ErrorHandler(getResp, err, c, key, updateFn)
			}
			var authErr *kvapi.AuthError
			if errors.As(apiErr.Unwrap(), &authErr) {
				return apiErr
			}
			var conflictErr *kvapi.ConflictError
			if errors.As(apiErr.Unwrap(), &conflictErr) {
				if lastSetVersion == conflictErr.ExpectedVersion {
					return nil
				}
				getResp, err = c.Get(&kvapi.GetRequest{Key: key})
				return ErrorHandler(getResp, err, c, key, updateFn)
			}
			if apiErr.Method == "get" {
				getResp, err = c.Get(&kvapi.GetRequest{Key: key})
				return ErrorHandler(getResp, err, c, key, updateFn)
			}

			if apiErr.Method == "set" {
				val, updErr := updateFn(&getResp.Value)
				if updErr != nil {
					return updErr
				}
				_, err = c.Set(&kvapi.SetRequest{Key: key, Value: val, OldVersion: getResp.Version, NewVersion: lastSetVersion})
				if err == nil {
					return nil
				}
				return ErrorHandler(getResp, err, c, key, updateFn)
			}
		}
		return ErrorHandler(getResp, err, c, key, updateFn)
	} else {
		if err != nil {
			return ErrorHandler(getResp, err, c, key, updateFn)
		}
		val, updErr := updateFn(&getResp.Value)
		if updErr != nil {
			return updErr
		}
		lastSetVersion = uuid.Must(uuid.NewV4())
		_, err = c.Set(&kvapi.SetRequest{Key: key, Value: val, OldVersion: getResp.Version, NewVersion: lastSetVersion})

		if err != nil {
			return ErrorHandler(getResp, err, c, key, updateFn)
		}
		return nil
	}

}

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	getResp, err := c.Get(&kvapi.GetRequest{Key: key})
	return ErrorHandler(getResp, err, c, key, updateFn)
}
