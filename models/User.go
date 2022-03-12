package models

import (
	"fmt"
	"time"

	"github.com/nadhirfr/codefood/helpers"

	"gorm.io/gorm"
)

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id" form:"id" swaggertype:"integer"`
	Username  string     `gorm:"type:varchar(300);unique_index" form:"username" json:"username" binding:"required,max=300" swaggertype:"string" example:"rozam"`
	Password  string     `gorm:"size:300" form:"password" json:"password,omitempty" binding:"required,min=6,max=300"`
	Serves    []Serve    `gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `form:"createdAt" json:"createdAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	UpdatedAt time.Time  `form:"updatedAt" json:"updatedAt" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
	DeletedAt *time.Time `form:"deletedAt" json:"deletedAt" gorm:"index" swaggertype:"string" example:"2021-04-12T00:39:11.652+07:00"`
}

// BeforeUpdate : hook before a user is updated
func (u *User) BeforeSave(tx *gorm.DB) (err error) {

	if u.Password != "" {
		hash, err := helpers.MakePassword(u.Password)
		if err != nil {
			return nil
		}
		u.Password = hash
	}

	fmt.Println(u.Password)
	return
}

type UserLogin struct {
	Username string `gorm:"type:varchar(20);unique_index" form:"username" json:"username" binding:"required,max=300" swaggertype:"string" example:"rozam"`
	Password string `gorm:"size:72" form:"password" json:"password,omitempty" binding:"required,min=6,max=300"`
}

type UserResult201 struct {
	ID       uint   `json:"id" swaggertype:"integer" example:"12345"`
	Username string `json:"username" swaggertype:"string" example:"account name"`
}
