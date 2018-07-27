package auth

import (
	"fmt"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type User struct {
	ID uint64 `gorm:"primary_key;AUTO_INCREMENT"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	UserName string `gorm:"size:255;unique"`
	Password string `gorm:"size:64"`

	Level int

	Groups []Group `gorm:"many2many:user_groups"`
}

func _Root() User {
	return User{
		Level: 0,
	}
}

func ReadPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error reading password")
		return "", err
	}
	return string(password), nil
}

func (u User) Authorize(
	Category string,
	Resource string,
	Action string,
) (bool, error) {
	if u.Level == Root {
		password, _ := ReadPassword()
		if user, err := BasicAuth("root", password); user == nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
	// return false, nil
	// var returned int
	// err := res.DB.Raw(`
	// 	SELECT r.id
	// 	FROM users AS u
	// 	JOIN user_groups AS ug
	// 		ON u.id = ug.user_id
	// 	JOIN groups as g
	// 		ON g.id = ug.group_id
	// 	JOIN resource_groups as rg
	// 		ON rg.group_id = ug.group_id
	// 	JOIN resources as r
	// 		ON r.id = rg.resource_id
	// 	WHERE r.Category = ?
	// 	AND r.Resource = ?
	// 	AND r.Action = ?`,
	// 	Category, Resource, Action,
	// ).Scan(&returned).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	log.WithFields(log.Fields{
	// 		"Error": err,
	// 	}).Error("Error querying authorization")
	// 	return false, err
	// } else if err == gorm.ErrRecordNotFound {
	// 	log.Info("Unauthorized")
	// 	return false, nil
	// }
	// return true, nil
}

// func (u User) AddResource(
// 	Resource *models.Resource,
// ) (bool, error) {
// 	// Check if user has permission to add a resource to a group
// 	if ok, err := u.Authorize("auth", "resources", "add"); err != nil || !ok {
// 		return ok, err
// 	}

// 	err := res.DB.Create(&Resource).Error
// 	return err != gorm.ErrRecordNotFound && err == nil, err
// }

// func (u User) AddActionToResource(
// 	Resource models.Resource,
// 	Action models.Action,
// ) (bool, error) {
// 	// Check if user has permission to add an action to a resource
// 	if ok, err := u.Authorize("actions", Action.Name, "modify"); err != nil {
// 		return false, err
// 	} else if !ok {
// 		return false, nil
// 	}
// 	// Check if user has permission to access the resource
// 	if ok, err := u.Authorize("resources.", string(Resource.ID), "add_to_group"); err != nil {
// 		return false, err
// 	} else if !ok {
// 		return false, nil
// 	}

// 	err := res.DB.Model(&Group).Association("resources").Append(&Resource).Error
// 	return err != gorm.ErrRecordNotFound && err == nil, err
// }

// func (u User) AddGroup(
// 	Name string,
// ) (*models.Group, error) {
// 	// Check if user has permission to add a group
// 	if ok, err := u.Authorize("auth", "groups", "add"); err != nil {
// 		return nil, err
// 	} else if !ok {
// 		return nil, nil
// 	}
// 	var Group *models.Group
// 	Group.Name = Name
// 	err := res.DB.Create(Group).Error
// 	return Group, err
// }

// func (u User) AddUserToGroup(
// 	User models.User,
// 	Group models.Group,
// ) (bool, error) {
// 	// Check if user has permission to add to group
// 	if ok, err := u.Authorize("groups", Group.Name, "add_user"); err != nil {
// 		return false, err
// 	} else if !ok {
// 		return false, nil
// 	}

// 	err := res.DB.Model(&Group).Association("Users").Append(User).Error
// 	return err != gorm.ErrRecordNotFound && err == nil, err
// }

// func (u User) AddUser(
// 	un string,
// 	pw string,
// 	level Level,
// ) (*User, error) {
// 	// Check if user has permission to add to group
// 	if ok, err := u.Authorize("auth", "users", "add"); err != nil {
// 		return nil, err
// 	} else if !ok {
// 		return nil, nil
// 	}
// 	if !ValidatePassword(pw) {
// 		log.Info("Invalid Password")
// 		return nil, nil
// 	}
// 	var user *User
// 	user.User = models.User{
// 		UserName: un,
// 		Password: GeneratePasswordHash(pw),
// 	}
// 	err := res.DB.Create(&user.User).Error
// 	return user, err
// }
