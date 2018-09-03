package auth

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Jeff-All/magi/errors"
	gomail "gopkg.in/gomail.v2"
)

type Application struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	UpdatedAt time.Time
	CreatedAt *time.Time
	DeletedAt *time.Time `sql:"index"`

	Hash string `gorm:"unique_index;type:varchar(64)"`

	User   *User
	UserID uint
}

func GetNextID() (
	uint64,
	error,
) {
	var app Application
	err := DB.Unscoped().Last(&app).GetError()
	if err != nil && err == gorm.ErrRecordNotFound {
		return 0, nil
	} else if err != nil {
		return 0, errors.CodedError{
			Message:  "error querying last id",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	return app.ID + 1, nil
}

func GetNextHash() (
	string,
	error,
) {
	nextID, err := GetNextID()
	if err != nil {
		return "", errors.CodedError{
			Message:  "error getting next id",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	hashArray := sha512.Sum512([]byte(fmt.Sprintf("%s%d%s", PrivateKey, nextID, Salt)))
	return base64.RawURLEncoding.EncodeToString(hashArray[:64]), nil
}

func CreateApplication(
	email string,
	role string,
) (
	*Application,
	error,
) {
	roleObj := &Role{}
	if err := DB.Where("name = ? AND name != 'root'", role).First(roleObj).GetError(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.CodedError{
				Message:  "invalid role",
				HTTPCode: http.StatusUnprocessableEntity,
				Err:      err,
				Fields: logrus.Fields{
					"role": role,
				},
			}
		}
		return nil, errors.CodedError{
			Message:  "error querying roles",
			HTTPCode: 500,
			Err:      err,
		}
	}
	application := Application{
		User: &User{
			Email:  email,
			Active: false,
			Roles:  []Role{*roleObj},
		},
	}

	var err error
	if application.Hash, err = GetNextHash(); err != nil {
		return nil, errors.CodedError{
			Message:  "error generating next hash",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	err = DB.Create(application.User).Create(&application).GetError()
	if err != nil {
		return nil, errors.CodedError{
			Message:  "error writing application",
			HTTPCode: 500,
			Err:      err,
		}
	}
	if err = application.Email(); err != nil {
		return nil, errors.CodedError{
			Message:  "error sending email",
			HTTPCode: 500,
			Err:      err,
		}
	}
	return &application, nil
}

func (app *Application) Email() error {
	message := gomail.NewMessage()
	message.SetHeader("Subject", "Magi Account Creation")
	message.SetBody(
		"text/html",
		fmt.Sprintf(
			`<html><body><p>Please follow the <a href="https://%s/login/%s">link</a> to register your account.</p></body></html>`,
			viper.GetString("server.domain"),
			app.Hash,
		),
	)
	return app.User.Mail(message)
}

func (app *Application) Activate(
	email string,
	password string,
) error {
	var err error
	if app.User == nil {
		app.User = &User{}
		err = DB.Model(app).Related(app.User).GetError()
		// err = DB.Model(app).Association("users").Find(app.User).GetError()
		// m := DB.Model(app)
		// a := m.Association("user")
		// f := a.Find(app.User)
		// err = f.GetError()
	}
	if err != nil || app.User == nil {
		return errors.CodedError{
			Message:  "error querying user",
			HTTPCode: 500,
			Err:      err,
		}
	}

	if email != app.User.Email {
		return errors.CodedError{
			Message:  "invalid email",
			HTTPCode: http.StatusUnprocessableEntity,
		}
	}

	app.User.Password = GeneratePasswordHash(password)
	app.User.Active = false

	err = DB.Save(app.User).GetError()
	if err != nil {
		return errors.CodedError{
			Message:  "error updating user",
			HTTPCode: 500,
			Err:      err,
		}
	}

	err = DB.Delete(app).GetError()
	if err != nil {
		return errors.CodedError{
			Message:  "error deleting application",
			HTTPCode: 500,
			Err:      err,
		}
	}

	app.User.Active = true

	err = DB.Save(app.User).GetError()
	if err != nil {
		return errors.CodedError{
			Message:  "error updating user",
			HTTPCode: 500,
			Err:      err,
		}
	}
	return nil
}

func (app *Application) GetByHash() error {
	if err := DB.Where("hash = ?", app.Hash).First(app).GetError(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.CodedError{
				Message:  "Unable to locate User Application",
				HTTPCode: http.StatusNoContent,
				Package:  "auth",
				Struct:   "Application",
				Function: "GetByHash",
				Fields: logrus.Fields{
					"hash": app.Hash,
				},
			}
		} else {
			return errors.CodedError{
				Message:  "error looking up application",
				HTTPCode: 500,
				Package:  "auth",
				Struct:   "Application",
				Function: "GetByHash",
				Err:      err,
			}
		}
	}
	return nil
}
