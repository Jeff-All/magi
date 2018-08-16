package request

import (
	"net/http"

	"github.com/Jeff-All/magi/mock"
	models "github.com/Jeff-All/magi/models"
	util "github.com/Jeff-All/magi/util"
	log "github.com/sirupsen/logrus"
)

var Request iRequest = BaseRequest{}

type iRequest interface {
	PUT(http.ResponseWriter, *http.Request) error
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

	// Parse into models.Request
	var request models.Request
	if err = util.Json.Unmarshal(body, &request); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error unmarshaling request")
		return err
	}

	// Create model entry
	if err = models.Requests.Create(&request); err != nil {
		return err
	}

	// Generate response string
	requestJSONString, err := util.Json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error marshaling request")
		return err
	}

	// Write the response
	if _, err = w.Write(requestJSONString); err != nil {
		return err
	}

	return nil
}
