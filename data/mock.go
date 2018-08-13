package data

import "github.com/Jeff-All/magi/mock"

type Mock struct {
	Mock mock.Mock
}

func NewMock() Mock {
	return Mock{Mock: mock.NewMock()}
}

func (m *Mock) Create(value interface{}) Data {
	m.Mock.Call("Create", value)
	return m
}

func (m *Mock) GetError() error {
	call := m.Mock.Call("GetError")
	if call.Return[0] == nil {
		return nil
	}
	return call.Return[0].(error)
}
