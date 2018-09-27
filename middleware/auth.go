package middleware

import (
	"net/http"

	"github.com/casbin/casbin"

	"github.com/Jeff-All/magi/session"
)

// func AuthJWT() func(next ErrorHandler) ErrorHandler {
// 	return func(
// 		w http.ResponseWriter,
// 		r *http.Request,
// 	) error {

// 	}
// }

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
			return nil
			// auth.GetUserBySubClaim(r.Context().value("user"))
		}
	}
}
