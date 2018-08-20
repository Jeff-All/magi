package models

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/Jeff-All/magi/errors"
	"github.com/sirupsen/logrus"
)

var Requests iRequests = _Requests{}

type iRequests interface {
	Create(...*Request) error
	Get(string) (*Request, error)
	Delete(...interface{}) error
	Page(int, int) ([]Request, error)
	CreateGift(uint64, []*Gift) error
}

type _Requests struct{}

type Request struct {
	BaseModel

	ID uint64 `gorm:"primary_key;AUTO_INCREMENT"`

	Agency string

	Gifts []Gift `gorm:"foreignkey:RequestID"`
}

func (request _Requests) RemoveGift(
	ID uint64,
	giftID uint64,
) error {
	if err := DB.Model(&Request{ID: ID}).Association("Gifts").Delete(Gift{ID: giftID}).GetError(); err != nil {
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

func (request _Requests) CreateGift(
	ID uint64,
	gifts []*Gift,
) error {
	for _, cur := range gifts {
		if cur != nil {
			if err := DB.Model(&Request{ID: ID}).Association("Gifts").Append(cur).GetError(); err != nil {
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

func (requests _Requests) Create(request ...*Request) error {
	if request == nil {
		return fmt.Errorf("request was nil")
	}
	for _, cur := range request {
		if cur != nil {
			err := DB.Create(cur).GetError()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (requests _Requests) Get(id string) (*Request, error) {
	var request Request
	err := DB.Preload("Gifts").First(&request, id).GetError()
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

func (requests _Requests) GetWithGifts(id string) (*Request, error) {
	var request Request
	err := DB.Preload("Gifts").First(&request, id).GetError()
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

func (requests _Requests) Delete(ids ...interface{}) error {
	if err := DB.Where("id IN (?)", ids...).Delete(Request{}).GetError(); err != nil {
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

func (requests _Requests) Page(
	limit int,
	offset int,
) ([]Request, error) {
	logrus.WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	}).Debug("models.Request.Page()")
	requestArray := make([]Request, 0, limit)
	if err := DB.Limit(limit).Offset(offset).Preload("Gifts").Find(&requestArray).GetError(); err != nil {
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
