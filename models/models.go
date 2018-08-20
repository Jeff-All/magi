package models

import (
	"time"

	"github.com/Jeff-All/magi/data"
	res "github.com/Jeff-All/magi/resources"
)

var DB data.Data

func AutoMigrate() {
	res.DB.AutoMigrate(
		&Request{},
		// &Agency{},
		&Gift{},
		&Tag{},
		&Endpoint{},
		&Request_HTTP{},
		&Response_HTTP{},
	)
}

type BaseModel struct {
	CreatedAt *time.Time `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type Gift struct {
	BaseModel

	ID          uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	Description string `gorm:"size:255"`

	RequestID uint64
}

type Tag struct {
	BaseModel
}

type Endpoint struct {
	BaseModel

	Name string `gorm:"size:255"`
}

type Request_HTTP struct {
	BaseModel

	Endpoint Endpoint
}

type Response_HTTP struct {
	BaseModel

	Request_HTTP Request_HTTP
}
