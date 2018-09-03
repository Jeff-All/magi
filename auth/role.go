package auth

import (
	"time"
)

type Role struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name string
}

func AddRole(name string) error {
	return DB.Create(&Role{Name: name}).GetError()
}
