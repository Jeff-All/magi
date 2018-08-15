package auth

import (
	"time"

	res "github.com/Jeff-All/magi/resources"
)

type Resource struct {
	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name string

	Actions []Action
}

func CreateResource(
	User User,
	Name string,
) (
	Resource *Resource,
	err error,
) {
	if ok, err := User.Authorize("auth", "resources", "create"); !ok || err != nil {
		return nil, err
	}
	err = res.DB.Create(Resource).GetError()
	return Resource, err
}
