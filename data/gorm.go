package data

import "github.com/jinzhu/gorm"

type Gorm struct {
	*gorm.DB
}

type GormAssociation struct {
	*gorm.Association
}

func (g *Gorm) Create(value interface{}) Data {
	// g.DB = g.DB.Create(value)
	// return g

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

func (g *Gorm) First(value interface{}) Data {
	// g.DB = g.DB.First(value)
	// return g

	db := g.DB.First(value)
	return &Gorm{DB: db}
}

func (g *Gorm) Delete(value interface{}) Data {
	// g.DB = g.DB.Delete(value)
	// return g

	db := g.DB.Delete(value)
	return &Gorm{DB: db}
}

func (g *Gorm) Model(value interface{}) Data {
	// g.DB = g.DB.Model(value)
	// return g

	db := g.DB.Model(value)
	return &Gorm{DB: db}
}

// func (g *Gorm) Append(value interface{}) Data {
// 	// g.DB = g.DB.Model(value)
// 	// return g

// 	g.DB.
// 	db := g.DB.Append(value)
// 	return &Gorm{DB: db}
// }

func (g *Gorm) Preload(column string, conditions ...interface{}) Data {
	db := g.DB.Preload(column, conditions...)
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
