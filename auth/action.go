package auth

import "time"

type Action struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name string

	ResourceID uint
}
