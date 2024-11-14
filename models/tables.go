package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Enum for OrderStatus
type OrderStatus string

const (
	Pending    OrderStatus = "PENDING"
	Processing OrderStatus = "PROCESSING"
	Shipped    OrderStatus = "SHIPPED"
	Delivered  OrderStatus = "DELIVERED"
	Canceled   OrderStatus = "CANCELED"
	Returned   OrderStatus = "RETURNED"
)

// OrderItem model
type Item struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	ProductID string  `gorm:"type:string;not null;index" json:"product_id"`
	Quantity  int32   `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"`
	OrderID   string  `gorm:"type:string;not null;index" json:"order_id"`
}

// Order model
type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	OrderID    string      `gorm:"type:string;uniqueIndex" json:"order_id"`
	CustomerID uuid.UUID   `gorm:"type:uuid;not null" json:"customer_id"`
	Items      []Item      `gorm:"foreignKey:OrderID" json:"items"`
	Address    Address     `gorm:"foreignKey:OrderID" json:"address"`
	TotalPrice float64     `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Phone      string      `json:"phone"`
	Status     OrderStatus `gorm:"type:varchar(20);default:PENDING;not null" json:"status"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type Address struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	OrderID      string `gorm:"type:string;uniqueIndex;not null" json:"order_id"`
	AddressLine1 string `gorm:"type:string;not null" json:"address_line_1"`
	AddressLine2 string `gorm:"type:string;not null" json:"address_line_2"`
	City         string `gorm:"type:string;not null" json:"city"`
	State        string `gorm:"type:string;not null" json:"state"`
	Zip          string `gorm:"type:string;not null" json:"zip"`
	Country      string `gorm:"type:string;default:IN;not null" json:"country"`
}

// BeforeCreate hook to set UUIDs before inserting a new record
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.OrderID == "" {
		o.OrderID = uuid.New().String()
	}

	o.Address.OrderID = o.OrderID //.String()

	for i := range o.Items {
		if o.Items[i].ID == 0 {
			o.Items[i].OrderID = o.OrderID //.String()
		}
	}
	return
}

type Product struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;" json:"product_id"`
	Name                  string    `gorm:"type:varchar(255);not null" json:"name"`
	Quantity              int32     `gorm:"not null" json:"quantity"`
	Type                  string    `gorm:"type:varchar(20);not null" json:"type"`
	Category              string    `gorm:"type:varchar(100);not null" json:"category"`
	ImageUrls             string    `gorm:"type:json" json:"image_urls"`
	Price                 float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Width                 float64   `gorm:"type:decimal(5,2)" json:"width"`
	Height                float64   `gorm:"type:decimal(5,2)" json:"height"`
	Weight                float64   `gorm:"type:decimal(5,2)" json:"weight"`
	ShippingBasePrice     float64   `gorm:"type:decimal(10,2);not null" json:"shipping_base_price"`
	BaseDeliveryTimelines int32     `gorm:"not null" json:"base_delivery_timelines"`
	SellerId              uuid.UUID `gorm:"type:uuid;not null" json:"seller_id"`
}
