package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func Log(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"route":  r.URL.Path,
		}).Debug("request")
		next.ServeHTTP(w, r)
	})
}
