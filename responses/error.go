package responses

import "magi/errors"

type Error struct {
	Code  errors.Code
	Error string
}
