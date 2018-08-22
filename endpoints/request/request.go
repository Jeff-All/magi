package request

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/mock"
	models "github.com/Jeff-All/magi/models"
	util "github.com/Jeff-All/magi/util"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var Request iRequest = BaseRequest{}

type iRequest interface {
	PUT(http.ResponseWriter, *http.Request) error
	GET(http.ResponseWriter, *http.Request) error
	DELETE(http.ResponseWriter, *http.Request) error
	GETPAGE(http.ResponseWriter, *http.Request) error
	PUTGift(http.ResponseWriter, *http.Request) error
}

type MockRequest struct {
	Mock mock.Mock
}

type BaseRequest struct{}

func (i BaseRequest) PUT(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Read body
	body, err := util.IOUtil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"body": string(body),
	}).Debug("/Request:PUT")

	// Parse into models.Request
	var requests []*models.Request
	if err = util.Json.Unmarshal(body, &requests); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error unmarshaling request")
		return err
	}

	// Create model entry
	if err = models.Requests.Create(requests...); err != nil {
		return err
	}

	// Generate response string
	requestJSONString, err := util.Json.Marshal(requests)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error marshaling request")
		return err
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(requestJSONString); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (i BaseRequest) GETPAGE(
	w http.ResponseWriter,
	r *http.Request,
) error {
	log.Debug("endpoints.request.GET()")

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

	requests, err := models.Requests.Page(limit, offset)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints.request.GET(): Error while querying request",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}

	responseBody, err := json.Marshal(requests)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints.Request.GET(): Error whil marshaling response",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints.Request.GET(): Error while writing response",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	w.WriteHeader(http.StatusFound)
	return nil
}

func (i BaseRequest) GET(
	w http.ResponseWriter,
	r *http.Request,
) error {
	log.Debug("endpoints.request.GET()")
	vars := mux.Vars(r)
	if vars == nil {
		return errors.CodedError{
			Message:  "endpoints.request.GET(): Error while pulling mux.Vars",
			HTTPCode: http.StatusInternalServerError,
		}
	}
	id, ok := vars["id"]
	if !ok {
		return errors.CodedError{
			Message:  "endpoints.request.GET(): {id} is undefined",
			HTTPCode: http.StatusInternalServerError,
		}
	}
	request, err := models.Requests.Get(id)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints.request.GET(): Error while querying request",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	responseBody, err := json.Marshal(request)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints.Request.GET(): Error whil marshaling response",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints.Request.GET(): Error while writing response",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func (i BaseRequest) DELETE(
	w http.ResponseWriter,
	r *http.Request,
) error {
	log.Debug("endpoints.request.DELETE()")
	vars := mux.Vars(r)
	if vars == nil {
		return errors.CodedError{
			Message:  "endpoints.request.DELETE(): Error while pulling mux.Vars",
			HTTPCode: http.StatusInternalServerError,
			Code:     0,
		}
	}
	ids, ok := vars["id"]
	if !ok {
		return errors.CodedError{
			Message:  "endpoints.request.DELETE(): {id} is undefined",
			HTTPCode: http.StatusInternalServerError,
			Code:     1,
		}
	}

	err := models.Requests.Delete([]interface{}{strings.Split(ids, ",")}...)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints.request.DELETE(): Error while querying request",
			HTTPCode: http.StatusInternalServerError,
			Err:      err,
		}
	}
	w.WriteHeader(http.StatusFound)
	return nil
}

func (i BaseRequest) PUTGift(
	w http.ResponseWriter,
	r *http.Request,
) error {
	vars := mux.Vars(r)
	if vars == nil {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift()",
			HTTPCode: http.StatusInternalServerError,
			Code:     0,
		}
	}
	idString, ok := vars["id"]
	if !ok {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift(): {id} is undefined",
			HTTPCode: http.StatusInternalServerError,
			Code:     1,
		}
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift()",
			HTTPCode: http.StatusInternalServerError,
			Code:     2,
			Err:      err,
		}
	}
	body, err := util.IOUtil.ReadAll(r.Body)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift()",
			HTTPCode: http.StatusInternalServerError,
			Code:     3,
			Err:      err,
		}
	}
	var gifts []*models.Gift
	if err = util.Json.Unmarshal(body, &gifts); err != nil {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift()",
			HTTPCode: http.StatusInternalServerError,
			Code:     4,
			Err:      err,
		}
	}
	if err = models.Requests.CreateGift(uint64(id), gifts); err != nil {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift()",
			HTTPCode: http.StatusInternalServerError,
			Code:     5,
			Err:      err,
		}
	}

	requestJSONString, err := util.Json.Marshal(gifts)
	if err != nil {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift()",
			HTTPCode: http.StatusInternalServerError,
			Code:     6,
			Err:      err,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(requestJSONString); err != nil {
		return errors.CodedError{
			Message:  "endpoints/request.PUTGift()",
			HTTPCode: http.StatusInternalServerError,
			Code:     7,
			Err:      err,
		}
	}
	w.WriteHeader(http.StatusCreated)
	return nil

}
