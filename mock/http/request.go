package http

import (
	"github.com/Jeff-All/magi/mock"

	_http "net/http"
)

type Request struct {
	*_http.Request

	Mock mock.Mock
}
