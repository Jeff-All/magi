package errors

type CodedError struct {
	Message  string
	Code     int
	HTTPCode int
	Err      error `json:"-"`
}

func (err CodedError) Error() string {
	return err.Err.Error()
}
