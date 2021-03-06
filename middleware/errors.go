package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/magi/errors"
	log "github.com/sirupsen/logrus"
)

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func HandleError(
	next ErrorHandler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if err := next(w, r); err != nil {
			var ok bool
			var codedError errors.CodedError
			if codedError, ok = err.(errors.CodedError); !ok {
				log.Debug("Not a Coded Error")
				codedError.Code = 0
				codedError.Message = "Internal Server Error."
				codedError.HTTPCode = http.StatusInternalServerError
				codedError.Err = err
			}
			root := codedError.Root()
			fields := log.Fields{
				"message":   codedError.Message,
				"http_code": codedError.HTTPCode,
				"code":      codedError.Code,
				"error":     codedError.Err,
				"package":   root.Package,
				"struct":    root.Struct,
				"function":  root.Function,
				"err":       root.Err,
			}
			for key, cur := range root.Fields {
				if _, ok := fields[key]; !ok {
					fields[key] = cur
				} else {
					fields[":"+key] = cur
				}
			}
			log.WithFields(fields).Errorf(
				"Error executing '%s':'%s'",
				r.URL.Path,
				r.Method,
			)

			var errorJSON []byte
			if errorJSON, err = json.Marshal(root); err != nil {
				log.WithFields(log.Fields{
					"route":  r.URL.Path,
					"method": r.Method,
					"error":  err.Error(),
				}).Error("Error while marshaling an error into JSON")
				w.Write([]byte("Internal Server Error"))
			}
			w.WriteHeader(root.HTTPCode)
			w.Write(errorJSON)
		}
	})
}
