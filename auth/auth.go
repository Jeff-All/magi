package auth

import (
	log "github.com/sirupsen/logrus"
)

func Init() error {
	err := DB.AutoMigrate(&User{}, &Role{}).GetError()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error Migrating Auth tables")
		return err
	}
	return InitRoles()
}
