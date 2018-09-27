package models

type Request struct {
	Base

	ID uint64 `gorm:"primary_key;AUTO_INCREMENT"`

	Sheet string `gorm:"-"`
	Row   int    `gorm:"-"`

	Batch    string `gorm:"unique_index:batch_id"`
	FamilyID string `gorm:"unique_index:batch_id"`
	Response string `gorm:"unique_index:batch_id"`

	FamilyName string
	FirstName  string

	Program    string
	Age        int
	Gender     string
	AdultChild string

	Gifts []Gift `gorm:"foreignkey:RequestID"`
}
