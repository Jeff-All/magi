package mail

import (
	"github.com/Jeff-All/magi/errors"
	"gopkg.in/gomail.v2"
)

var dialer gomail.Dialer

func Init(
	host string,
	port int,
	user string,
	pass string,
) {
	dialer = gomail.Dialer{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
	}
}

func Send(
	message *gomail.Message,
) error {
	if err := dialer.DialAndSend(message); err != nil {
		return errors.CodedError{
			Message:  "error sending e-mail",
			HTTPCode: 500,
			Code:     0,
			Err:      err,
		}
	}
	return nil
}
