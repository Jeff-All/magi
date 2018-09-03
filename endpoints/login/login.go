package login

import (
	"net/http"

	"github.com/Jeff-All/magi/auth"
	"github.com/gorilla/mux"

	"github.com/Jeff-All/magi/endpoints"
	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/resources"
)

type validation struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=20"`
}

func PUTHash() func(
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
		app := auth.Application{Hash: mux.Vars(r)["hash"]}
		if err := app.GetByHash(); err != nil {
			return errors.CodedError{
				Message:  "unable to retrieve application",
				HTTPCode: 500,
				Err:      err,
			}
		}
		if err := app.Activate(input.Email, input.Password); err != nil {
			return errors.CodedError{
				Message:  "unable to activate account",
				HTTPCode: 500,
				Err:      err,
			}
		}
		return resources.Session.LogOut(w, r)
	}
}
