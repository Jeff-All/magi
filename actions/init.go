package actions

import "github.com/Jeff-All/magi/models"

func AutoMigrate() {
	DB.AutoMigrate(
		&models.Request{},
		// &Agency{},
		&models.Gift{},
		&models.Tag{},
		&models.Endpoint{},
		&models.Request_HTTP{},
		&models.Response_HTTP{},
	)
}
