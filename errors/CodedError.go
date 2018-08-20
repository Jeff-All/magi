package errors

import (
	"github.com/sirupsen/logrus"
)

type CodedError struct {
	Message  string
	Code     int
	HTTPCode int
	Err      error `json:"-"`
	Fields   logrus.Fields
}

func (err CodedError) Error() string {
	if err.Err != nil {
		return err.Message + ": " + err.Err.Error()
	}
	return err.Message
}

func (err CodedError) Root() CodedError {
	if err.Err != nil {
		if codedError, ok := err.Err.(CodedError); ok {
			return codedError
		}
	}
	return err
}
