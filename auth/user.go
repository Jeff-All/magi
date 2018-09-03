package auth

import (
	"syscall"
	"time"

	"github.com/Jeff-All/magi/mail"
	"golang.org/x/crypto/ssh/terminal"
	gomail "gopkg.in/gomail.v2"
)

type User struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Password string `gorm:"size:64"`
	Active   bool   `gorm:"default:false"`
	Email    string `gorm:"unique_index;type:varchar(254)"`

	Roles []Role `gorm:"many2many:user_roles"`
}

func ReadPassword() (string, error) {
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func (u *User) Mail(
	message *gomail.Message,
) error {
	message.SetHeader("To", u.Email)
	message.SetHeader("From", "noreply@magi.com")
	return mail.Send(message)
}

func (u *User) GetRoles() []string {
	toReturn := make([]string, len(u.Roles))
	for _, cur := range u.Roles {
		toReturn = append(toReturn, cur.Name)
	}
	return toReturn
}
