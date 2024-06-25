package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Orders   []Order
}

type Order struct {
    gorm.Model
    UserID      uint
    User        User `gorm:"foreignKey:UserID"`
    TotalAmount float64
    Status      string
}

type Product struct {
    gorm.Model
    Name        string
    Description string
    Price       float64
    Stock       int
    Category    string
}

type OrderItem struct {
    gorm.Model
    OrderID   uint
    ProductID uint
    Quantity  int
    Price     float64
}

type CartItem struct {
    gorm.Model
    UserID    uint
    ProductID uint
    Quantity  int
}