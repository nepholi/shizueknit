package main

import (
	"log"
	"net/http"

	"shizueknit/internal/database"
	"shizueknit/internal/handlers"
	"shizueknit/internal/models"
)

func initializeTestData() {
    products := []models.Product{
        {Name: "經典羊毛圍巾", Description: "100% 純羊毛製作，柔軟保暖", Price: 1200, Stock: 50, Category: "圍巾"},
        {Name: "針織毛帽", Description: "舒適保暖的針織毛帽，適合冬季使用", Price: 800, Stock: 100, Category: "帽子"},
        {Name: "手工編織披肩", Description: "精美手工編織，適合各種場合", Price: 1500, Stock: 30, Category: "披肩"},
    }

    for _, product := range products {
        database.DB.Create(&product)
    }
}

func main() {
    // 連接數據庫
    err := database.InitDB()
    if err != nil {
        log.Fatal(err)
    }
	//測試數據 TODO:刪除
	//initializeTestData()
    
	// 設置路由
    http.HandleFunc("/", handlers.HomeHandler)
    http.HandleFunc("/register", handlers.RegisterHandler)
    http.HandleFunc("/login", handlers.LoginHandler)
    http.HandleFunc("/logout", handlers.LogoutHandler)
    http.HandleFunc("/products", handlers.ProductListHandler)
    http.HandleFunc("/product", handlers.ProductDetailHandler)

    // 設置靜態文件服務
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // 啟動服務器
    log.Println("ShizueKnit Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}