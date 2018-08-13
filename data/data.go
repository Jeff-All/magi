package data

type Data interface {
	Create(value interface{}) Data
	GetError() error
}
