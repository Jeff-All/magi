package output

import "time"

type User struct {
	ID        uint64
	CreatedAt *time.Time `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	Name       string
	NickName   string
	GivenName  string
	FamilyName string

	Active   bool
	Email    string
	SubClaim string `json:"-"`

	Role string
}
