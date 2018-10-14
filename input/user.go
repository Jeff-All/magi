package input

type User struct {
	Name       Nullable `validate:"min=1,max=255,omitempty"`
	NickName   Nullable `validate:"min=1,max=255,omitempty"`
	GivenName  Nullable `validate:"min=1,max=255,omitempty"`
	FamilyName Nullable `validate:"min=1,max=255,omitempty"`

	Active Nullable `validate:"true"`
	Locked Nullable `validate:"-"`

	Email Nullable `validate:"min=1,max=255,omitempty"`

	Role Nullable `validate:"oneof=admin user,omitempty"`
}

// type User struct {
// 	ID         uint64 `validate:"required"`
// 	Name       string `validate:"min=1,max=255,omitempty"`
// 	NickName   string `validate:"min=1,max=255,omitempty"`
// 	GivenName  string `validate:"min=1,max=255,omitempty"`
// 	FamilyName string `validate:"min=1,max=255,omitempty"`

// 	Active bool `validate:"-"`
// 	Locked bool `validate:"-"`

// 	Email string `validate:"min=1,max=255,omitempty"`

// 	Role string `validate:"oneof=admin user,omitempty"`
// }

// type User struct {
// 	ID        uint64 `gorm:"primary_key;AUTO_INCREMENT"`
// 	CreatedAt *time.Time
// 	UpdatedAt time.Time
// 	DeletedAt *time.Time `sql:"index"`

// 	Name       string
// 	NickName   string
// 	GivenName  string
// 	FamilyName string

// 	Active bool `gorm:"default:false"`
// 	Locked bool `gorm:"default:false"`

// 	Email    string `gorm:"type:varchar(254)"`
// 	SubClaim string `gorm:"unique_index:varchar(254)`

// 	Role string
// }
