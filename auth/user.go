package auth

import (
	res "github.com/Jeff-All/magi/resources"

	"github.com/Jeff-All/magi/models"

	log "github.com/sirupsen/logrus"
)

type User struct {
	models.User
}

func (u User) AddUser(
	un string,
	pw string,
	level Level,
) (*User, error) {
	log.WithFields(log.Fields{
		"User.Level": u.Level,
		"Root.Level": Root,
	}).Debug("auth.User.AddUser()")
	err := ValidatePassword(pw)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Password")
		return nil, err
	}
	switch u.Level {
	case Root, Admin:
		log.Debug("Root/Admin")
		if level > Level(u.Level) {
			newUser := User{User: models.User{
				UserName: un,
				Password: GeneratePasswordHash(pw),
				Level:    int(level),
			}}
			err := res.DB.Create(&newUser).Error
			if err != nil {
				log.WithFields(log.Fields{
					"Error":    err.Error(),
					"UserName": un,
					"Level":    level,
				}).Error("Failed to Add User")
				return nil, err
			}
			return &newUser, nil
		} else {
			log.WithFields(log.Fields{
				"Creator":       u.UserName,
				"Creator Level": u.Level,
				"New Level":     level,
			}).Debug("Creator level must exceed new level")
			return nil, nil
		}
	default:
		log.Debug("Need at least Admin level to add users")
		break
	}
	log.Debug("Reached the end")
	return nil, nil
}

func (u User) AddUserToGroup(
	user User,
	group string,
) {

}
