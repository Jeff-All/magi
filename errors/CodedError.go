package errors

import (
	"github.com/sirupsen/logrus"
)

type CodedError struct {
	Public   string
	Message  string
	HTTPCode int

	Package  string
	Struct   string
	Function string
	Code     int

	Err    error `json:"-"`
	Fields logrus.Fields
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
			return codedError.Root()
		}
	}
	return err
}
