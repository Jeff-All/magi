package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin"

	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/session"
	"github.com/sirupsen/logrus"
)

func Authorize(
	enforcer *casbin.Enforcer,
	sessionManager *session.Manager,
	loginURL string,
) func(next ErrorHandler) ErrorHandler {
	return func(next ErrorHandler) ErrorHandler {
		return func(
			w http.ResponseWriter,
			r *http.Request,
		) error {
			session, err := sessionManager.Load(r)
			if err != nil {
				return errors.CodedError{
					Message:  "Internal Server Error",
					HTTPCode: http.StatusInternalServerError,
					Err:      err,
				}
			}
			if session == nil {
				http.Redirect(w, r, loginURL+"?origin="+r.URL.Path, 302)
				logrus.WithFields(logrus.Fields{
					"redirect_url": loginURL,
					"route":        r.URL.Path,
				}).Debug("No session: Re-Directing to log-in")
				return nil
			}
			roles, ok := session.Values["roles"].([]string)
			if !ok {
				return errors.CodedError{
					Message:  "Internal Server Error",
					Code:     50,
					HTTPCode: 500,
					Err:      fmt.Errorf("Error pulling roles from session"),
				}
			}
			for _, curRole := range roles {
				ok, err := enforcer.EnforceSafe(curRole, r.URL.Path, r.Method)
				if err != nil {
					return errors.CodedError{
						Message:  "Internal Server Error",
						Code:     50,
						HTTPCode: 500,
						Err:      err,
					}
				}
				if ok {
					return next(w, r)
				}
			}
			return errors.CodedError{
				Message:  "User is not Authorized",
				Code:     50,
				HTTPCode: 401,
				Err:      fmt.Errorf("Invalid Authorization"),
			}
		}
	}
}
