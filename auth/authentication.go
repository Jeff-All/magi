package auth

import (
	"crypto/sha512"
	"net/http"

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
	var user User
	err := res.DB.Where("user_name = ?", un).First(&user).GetError()
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

	return &user, nil
}

func GeneratePasswordHash(pw string) string {
	toReturn := sha512.Sum512([]byte(pw + privateKey))
	return string(toReturn[:64])
}

func Init(pw string) (*User, error) {
	root, err := AddRootUser(pw)
	if err != nil {
		return nil, err
	}

	// superAdmins, err := root.AddGroup("super_admins")
	// admins, err := root.AddGroup("admins")
	// users, err := root.AddGroup("users")

	// authResources, err = root.AddResource("auth", "resources")
	// authActions, err = root.AddResource("auth", "actions")
	// authGroups, err = root.AddResource("auth", "groups")
	// authUsers, err = root.AddResource("auth", "users")

	// err = authResources.AddAction("create", "read", "update", "delete")
	// err = authActions.AddAction("create", "read", "update", "delete")
	// err = authGroups.AddAction("create", "read", "update", "delete")
	// err = authUsers.AddAction("create", "read", "update", "delete")

	// err = superAdmins.AddAction(authResources.Actions("create", "read", "update", "delete")...)
	// err = superAdmins.AddAction(authActions.Actions("create", "read", "update", "delete")...)

	// root.AddResource(&models.Resource{
	// 	Category: "auth",
	// 	Resource: "",
	// })

	return root, nil
}

func AddRootUser(pw string) (*User, error) {
	if !ValidatePassword(pw) {
		return nil, nil
	}
	// Check if Users has a root
	var user User
	err := res.DB.Where("user_name = ?", "root").First(&user).GetError()
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
	rootUser := User{
		UserName: "root",
		Password: GeneratePasswordHash(pw),
		Level:    int(Root),
	}

	err = res.DB.Create(&rootUser).GetError()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Creating Root User")
		return nil, err
	}
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
