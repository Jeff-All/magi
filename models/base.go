package models

import "time"

type Base struct {
	CreatedAt *time.Time `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
