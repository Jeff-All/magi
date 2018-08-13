package util

import (
	"io"
	"io/ioutil"

	"github.com/Jeff-All/magi/mock"
)

var IOUtil iIOUtil = _IOUtil{}

type iIOUtil interface {
	ReadAll(io.Reader) ([]byte, error)
}

type _IOUtil struct{}

func (i _IOUtil) ReadAll(
	r io.Reader,
) (
	[]byte,
	error,
) {
	return ioutil.ReadAll(r)
}

type MockIOUtil struct {
	Mock mock.Mock
}

func (i MockIOUtil) ReadAll(
	r io.Reader,
) (
	[]byte,
	error,
) {
	call := i.Mock.Call("ReadAll", r)
	return1 := call.Return[0].([]byte)
	return2 := call.Return[1]
	if return2 == nil {
		return return1, nil
	}
	return return1, return2.(error)
}
