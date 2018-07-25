package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"github.com/Jeff-All/magi/auth"
	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/models"
	"github.com/Jeff-All/magi/responses"
	"net/http"

	res "github.com/Jeff-All/magi/resources"

	log "github.com/sirupsen/logrus"
)

func PutGift(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Debugf("/gifts PUT")

	u, err := auth.AuthRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":    err.Error(),
			"Endpoint": "/gifts",
			"Action":   "PUT",
		}).Error("Error Authing User")
		response := responses.Error{
			Code:  errors.Default,
			Error: err.Error(),
		}

		responseString, _ := json.Marshal(response)

		w.Write(responseString)
		return
	} else if u == nil {
		log.Debug("Failed Auth")
		w.Write([]byte("Invalid User"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Endpoint": "/gifts",
			"Action":   "PUT",
			"Error":    err,
		}).Error("Error Reading Body")

		response := responses.Error{
			Code:  errors.Default,
			Error: err.Error(),
		}

		responseString, _ := json.Marshal(response)

		w.Write(responseString)
		return
	}

	var gift models.Gift
	json.Unmarshal(body, &gift)
	gift.ID = 0

	var errDB = res.DB.Create(&gift)
	jsonString, _ := json.Marshal(gift)
	if errDB.Error != nil {
		log.WithFields(log.Fields{
			"Endpoint": "/gifts",
			"Action":   "PUT",
			"Error":    errDB.Error,
			"Gift":     string(jsonString),
			"Table":    "gifts",
		}).Error("Error Inserting Into Table")

		response := responses.Error{
			Code:  errors.Default,
			Error: errDB.Error.Error(),
		}

		responseString, _ := json.Marshal(response)

		w.Write(responseString)
		return
	}

	w.Write(jsonString)
}
