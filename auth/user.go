package auth

import (
	"fmt"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type User struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time
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

func (u *User) Roles() []string {
	toReturn := make([]string, len(u.Groups))
	for _, curGroup := range u.Groups {
		toReturn = append(toReturn, curGroup.Name)
	}
	return toReturn
}
