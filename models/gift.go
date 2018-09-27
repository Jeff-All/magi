package models

type Gift struct {
	Base

	ID uint64 `gorm:"primary_key;AUTO_INCREMENT"`

	Category string
	Name     string
	Detail   string

	RequestID uint64
}
