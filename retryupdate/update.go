//go:build !solution

package retryupdate

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var res *kvapi.GetResponse
	var err error
	var oldValue *string
	correctValue := false
	newVersion := uuid.UUID{}
	for {
		if !correctValue {
			res, err = c.Get(&kvapi.GetRequest{Key: key})
		}
		authErr := &kvapi.AuthError{}
		switch {
		case errors.Is(err, kvapi.ErrKeyNotFound):
			oldValue = nil
		case errors.As(err, &authErr):
			return err
		case err == nil:
			if !correctValue {
				oldValue = &res.Value
			}
		default:
			continue
		}
		correctValue = true
		newValue, err := updateFn(oldValue)
		if err != nil {
			return err
		}
		if (newVersion != uuid.UUID{}) {
			newVersion = uuid.Must(uuid.NewV4())
		}
		lstVersion := uuid.UUID{}
		if res != nil {
			lstVersion = res.Version
		}

		req := kvapi.SetRequest{Key: key, Value: newValue, OldVersion: lstVersion, NewVersion: newVersion}
		_, err = c.Set(&req)
		fmt.Println("$", err, newVersion, lstVersion)
		var conflictErr *kvapi.ConflictError
		fmt.Println(err)
		switch {
		case errors.As(err, &conflictErr):
			correctValue = false
			continue
		case errors.Is(err, kvapi.ErrKeyNotFound):
			oldValue = nil
			res = nil
			err = nil
			continue
		case errors.As(err, &authErr):
			fmt.Println("MEOW")
			return err
		case err == nil:
			return nil
		default:
			err = nil
			continue
		}
	}
}
