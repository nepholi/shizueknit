package database

import (
	"shizueknit/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "username:password@tcp(local:3306)/shizueknit?charset=utf8mb4&parseTime=True&loc=Local"

    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
    })
    if err != nil {
        return err
    }

    // 自動遷移
    err = DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}, &models.CartItem{})
    if err != nil {
        return err
    }

    return nil
}