package auth

import (
	"time"
)

var Roles = []string{
	AdminRole,
	UserRole,
}

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type Role struct {
	// ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name string `gorm:"primary_key"`

	User User `gorm:"foreignkey:Role"`
}

func InitRoles() error {
	for _, role := range Roles {
		if err := AddRole(role); err != nil {
			return err
		}
	}
	return nil
}

func AddRole(name string) error {
	return DB.Create(&Role{Name: name}).GetError()
}
