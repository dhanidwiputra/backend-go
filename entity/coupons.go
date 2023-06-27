package entity

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Code         uuid.UUID `json:"code"`
	Description  string    `json:"description"`
	IssuerID     uint      `json:"issuer_id"`
	Discount     int       `json:"discount"`
	Availability bool      `json:"availability"`
	Users        []User    `gorm:"many2many:users_coupons;" json:"users,omitempty"`
}

// type Person struct {
// 	Age  int    `json:"age"`
// 	Name string `json:"name"`
// }

// type Person struct {
// 	Age  interface{} `json:"age"`
// 	Name interface{} `json:"name"`
// }
