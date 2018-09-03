package data

import (
	"github.com/jinzhu/gorm"
)

type Gorm struct {
	*gorm.DB
}

type GormAssociation struct {
	*gorm.Association
}

func (g *Gorm) Unscoped() Data {
	db := g.DB.Unscoped()
	return &Gorm{DB: db}
}

func (g *Gorm) Create(value interface{}) Data {
	db := g.DB.Create(value)
	return &Gorm{DB: db}
}

func (g *Gorm) GetError() error {
	return g.Error
}

func (g *Gorm) AutoMigrate(data ...interface{}) Data {
	// g.DB = g.DB.AutoMigrate(data...)
	// return g

	db := g.DB.AutoMigrate(data...)
	return &Gorm{DB: db}
}

func (g *Gorm) Where(
	statement interface{},
	value ...interface{},
) Data {
	// g.DB = g.DB.Where(statement, value...)
	// return g

	db := g.DB.Where(statement, value...)
	return &Gorm{DB: db}
}

func (g *Gorm) First(
	value interface{},
	values ...interface{},
) Data {
	db := g.DB.First(value, values...)
	return &Gorm{DB: db}
}

func (g *Gorm) Last(
	value interface{},
	values ...interface{},
) Data {
	db := g.DB.Last(value, values...)
	return &Gorm{DB: db}
}

func (g *Gorm) Delete(value interface{}) Data {
	db := g.DB.Delete(value)
	return &Gorm{DB: db}
}

func (g *Gorm) Save(value interface{}) Data {
	db := g.DB.Save(value)
	return &Gorm{DB: db}
}

func (g *Gorm) Model(value interface{}) Data {
	db := g.DB.Model(value)
	return &Gorm{DB: db}
}

func (g *Gorm) Find(value interface{}) Data {
	db := g.DB.Find(value)
	return &Gorm{DB: db}
}

func (g *Gorm) Preload(column string, conditions ...interface{}) Data {
	db := g.DB.Preload(column, conditions...)
	return &Gorm{DB: db}
}

func (g *Gorm) Related(value interface{}) Data {
	db := g.DB.Related(value)
	return &Gorm{DB: db}
}

func (g *Gorm) Association(value string) Association {
	return &GormAssociation{
		Association: g.DB.Association(value),
	}
}

func (g *Gorm) Close() error {
	return g.DB.Close()
}

func (g *Gorm) Offset(value int) Data {
	return &Gorm{DB: g.DB.Offset(value)}
}

func (g *Gorm) Limit(value int) Data {
	return &Gorm{DB: g.DB.Limit(value)}
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

func (g *GormAssociation) Find(value interface{}) Association {
	g.Association = g.Association.Find(value)
	return g
}
