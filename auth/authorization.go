package auth

func Authorize(
	Category int,
	Resource string,
	Action string,
	User User,
) (bool, error) {
	return false, nil
	// var resource models.Resource
	// err := res.DB.Model(&User)
	// 	.Association("Groups")

	// 	// .Association("Resources")
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	log.WithFields(log.Fields{
	// 		"Category": Category,
	// 		"Resource": Resource,
	// 		"Action":   Action,
	// 		"Error":    err.Error(),
	// 	}).Error("Failed to query resources")
	// 	return false, err
	// } else if err == gorm.ErrRecordNotFound {
	// 	log.Debug("Unable to find Resource")
	// 	return false, nil
	// }

}
