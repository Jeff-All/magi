package application

import (
	"net/http"

	"github.com/Jeff-All/magi/auth"
	"github.com/Jeff-All/magi/endpoints"
	"github.com/Jeff-All/magi/errors"
	"github.com/gorilla/mux"
)

type validation struct {
	Email string `validate:"required,email"`
	Role  string `validate:"required,oneof=recorder shopper manager"`
}

func PUT() func(
	w http.ResponseWriter,
	r *http.Request,
) error {
	unmarshalVerify := endpoints.UnmarshalVerify(validation{})
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) error {
		model, err := unmarshalVerify(r)
		if err != nil {
			return err
		}
		input, ok := model.(validation)
		if !ok {
			return errors.CodedError{
				Message:  "error asserting type",
				HTTPCode: 500,
			}
		}
		app, err := auth.CreateApplication(input.Email, input.Role)
		if err != nil {
			return errors.CodedError{
				Message:  "error creating application",
				HTTPCode: 500,
				Err:      err,
			}
		}
		return endpoints.MarshalWrite(w, app)
	}
}

func GET(
	filename string,
	update bool,
) (
	func(w http.ResponseWriter, r *http.Request) error,
	error,
) {
	getHTML, err := endpoints.GetHTML(
		filename,
		update,
	)
	if err != nil {
		return nil, errors.CodedError{
			Message:  "unable to build GetHTML in application.GET",
			HTTPCode: 500,
		}
	}
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) error {
		app := auth.Application{Hash: mux.Vars(r)["hash"]}
		if err := app.GetByHash(); err != nil {
			return errors.CodedError{
				Message:  "unable to retrieve application",
				HTTPCode: 500,
				Err:      err,
			}
		}
		return getHTML(w, r)
	}, nil
}
