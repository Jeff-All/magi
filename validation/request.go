package validation

import "github.com/Jeff-All/magi/models"

type Request struct {
	models.Base

	ID uint64

	Sheet string
	Row   int

	FamilyID string `validate:"min=1,max=255"`
	Response string `validate:"min=1,max=255"`

	FamilyName string `validate:"min=1,max=255"`
	FirstName  string `validate:"min=1,max=255"`

	Program    string `validate:"min=1,max=255"`
	Age        int    `validate:"min=0,max=100"`
	Gender     string `validate:"oneof=male female"`
	AdultChild string `validate:"oneof=adult child"`

	Gifts []models.Gift
}
