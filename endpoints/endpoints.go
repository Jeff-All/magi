package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Jeff-All/magi/models"
	"github.com/Jeff-All/magi/responses"

	res "github.com/Jeff-All/magi/resources"

	log "github.com/sirupsen/logrus"
)

// Agency
// CRUD - Admin

// Gifts
// 	Create --- Admin
// 	Retrieve - All
// 	Update --- Admin
// 	Delete --- Admin
//

// Tags
// 	Create --- Admin
//	Retrieve - All
//	Update --- Admin
// 	Delete --- Admin

func PutGift(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Debugf("/gifts PUT")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Endpoint": "/gifts",
			"Action":   "PUT",
			"Error":    err,
		}).Error("Error Reading Body")

		response := responses.Error{
			// Code:  errors.Default,
			Error: err.Error(),
		}

		responseString, _ := json.Marshal(response)

		w.Write(responseString)
		return
	}

	var gift models.Gift
	json.Unmarshal(body, &gift)
	gift.ID = 0

	err = res.DB.Create(&gift).Error
	jsonString, _ := json.Marshal(gift)
	if err != nil {
		log.WithFields(log.Fields{
			"Endpoint": "/gifts",
			"Action":   "PUT",
			"Error":    err.Error(),
			"Gift":     string(jsonString),
			"Table":    "gifts",
		}).Error("Error Inserting Into Table")

		response := responses.Error{
			// Code:  errors.Default,
			Error: err.Error(),
		}

		responseString, _ := json.Marshal(response)

		w.Write(responseString)
		return
	}

	w.Write(jsonString)
}
