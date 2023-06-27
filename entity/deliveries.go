package entity

import (
	"time"

	"gorm.io/gorm"
)

type Delivery struct {
	gorm.Model
	OrderID      uint      `json:"order_id"`
	Address      string    `json:"address"`
	DeliveryDate time.Time `json:"delivery_date"`
	Status       string    `json:"status"`
}
