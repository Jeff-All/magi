package handlers

import "net/http"

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func HandleError(
	next ErrorHandler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if err := next(w, r); err != nil {
			w.Write([]byte(err.Error()))
		}
	})
}
