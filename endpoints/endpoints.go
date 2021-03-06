package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/Jeff-All/magi/data"
	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"gopkg.in/go-playground/validator.v9"
)

func GETPage(
	model interface{},
	DB data.Data,
	preload ...string,
) func(
	w http.ResponseWriter,
	r *http.Request,
) error {
	modelType := reflect.TypeOf(model)
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) error {
		params := r.URL.Query()

		var limit int
		var err error
		if limit, err = strconv.Atoi(params.Get("limit")); err != nil || limit == 0 {
			limit = 20
		}

		var offset int
		if offset, err = strconv.Atoi(params.Get("offset")); err != nil {
			offset = 0
		}

		arrayBase := reflect.MakeSlice(reflect.SliceOf(modelType), 0, limit)

		array := reflect.New(arrayBase.Type())
		array.Elem().Set(arrayBase)

		db := DB.Limit(limit).Offset(offset)
		for _, cur := range preload {
			db = db.Preload(cur)
		}
		if err := db.Find(array.Interface()).GetError(); err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.CodedError{
					Message:  "unable to find records",
					HTTPCode: 500,
					Err:      err,
				}
			}
			return errors.CodedError{
				Message:  "error querying array",
				HTTPCode: 500,
				Err:      err,
			}
		}
		logrus.Debug(array)
		responseBody, err := json.Marshal(array.Interface())
		if err != nil {
			return errors.CodedError{
				Message:  "error marshalling response",
				HTTPCode: 500,
				Err:      err,
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseBody)
		if err != nil {
			return errors.CodedError{
				Message:  "error writing body",
				HTTPCode: 500,
				Err:      err,
			}
		}
		w.WriteHeader(http.StatusFound)
		return nil
	}
}

func PUT(
	model interface{},
) func(
	w http.ResponseWriter,
	r *http.Request,
) error {
	modelType := reflect.TypeOf(model)
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) error {
		body, err := util.IOUtil.ReadAll(r.Body)
		if err != nil {
			return errors.CodedError{
				Message:  "unable to read the body",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
			}
		}
		value := reflect.New(modelType)
		if err = util.Json.Unmarshal(body, &value); err != nil {
			return errors.CodedError{
				Message:  "unable to unmarshal",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
				Fields: logrus.Fields{
					"body": string(body),
				},
			}
		}
		return nil
	}
}

func MarshalWrite(
	w http.ResponseWriter,
	obj interface{},
) error {
	bytes, err := util.Json.Marshal(obj)
	if err != nil {
		return errors.CodedError{
			Message:  "error marshaling output",
			HTTPCode: 500,
			Err:      err,
		}
	}
	if _, err = w.Write(bytes); err != nil {
		return errors.CodedError{
			Message:  "error writing output",
			HTTPCode: 500,
			Err:      err,
		}
	}
	return nil
}

func UnmarshalVerify(
	model interface{},
) func(
	r *http.Request,
) (
	interface{},
	error,
) {
	modelType := reflect.TypeOf(model)
	logrus.Debug("UnmarshalVerify:" + modelType.Name())
	valid := validator.New()
	return func(
		r *http.Request,
	) (
		interface{},
		error,
	) {
		body, err := util.IOUtil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.CodedError{
				Message:  "unable to read the body",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
			}
		}
		value := reflect.New(modelType)
		if err = util.Json.Unmarshal(body, value.Interface()); err != nil {
			return nil, errors.CodedError{
				Message:  "unable to unmarshal",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
				Fields: logrus.Fields{
					"body": string(body),
				},
			}
		}
		err = valid.Struct(value)
		if err != nil {
			return nil, errors.CodedError{
				Message:  "failed validation",
				HTTPCode: http.StatusUnprocessableEntity,
				Err:      err,
			}
		}
		logrus.WithFields(logrus.Fields{
			"body":  string(body),
			"value": value,
		}).Debug("UnmarshalVerify")
		return value.Elem().Interface(), nil
	}
}

func GetResource(w http.ResponseWriter, r *http.Request) error {
	file, err := ioutil.ReadFile("." + r.URL.Path)
	if err != nil {
		return errors.CodedError{
			Message:  "Internal Server Error",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	_, err = w.Write(file)
	if err != nil {
		return errors.CodedError{
			Message:  "Internal Server Error",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	return nil
}

func GetHTML(
	filename string,
	update bool,
) (
	func(w http.ResponseWriter, r *http.Request) error,
	error,
) {
	var file []byte
	filename = "./" + filename + ".html"
	if !update {
		var err error
		if file, err = ioutil.ReadFile(filename); err != nil {
			return nil, errors.CodedError{
				Message:  "Internal Server Error",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
			}
		}
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		if update {
			file, err = ioutil.ReadFile(filename)
			if err != nil {
				return errors.CodedError{
					Message:  "Internal Server Error",
					HTTPCode: http.StatusInternalServerError,
					Err:      err,
				}
			}
		}
		_, err = w.Write(file)
		if err != nil {
			return errors.CodedError{
				Message:  "Internal Server Error",
				HTTPCode: http.StatusInternalServerError,
				Err:      err,
			}
		}
		return nil
	}, nil
}
