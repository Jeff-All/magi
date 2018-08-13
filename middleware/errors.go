package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Jeff-All/magi/responses"
	log "github.com/sirupsen/logrus"
)

type ErrorableHandler func(
	w http.ResponseWriter,
	r *http.Request,
) error

type ErrorHandler struct {
	Endpoint string
	ErrorableHandler
}

func (eh ErrorHandler) HandleErrors(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := eh.ErrorableHandler(w, r)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": eh.Endpoint,
			"error":    err.Error(),
		}).Error("error executing endpoint")
	}

	response, _ := json.Marshal(
		responses.Error{
			// Code:  errors.Default,
			Error: err.Error(),
		},
	)

	w.Write(response)
	return
}
