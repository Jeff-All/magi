package request

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/Jeff-All/magi/actions"
	"github.com/Jeff-All/magi/errors"
	"github.com/Jeff-All/magi/models"
)

func RemoveGift(
	ID uint64,
	giftID uint64,
) error {
	if err := actions.DB.Model(&models.Request{ID: ID}).Association("Gifts").Delete(models.Gift{ID: giftID}).GetError(); err != nil {
		return errors.CodedError{
			Message:  "models.Request.RemoveGift()",
			HTTPCode: 500,
			Fields: logrus.Fields{
				"ID":      ID,
				"gift ID": giftID,
			},
			Err: err,
		}
	}
	return nil
}

func CreateGift(
	ID uint64,
	gifts []*models.Gift,
) error {
	for _, cur := range gifts {
		if cur != nil {
			if err := actions.DB.Model(&models.Gift{ID: ID}).Association("Gifts").Append(cur).GetError(); err != nil {
				return errors.CodedError{
					Message:  "models.Request.CreateGift()",
					HTTPCode: 500,
					Fields: logrus.Fields{
						"ID": ID,
					},
					Err: err,
				}
			}
		}
	}
	return nil
}

func Create(request ...*models.Request) error {
	if request == nil {
		return fmt.Errorf("request was nil")
	}
	for _, cur := range request {
		if cur != nil {
			err := actions.DB.Create(cur).GetError()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Get(id string) (*models.Gift, error) {
	var request models.Gift
	err := actions.DB.Preload("Gifts").First(&request, id).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.CodedError{
			Message:  "models.Request.Get()",
			HTTPCode: 500,
			Fields: logrus.Fields{
				"id": id,
			},
			Err: err,
		}
	} else if err != nil {
		return nil, errors.CodedError{
			Message:  "Can't find requests",
			HTTPCode: 404,
			Fields: logrus.Fields{
				"id": id,
			},
			Err: err,
		}
	}
	return &request, nil
}

func GetWithGifts(id string) (*models.Request, error) {
	var request models.Request
	err := actions.DB.Preload("Gifts").First(&request, id).GetError()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.CodedError{
			Message:  "models.Request.Get()",
			HTTPCode: 500,
			Fields: logrus.Fields{
				"id": id,
			},
			Err: err,
		}
	} else if err != nil {
		return nil, errors.CodedError{
			Message:  "Can't find requests",
			HTTPCode: 404,
			Fields: logrus.Fields{
				"id": id,
			},
			Err: err,
		}
	}
	return &request, nil
}

func Delete(ids ...interface{}) error {
	if err := actions.DB.Where("id IN (?)", ids...).Delete(models.Gift{}).GetError(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.CodedError{
				Message:  "Can't find requests",
				HTTPCode: 404,
				Fields: logrus.Fields{
					"id": ids,
				},
				Err: err,
			}
		}
		return errors.CodedError{
			Message:  "models.Request.Delete()",
			HTTPCode: 500,
			Fields: logrus.Fields{
				"id": ids,
			},
			Err: err,
		}
	}
	return nil
}

func Page(
	limit int,
	offset int,
) ([]models.Request, error) {
	logrus.WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	}).Debug("models.Request.Page()")
	requestArray := make([]models.Request, 0, limit)
	if err := actions.DB.Limit(limit).Offset(offset).Preload("Gifts").Find(&requestArray).GetError(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.CodedError{
				Message:  "Can't find requests",
				HTTPCode: 404,
				Fields: logrus.Fields{
					"limit":  limit,
					"offset": offset,
				},
				Err: err,
			}
		}
		return nil, errors.CodedError{
			Message:  "models.Request.Page()",
			HTTPCode: 500,
			Fields: logrus.Fields{
				"limit":  limit,
				"offset": offset,
			},
			Err: err,
		}
	}
	return requestArray, nil
}
