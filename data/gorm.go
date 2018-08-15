package data

import "github.com/jinzhu/gorm"

type Gorm struct {
	*gorm.DB
}

type GormAssociation struct {
	*gorm.Association
}

func (g *Gorm) Create(value interface{}) Data {
	g.DB = g.DB.Create(value)
	return g
}

func (g *Gorm) GetError() error {
	return g.Error
}

func (g *Gorm) AutoMigrate(data ...interface{}) Data {
	g.DB = g.DB.AutoMigrate(data...)
	return g
}

func (g *Gorm) Where(
	statement interface{},
	value ...interface{},
) Data {
	g.DB = g.DB.Where(statement, value...)
	return g
}

func (g *Gorm) First(value interface{}) Data {
	g.DB = g.DB.First(value)
	return g
}

func (g *Gorm) Delete(value interface{}) Data {
	g.DB = g.DB.Delete(value)
	return g
}

func (g *Gorm) Model(value interface{}) Data {
	g.DB = g.DB.Model(value)
	return g
}

func (g *Gorm) Append(value interface{}) Data {
	g.DB = g.DB.Model(value)
	return g
}

func (g *Gorm) Association(value string) Association {
	return &GormAssociation{
		Association: g.DB.Association(value),
	}
}

func (g *Gorm) Close() error {
	return g.DB.Close()
}

// Gorm Assocciation

func (g *GormAssociation) Append(value interface{}) Association {
	g.Association = g.Association.Append(value)
	return g
}

func (g *GormAssociation) GetError() error {
	return g.Association.Error
}

func (g *GormAssociation) Delete(value interface{}) Association {
	g.Association = g.Association.Delete(value)
	return g
}
