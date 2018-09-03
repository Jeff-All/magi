package auth

import (
	"crypto/sha512"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/Jeff-All/magi/data"

	"github.com/Jeff-All/magi/errors"
	log "github.com/sirupsen/logrus"
)

type Level int

const (
	Root = iota
	Admin
)

func Init() error {
	err := DB.AutoMigrate(&User{}, &Role{}, &Application{}).GetError()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error Migrating Auth tables")
	}

	return err
}

func AuthRequest(r *http.Request) (*User, error) {
	un, pw, _ := r.BasicAuth()
	return BasicAuthentication(un, pw)
}

func BasicAuthentication(
	username string,
	password string,
) (*User, error) {
	var user User
	err := DB.Where("email = ? AND active = 1", username).Preload("Roles").First(&user).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.CodedError{
			Message:  "failed to query users",
			HTTPCode: 500,
			Code:     0,
			Err:      err,
		}
	} else if err == gorm.ErrRecordNotFound {
		return nil, errors.CodedError{
			Message:  "invalid authentication",
			HTTPCode: http.StatusUnauthorized,
			Code:     1,
		}
	}
	pwHash := GeneratePasswordHash(password)
	if pwHash != user.Password {
		return nil, errors.CodedError{
			Message:  "invalid password",
			HTTPCode: http.StatusUnauthorized,
			Code:     2,
		}
	}
	return &user, nil
}

func GeneratePasswordHash(pw string) string {
	toReturn := sha512.Sum512([]byte(Salt + pw + PrivateKey))
	return string(toReturn[:64])
}

func AddRootUser(
	password string,
) (*User, error) {
	log.Debug("auth.AddRootUser()")
	if !ValidatePassword(password) {
		log.Info("Invalid Password")
		return nil, nil
	}
	// Check if Users has a root
	var user User
	err := DB.Where("email = ?", "root").First(&user).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error checking users table for Root")
		return nil, err
	} else if err == nil {
		log.Error("Root user already exists")
		return nil, nil
	}

	var rootRole Role
	err = DB.Where("name = ?", "root").First(&rootRole).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error checking group table for Root")
		return nil, err
	} else if err == nil {
		log.Error("Root group already exists")
		return nil, nil
	}

	rootRole = Role{Name: "root"}
	db, ok := DB.(*data.Gorm)
	if !ok {
		log.Error("Unable to convert to gorm")
	}
	groupstring, _ := json.Marshal(rootRole)
	err = db.DB.Create(&rootRole).Error
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"group": string(groupstring),
		}).Error("Error Creating Root Group")
	}

	rootUser := User{
		Email:    "root",
		Password: GeneratePasswordHash(password),
		Roles:    []Role{rootRole},
		Active:   true,
	}

	userstring, _ := json.Marshal(rootUser)
	err = DB.Create(&rootUser).GetError()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"ID":    rootUser.ID,
			"user":  string(userstring),
		}).Error("Error Creating Root User")
		return nil, err
	}
	log.Debug("auth.AddRootUsers(): Success")
	return &rootUser, nil
}

func ValidatePassword(pw string) bool {
	// Check length is at least 5 characters long
	if len([]rune(pw)) < 5 {
		log.WithFields(log.Fields{
			"Password": pw,
		}).Debug("Password has to be atleast 5 characters long")
		return false
	}
	return true
}
