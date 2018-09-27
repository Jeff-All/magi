package session

import (
	"net/http"

	"github.com/Jeff-All/magi/errors"
	"github.com/gorilla/sessions"
)

type Manager struct {
	sessions.Store
}

func (m *Manager) LogOut(
	w http.ResponseWriter,
	r *http.Request,
) error {
	session, err := m.Store.Get(r, "user")
	if err != nil {
		return errors.CodedError{
			Message:  "unable to logout",
			HTTPCode: 500,
			Err:      err,
		}
	}

	session.Values["id"] = ""
	session.Values["roles"] = ""
	session.Options.MaxAge = -1

	if err := m.Store.Save(r, w, session); err != nil {
		return errors.CodedError{
			Message:  "unable to store session",
			HTTPCode: 500,
			Err:      err,
		}
	}

	if r.Method == "GET" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	return nil
}

// func (m *Manager) Login(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) error {
// 	session, err := m.Store.New(r, "user")
// 	if err != nil {
// 		return errors.CodedError{
// 			Message:  "Unable to load sesssion",
// 			HTTPCode: http.StatusInternalServerError,
// 			Err:      err,
// 		}
// 	}
// 	user, err := auth.AuthRequest(r)
// 	if err != nil {
// 		return errors.CodedError{
// 			Message:  "Unable to auth request",
// 			HTTPCode: http.StatusUnauthorized,
// 			Err:      err,
// 		}
// 	}
// 	if user == nil {
// 		return errors.CodedError{
// 			Message:  "Unable to auth request",
// 			HTTPCode: http.StatusUnauthorized,
// 		}
// 	}
// 	session.Values = make(map[interface{}]interface{})
// 	session.Values["id"] = user.ID
// 	session.Values["roles"] = user.GetRoles()
// 	err = m.Store.Save(r, w, session)
// 	if err != nil {
// 		return errors.CodedError{
// 			Message:  "Unable to save session",
// 			HTTPCode: http.StatusInternalServerError,
// 			Err:      err,
// 		}
// 	}
// 	return nil
// }

func (m *Manager) Load(
	r *http.Request,
) (
	*sessions.Session,
	error,
) {
	session, err := m.Store.Get(r, "user")
	if err != nil {
		return nil, errors.CodedError{
			Message:  "Unable to load sesssion",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	if !session.IsNew {
		return session, nil
	}
	return nil, nil
}
