package session

import (
	"net/http"

	"github.com/Jeff-All/magi/auth"
	"github.com/Jeff-All/magi/errors"
	"github.com/gorilla/sessions"
)

type Manager struct {
	sessions.Store
}

func (m *Manager) Load(r *http.Request) (*sessions.Session, error) {
	session, err := m.Store.Get(r, "user")
	if err != nil {
		return nil, errors.CodedError{
			Message:  "Unable to load sessiion",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	if !session.IsNew {
		return session, nil
	}
	user, err := auth.AuthRequest(r)
	if err != nil {
		return nil, errors.CodedError{
			Message:  "Unable to auth request",
			HTTPCode: http.StatusUnauthorized,
			Err:      err,
		}
	}
	if user == nil {
		return nil, errors.CodedError{
			Message:  "Unable to auth request",
			HTTPCode: http.StatusUnauthorized,
		}
	}
	session.Values = make(map[interface{}]interface{})
	session.Values["id"] = user.ID
	session.Values["roles"] = user.Roles()
	return session, nil
}
