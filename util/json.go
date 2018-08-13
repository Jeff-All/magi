package util

import (
	"encoding/json"

	"github.com/Jeff-All/magi/mock"
)

var Json iJson = _Json{}

type iJson interface {
	Unmarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
}

type _Json struct{}

func (j _Json) Unmarshal(
	source []byte,
	object interface{},
) error {
	return json.Unmarshal(source, object)
}

func (j _Json) Marshal(
	object interface{},
) (
	[]byte,
	error,
) {
	return json.Marshal(object)
}

type MockJson struct {
	Mock mock.Mock
}

func (j *MockJson) Unmarshal(
	source []byte,
	object interface{},
) error {
	call := j.Mock.Call("Unmarshal", source, object)
	if call.Return[0] == nil {
		return nil
	}
	return call.Return[0].(error)
}

func (j *MockJson) Marshal(
	object interface{},
) (
	[]byte,
	error,
) {
	call := j.Mock.Call("Marshal", object)
	return1 := call.Return[0].([]byte)
	return2 := call.Return[1]
	if return2 == nil {
		return return1, nil
	}
	return return1, return2.(error)
}
