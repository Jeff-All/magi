package auth

import (
	"time"

	res "github.com/Jeff-All/magi/resources"
)

type Group struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name string

	Users   []User   `gorm:"many2many:user_groups"`
	Actions []Action `gorm:"many2many:resource_groups"`
}

func CreateGroup(
	User User,
) (
	*Group,
	error,
) {
	if ok, err := User.Authorize("auth", "groups", "create"); !ok || err != nil {
		return nil, err
	}
	var group *Group
	err := res.DB.Create(group).Error
	return group, err
}

func ReadGroup(
	User User,
	Name string,
) (
	*Group,
	error,
) {
	if ok, err := User.Authorize("auth", "groups", "read"); !ok || err != nil {
		return nil, err
	}
	var group *Group
	err := res.DB.Where("name = ?", Name).First(group).Error
	return group, err
}

func (g *Group) Delete(
	User User,
) (
	bool,
	error,
) {
	if ok, err := User.Authorize("groups", g.Name, "delete"); !ok || err != nil {
		return false, err
	}
	err := res.DB.Delete(g).Error
	return err == nil, err
}

func (g *Group) Update(
	User User,
) (
	bool,
	error,
) {
	if ok, err := User.Authorize("groups", g.Name, "update"); !ok || err != nil {
		return false, err
	}
	err := res.DB.Delete(g).Error
	return err == nil, err
}

// Add user to group
func (g *Group) AddUser(
	User User,
	toAdd User,
) (
	bool,
	error,
) {
	if ok, err := User.Authorize("groups", g.Name, "add_action"); !ok || err != nil {
		return false, err
	}
	err := res.DB.Model(g).Association("Users").Append(toAdd).Error
	return err == nil, err
}

func (g *Group) RemoveUser(
	User User,
	toRemove User,
) (
	bool,
	error,
) {
	if ok, err := User.Authorize("groups", g.Name, "add_action"); !ok || err != nil {
		return false, err
	}
	err := res.DB.Model(g).Association("Users").Append(toRemove).Error
	return err == nil, err
}

func (g *Group) AddAction(
	User User,
	Action Action,
) (
	bool,
	error,
) {
	if ok, err := User.Authorize("groups", g.Name, "add_action"); !ok || err != nil {
		return false, err
	}
	err := res.DB.Model(g).Association("Actions").Append(Action).Error
	return err == nil, err
}

func (g *Group) RemoveAction(
	User User,
	Action Action,
) (
	bool,
	error,
) {
	if ok, err := User.Authorize("groups", g.Name, "remove_action"); !ok || err != nil {
		return false, err
	}
	err := res.DB.Model(g).Association("Actions").Delete(&Action).Error
	return err == nil, err
}
