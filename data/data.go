package data

type Data interface {
	Create(value interface{}) Data
	GetError() error
	AutoMigrate(...interface{}) Data
	Where(interface{}, ...interface{}) Data
	First(interface{}) Data
	Delete(interface{}) Data
	Model(interface{}) Data
	Append(interface{}) Data
	Association(string) Association
	Close() error
}

type Association interface {
	Append(interface{}) Association
	GetError() error
	Delete(interface{}) Association
}
