package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/magi/errors"
	log "github.com/sirupsen/logrus"
)

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func HandleError(
	name string,
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
				codedError.Code = 0
				codedError.Message = "Internal Server Error."
				codedError.HTTPCode = http.StatusInternalServerError
				codedError.Err = err
			}
			log.WithFields(log.Fields{
				"code":      codedError.Code,
				"message":   codedError.Message,
				"http_code": codedError.HTTPCode,
				"error":     err.Error(),
			}).Errorf("Error executing '%s'", name)

			var errorJSON []byte
			if errorJSON, err = json.Marshal(codedError); err != nil {
				log.WithFields(log.Fields{
					"name":  name,
					"error": err.Error(),
				}).Error("Error while marshaling an error into JSON")
				w.Write([]byte("Internal Server Error"))
			}
			w.WriteHeader(codedError.HTTPCode)
			w.Write(errorJSON)
		}
	})
}
