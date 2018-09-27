package models

type Tag struct {
	Base
}

type Endpoint struct {
	Base

	Name string `gorm:"size:255"`
}

type Request_HTTP struct {
	Base

	Endpoint Endpoint
}

type Response_HTTP struct {
	Base

	Request_HTTP Request_HTTP
}
