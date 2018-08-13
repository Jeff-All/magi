package data

import "github.com/jinzhu/gorm"

type Gorm gorm.DB

func (g *Gorm) Create(value interface{}) Data {
	return g.Create(value)
}

func (g *Gorm) GetError() error {
	return g.Error
}
