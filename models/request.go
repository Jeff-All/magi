package models

import (
	"fmt"
)

var Requests iRequests = _Requests{}

type iRequests interface {
	Create(*Request) error
}

type _Requests struct{}

type Request struct {
	BaseModel

	Agency *Agency
}

func (requests _Requests) Create(request *Request) error {
	// ToDo: Access Control Check

	if request.CreatedAt != nil {
		return fmt.Errorf("Request was already created")
	}

	if request == nil {
		return fmt.Errorf("request was nil")
	}
	err := DB.Create(request)
	return err.GetError()
}