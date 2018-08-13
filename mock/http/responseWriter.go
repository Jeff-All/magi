package http

import (
	"github.com/Jeff-All/magi/mock"

	_http "net/http"
)

type ResponseWriter struct {
	Mock mock.Mock
}

func (w *ResponseWriter) Header() _http.Header {
	call := w.Mock.Call("Header")

	return call.Return[0].(_http.Header)
}

func (w *ResponseWriter) Write(input []byte) (int, error) {
	call := w.Mock.Call("Writer", input)

	return1 := call.Return[0].(int)
	return2 := call.Return[1]
	if return2 == nil {
		return return1, nil
	}
	return return1, return2.(error)
}

func (w *ResponseWriter) WriteHeader(statuscode int) {
	w.Mock.Call("WriterHeader", statuscode)
}
