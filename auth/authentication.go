package auth

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/Jeff-All/magi/data"

	log "github.com/sirupsen/logrus"
)

type Level int

const (
	Root = iota
	Admin
)

func Init() error {
	err := DB.AutoMigrate(&User{}, &Group{}).GetError()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error Migrating Auth tables")
		// return err
	}
	// _, err = AddRootUser(pw)

	return err
}

func AuthRequest(r *http.Request) (*User, error) {
	un, pw, _ := r.BasicAuth()
	return BasicAuthentication(un, pw)
}

func BasicAuthentication(
	un string,
	pw string,
) (*User, error) {
	var user User
	err := DB.Where("user_name = ?", un).Preload("Groups").First(&user).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"Username": un,
			"Error":    err.Error(),
		}).Error("Failed to query users")
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		log.Debug("Unable to find User")
		return nil, fmt.Errorf("Unable to find user")
	}

	pwHash := GeneratePasswordHash(pw)

	if pwHash != user.Password {
		log.WithFields(log.Fields{
			"Username": un,
		}).Debug("Password Did Not Match")
		return nil, fmt.Errorf("Password did not match")
	}

	return &user, nil
}

func GeneratePasswordHash(pw string) string {
	toReturn := sha512.Sum512([]byte(pw + PrivateKey))
	return string(toReturn[:64])
}

func AddRootUser(
	pw string,
) (*User, error) {
	log.Debug("auth.AddRootUser()")
	if !ValidatePassword(pw) {
		log.Info("Invalid Password")
		return nil, nil
	}
	// Check if Users has a root
	var user User
	err := DB.Where("user_name = ?", "root").First(&user).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error checking users table for Root")
		return nil, err
	} else if err == nil {
		log.Error("Root user already exists")
		return nil, nil
	}

	var rootGroup Group
	err = DB.Where("name = ?", "root").First(&rootGroup).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error checking group table for Root")
		return nil, err
	} else if err == nil {
		log.Error("Root group already exists")
		return nil, nil
	}

	// testModel := models.Request{}

	// err = DB.Create(&testModel).GetError()
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"Error": err,
	// 		// "group": string(groupstring),
	// 	}).Error("Error Creating temp")
	// }

	rootGroup = Group{Name: "root"}
	db, ok := DB.(*data.Gorm)
	if !ok {
		log.Error("Unable to convert to gorm")
	}
	// db, ok := data.Gorm.(DB)
	groupstring, _ := json.Marshal(rootGroup)
	err = db.DB.Create(&rootGroup).Error
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"group": string(groupstring),
		}).Error("Error Creating Root Group")
	}

	// Create Root entry in users table
	rootUser := User{
		UserName: "root",
		Password: GeneratePasswordHash(pw),
		Level:    int(Root),
		Groups:   []Group{rootGroup},
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
	// Return Root User
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
