package auth

/*
auth:group:add_user
*/

func Authorize(
	Category int,
	Resource string,
	Action string,
	User User,
) (bool, error) {
	return false, nil

	/*
		SELECT r.id
		FROM users AS u
		JOIN user_groups AS ug
			ON u.id = ug.user_id
		JOIN groups as g
			ON g.id = ug.group_id
		JOIN resource_groups as rg
			ON rg.group_id = ug.group_id
		JOIN resources as r
			ON r.id = rg.resource_id
		WHERE r.Category = ?
		AND r.Resource = ?
		AND r.Action = ?
	*/

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
