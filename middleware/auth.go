package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin"

	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/session"
)

func Authorize(
	enforcer *casbin.Enforcer,
	sessionManager *session.Manager,
) func(next ErrorHandler) ErrorHandler {
	return func(next ErrorHandler) ErrorHandler {
		return func(
			w http.ResponseWriter,
			r *http.Request,
		) error {
			session, err := sessionManager.Load(r)
			if err != nil {
				return err
				// return errors.CodedError{
				// 	Message:  "Internal Server Error",
				// 	Code:     50,
				// 	HTTPCode: 500,
				// 	Err:      err,
				// }
			}
			if session == nil {
				return errors.CodedError{
					Message:  "Invalid Authentication",
					Code:     50,
					HTTPCode: 401,
					Err:      fmt.Errorf("Invalid Authentication"),
				}
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
