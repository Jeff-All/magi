package data

import "github.com/Jeff-All/magi/mock"

type Mock struct {
	Mock mock.Mock
}

type MockAssociation struct {
	*Mock
}

func NewMock() Mock {
	return Mock{Mock: mock.NewMock()}
}

func (m *Mock) Unscoped() Data {
	m.Mock.Call("Unscoped")
	return m
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

func (m *Mock) AutoMigrate(value ...interface{}) Data {
	m.Mock.Call("AutoMigrate", value...)
	return m
}

func (m *Mock) Order(value interface{}) Data {
	m.Mock.Call("Order", value)
	return m
}

func (m *Mock) Where(
	statement interface{},
	value ...interface{},
) Data {
	m.Mock.Call("Where", ([]interface{}{statement, value})...)
	return m
}

func (m *Mock) First(
	value interface{},
	values ...interface{},
) Data {
	m.Mock.Call("First", value, values)
	return m
}

func (m *Mock) Last(
	value interface{},
	values ...interface{},
) Data {
	m.Mock.Call("Last", value, values)
	return m
}

func (m *Mock) Delete(value interface{}) Data {
	m.Mock.Call("Delete", value)
	return m
}

func (m *Mock) Save(value interface{}) Data {
	m.Mock.Call("Save", value)
	return m
}

func (m *Mock) Updates(value interface{}) Data {
	m.Mock.Call("Updates", value)
	return m
}

func (m *Mock) Model(value interface{}) Data {
	m.Mock.Call("Model", value)
	return m
}

func (m *Mock) Find(value interface{}) Data {
	m.Mock.Call("Find", value)
	return m
}

func (m *Mock) Select(value interface{}) Data {
	m.Mock.Call("Select", value)
	return m
}

func (m *Mock) Joins(value string) Data {
	m.Mock.Call("Joins", value)
	return m
}

func (m *Mock) Append(value interface{}) Data {
	m.Mock.Call("Append", value)
	return m
}

func (m *Mock) Preload(column string, conditions ...interface{}) Data {
	m.Mock.Call("Preload", []interface{}{column, conditions}...)
	return m
}

func (m *Mock) Related(value interface{}) Data {
	m.Mock.Call("Related", value)
	return m
}

func (m *Mock) Close() error {
	call := m.Mock.Call("Close")
	if call.Return[0] == nil {
		return nil
	}
	return call.Return[0].(error)
}

func (m *Mock) Offset(value int) Data {
	m.Mock.Call("Offset", value)
	return m
}

func (m *Mock) Limit(value int) Data {
	m.Mock.Call("Limit", value)
	return m
}

func (m *Mock) Association(value string) Association {
	m.Mock.Call("Association", value)
	return &MockAssociation{Mock: m}
}

func (m *MockAssociation) Append(value interface{}) Association {
	m.Mock.Mock.Call("Association", value)
	return m
}

func (m *MockAssociation) Delete(value interface{}) Association {
	m.Mock.Mock.Call("Delete", value)
	return m
}

func (m *MockAssociation) Find(value interface{}) Association {
	m.Mock.Mock.Call("Find", value)
	return m
}
