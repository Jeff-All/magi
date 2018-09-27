package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/Jeff-All/magi/data"
	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"gopkg.in/go-playground/validator.v9"
)

func ParseWhereParam(model interface{}, param string) (string, []string, error) {

	var s scanner.Scanner
	s.Init(strings.NewReader(param))
	s.Filename = "example"
	queryBuilder := &strings.Builder{}
	arguments, err := parseStatement(&s, queryBuilder, reflect.TypeOf(model))
	if err != nil {
		return "", nil, err
	}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch s.TokenText() {
		case "+":
			queryBuilder.Write([]byte(" AND "))
			break
		case "|":
			queryBuilder.Write([]byte(" OR "))
			break
		}
		if args, err := parseStatement(&s, queryBuilder, reflect.TypeOf(model)); err != nil {
			return "", nil, err
		} else {
			arguments = append(arguments, args...)
		}
	}
	return queryBuilder.String(), arguments, nil
}

func parseStatement(
	s *scanner.Scanner,
	w io.Writer,
	typ reflect.Type,
) ([]string, error) {
	if err := parseColumn(s, w, typ); err != nil {
		return nil, err
	}
	if arguments, err := parseOperator(s, w); err != nil {
		return nil, err
	} else {
		return arguments, nil
	}
}

func parseArgument(
	s *scanner.Scanner,
) string {
	toReturn := ""
	for tok := s.Peek(); tok != scanner.EOF; tok = s.Peek() {
		switch tok {
		case '+':
		case '|':
		case scanner.EOF:
			return toReturn
		}
		tok = s.Scan()
		toReturn += s.TokenText()
	}
	return toReturn
}

func parseArray(
	s *scanner.Scanner,
) ([]string, error) {
	toReturn := []string{}
	current := ""
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		log.Printf("parseArray: %v", s.TokenText())
		switch tok {
		case ']':
			return append(toReturn, current), nil
		case ',':
			toReturn = append(toReturn, current)
			current = ""
		default:
			current += s.TokenText()
		}
	}
	return nil, fmt.Errorf("parseArray: unexpected end of file")
}

func parseColumn(
	s *scanner.Scanner,
	w io.Writer,
	typ reflect.Type,
) error {
	if tok := s.Scan(); tok != scanner.EOF {
		if _, ok := typ.FieldByName(s.TokenText()); !ok {
			return fmt.Errorf("%s is not a valid column", s.TokenText())
		}
		_, err := w.Write([]byte(s.TokenText()))
		return err
	}

	return fmt.Errorf("parseColumn: unexpected end of file")
}

func parseOperator(s *scanner.Scanner, w io.Writer) ([]string, error) {
	not := false
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '!':
			log.Println("Not")
			if not {
				return nil, fmt.Errorf("unexpected character '!'")
			}
			switch s.Peek() {
			case '=':
				w.Write([]byte("!"))
				break
			case '%':
				fallthrough
			case '[':
				w.Write([]byte(" NOT"))
				break
			default:
				return nil, fmt.Errorf("unexpected character")
			}
			not = true
			break
		case '=':
			w.Write([]byte("=?"))
			return []string{parseArgument(s)}, nil
		case '%':
			w.Write([]byte(" LIKE ?"))
			return []string{"%" + parseArgument(s)}, nil
		case '[':
			w.Write([]byte(" IN(?"))
			if arguments, err := parseArray(s); err != nil {
				return nil, err
			} else {
				if len(arguments) == 0 {
					return nil, fmt.Errorf("arrays must have a value")
				}
				w.Write([]byte(strings.Repeat(",?", len(arguments)-1)))
				w.Write([]byte(")"))
				return arguments, nil
			}
		default:
			return nil, fmt.Errorf("parseOperator: unexpected character")
		}
	}
	return nil, fmt.Errorf("parseOperator: unexpected end of file")
}

func GetPage(
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

		var order string
		var err error
		if order = params.Get("order"); err != nil {
			order = ""
		}

		var limit int
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
		if order != "" {
			db = db.Order(order)
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
		// logrus.Debug(array)
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
