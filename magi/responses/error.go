package responses

import "github.com/Jeff-All/magi/errors"

type Error struct {
	Code  errors.Code
	Error string
}
