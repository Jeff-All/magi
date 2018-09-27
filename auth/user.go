package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/mail"
	"github.com/Jeff-All/magi/output"
	"github.com/jinzhu/gorm"
	gomail "gopkg.in/gomail.v2"

	log "github.com/sirupsen/logrus"
)

type User struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name       string
	NickName   string
	GivenName  string
	FamilyName string

	Active   bool   `gorm:"default:false"`
	Email    string `gorm:"type:varchar(254)"`
	SubClaim string `gorm:"unique_index:varchar(254)`

	Role string
}

func GetUser(claims map[string]interface{}) (*User, error) {
	user := User{}
	if err := DB.Where("sub_claim = ?", claims["sub"]).First(&user).GetError(); err != nil {
		if err == gorm.ErrRecordNotFound {
			if sub, ok := claims["sub"]; !ok {
				return nil, errors.CodedError{
					Message:  "sub claim undefined",
					HTTPCode: http.StatusBadRequest,
				}
			} else {
				user.SubClaim = sub.(string)
			}
			if name, ok := claims["name"]; ok {
				user.Name = name.(string)
			}
			if nickname, ok := claims["nickname"]; ok {
				user.NickName = nickname.(string)
			}
			if familyName, ok := claims["family_name"]; ok {
				user.FamilyName = familyName.(string)
			}
			if givenName, ok := claims["given_ame"]; ok {
				user.GivenName = givenName.(string)
			}
			user.Active = false
			user.Role = UserRole
			if err = DB.Create(&user).GetError(); err != nil {
				return nil, errors.CodedError{
					Message:  "error creating user",
					Err:      err,
					HTTPCode: 500,
				}
			}
		} else {
			return nil, errors.CodedError{
				Message:  "error querying users",
				Err:      err,
				HTTPCode: 500,
			}
		}
	}
	if !user.Active {
		log.WithFields(log.Fields{
			"user": user,
			"role": user.Role,
		}).Debug("authorization: inactive user")
		return nil, errors.CodedError{
			Message:  "inactive User",
			HTTPCode: http.StatusUnauthorized,
		}
	}
	return &user, nil
}

func (u *User) EnforceRole(r *http.Request) error {
	log.WithFields(log.Fields{
		"role": u.Role,
	}).Debug("auth.Users.EnforceRole")
	if ok, err := Enforcer.EnforceSafe(u.Role, r.URL.Path, r.Method); err != nil {
		return errors.CodedError{
			Message:  "error authorizing user",
			HTTPCode: 500,
		}
	} else if !ok {
		return errors.CodedError{
			Message:  "unauthorized request",
			HTTPCode: http.StatusUnauthorized,
		}
	}
	return nil
}

func (u *User) Mail(
	message *gomail.Message,
) error {
	message.SetHeader("To", u.Email)
	message.SetHeader("From", "noreply@magi.com")
	return mail.Send(message)
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(output.User(u))
}
