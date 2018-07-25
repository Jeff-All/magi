package controllers

import (
	"magi/resources"
	"magi/models"
)

func CreateGift(
	g models.Gift,
) {
	resources.DB.Exec(`
		INSERT INTO gifts 
		(
			
		)
	`
	)
}
