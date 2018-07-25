package auth

import (
	"crypto/sha512"
	"fmt"
	"net/http"

	"github.com/Jeff-All/magi/models"

	"github.com/jinzhu/gorm"

	res "github.com/Jeff-All/magi/resources"

	log "github.com/sirupsen/logrus"
)

var privateKey string

type Level int

const (
	Root = iota
	Admin
)

func AuthRequest(r *http.Request) (*User, error) {
	un, pw, _ := r.BasicAuth()
	return BasicAuth(un, pw)
}

func BasicAuth(
	un string,
	pw string,
) (*User, error) {
	user := models.User{}
	err := res.DB.Where("user_name = ?", un).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"Username": un,
			"Error":    err.Error(),
		}).Error("Failed to query users")
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		log.Debug("Unable to find User")
		return nil, nil
	}

	pwHash := GeneratePasswordHash(pw)

	if pwHash != user.Password {
		log.WithFields(log.Fields{
			"Username": un,
		}).Debug("Password Did Not Match")
		return nil, nil
	}

	return &User{User: user}, nil
}

func GeneratePasswordHash(pw string) string {
	toReturn := sha512.Sum512([]byte(pw + privateKey))
	return string(toReturn[:64])
}

func Init(pw string) (*User, error) {
	err := ValidatePassword(pw)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Password")
		return nil, err
	}
	// Check if Users has a root
	var user models.User
	err = res.DB.Where("user_name = ?", "root").First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error checking users table for Root")
		return nil, err
	} else if err == nil {
		log.Error("Root user already exists")
		return nil, nil
	}

	// Create Root entry in users table
	rootUser := User{User: models.User{
		UserName: "root",
		Password: GeneratePasswordHash(pw),
		Level:    int(Root),
	}}

	err = res.DB.Create(&rootUser).Error
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Creating Root User")
		return nil, err
	}
	// Return Root User
	return &rootUser, nil
}

func ValidatePassword(pw string) error {
	// Check length is at least 5 characters long
	if len([]rune(pw)) < 5 {
		log.WithFields(log.Fields{
			"Password": pw,
		}).Debug("Password has to be atleast 5 characters long")
		return fmt.Errorf("Password has to be at least 5 characters long")
	}
	return nil
}
