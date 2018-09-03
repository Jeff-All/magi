package data

type Data interface {
	Unscoped() Data
	Create(interface{}) Data
	GetError() error
	AutoMigrate(...interface{}) Data
	Where(interface{}, ...interface{}) Data
	First(interface{}, ...interface{}) Data
	Last(interface{}, ...interface{}) Data
	Save(interface{}) Data
	Delete(interface{}) Data
	Model(interface{}) Data
	Find(interface{}) Data
	Preload(string, ...interface{}) Data
	Related(interface{}) Data
	// Append(interface{}) Data
	Association(string) Association
	Close() error
	Offset(int) Data
	Limit(int) Data
}

type Association interface {
	Append(interface{}) Association
	GetError() error
	Delete(interface{}) Association
	Find(interface{}) Association
}
