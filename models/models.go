package models

import (
	"time"

	res "github.com/Jeff-All/magi/resources"
)

func AutoMigrate() {
	res.DB.AutoMigrate(
		&Request{},
		&Agency{},
		&Gift{},
		&Tag{},
		&Endpoint{},
		&Request_HTTP{},
		&Response_HTTP{},
	)
}

type BaseModel struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Request struct {
	BaseModel

	Agency Agency
}

type Agency struct {
	BaseModel
}

type Gift struct {
	BaseModel

	Description string `gorm:"size:255"`

	Request Request
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
