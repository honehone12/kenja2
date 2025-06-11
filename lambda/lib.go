package lambda

import "errors"

const REQUEST_BODY_LIMIT = 1024

func CheckBodyLimit(req []byte) error {
	if len(req) > REQUEST_BODY_LIMIT {
		return errors.New("request body size over limit")
	}
	return nil
}
