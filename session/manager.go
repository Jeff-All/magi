package session

import (
	"fmt"
	"net/http"

	"github.com/Jeff-All/magi/auth"
	"github.com/gorilla/sessions"
)

type Manager struct {
	sessions.Store
}

func (m *Manager) Load(r *http.Request) (*sessions.Session, error) {
	session, err := m.Store.Get(r, "user")
	if err != nil {
		return nil, err
	}
	if !session.IsNew {
		return session, nil
	}
	user, err := auth.AuthRequest(r)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("Couldnt locate user")
	}
	session.Values = make(map[interface{}]interface{})
	session.Values["id"] = user.ID
	session.Values["roles"] = user.Roles()
	return session, nil
}
