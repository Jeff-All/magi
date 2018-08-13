package mock

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type Function struct {
	Name      string
	Arguments []interface{}
	Return    []interface{}
}

type Mock struct {
	FunctionCalls []*Function
	Errors        []error
	i             int
	j             int
}

func NewMock() Mock {
	return Mock{FunctionCalls: []*Function{}}
}

func (m *Mock) AssertEnd(
	t *testing.T,
) {
	if len(m.FunctionCalls) != m.j {
		t.Log("calls remain unasserted")
		t.Fail()
	}
}

func (m *Mock) AddCall(
	name string,
	returns ...interface{},
) {
	m.FunctionCalls = append(
		m.FunctionCalls,
		&Function{
			Name:   name,
			Return: returns,
		},
	)
}

func (m *Mock) AssertCall(
	t *testing.T,
	name string,
	arguments ...interface{},
) {
	if m.Errors != nil {
		for _, curError := range m.Errors {
			t.Log(curError.Error())
		}
		t.Fail()
	}

	if m.j >= len(m.FunctionCalls) {
		t.Log("no calls left to assert")
		t.Fail()
		return
	}
	call := m.FunctionCalls[m.j]
	m.j++

	convey.So(name, convey.ShouldEqual, call.Name)
	for i, curArg := range arguments {
		convey.So(curArg, convey.ShouldEqual, call.Arguments[i])
	}
}

func (m *Mock) Call(name string, arguments ...interface{}) *Function {
	if m.i < len(m.FunctionCalls) {
		call := m.FunctionCalls[m.i]
		call.Name = name
		call.Arguments = arguments
		m.i++
		return call
	}
	m.Errors = append(m.Errors, fmt.Errorf("Unexpected Function Call"))
	return nil
}
