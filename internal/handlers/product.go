package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"shizueknit/internal/database"
	"shizueknit/internal/models"
	"shizueknit/internal/utils"
)

func ProductListHandler(w http.ResponseWriter, r *http.Request) {
    var products []models.Product
    result := database.DB.Find(&products)
    if result.Error != nil {
        http.Error(w, "Error fetching products", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/layout.html", "templates/product_list.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userID, loggedIn := utils.GetUserID(r)
    data := map[string]interface{}{
        "Products": products,
        "LoggedIn": loggedIn,
        "UserID":   userID,
    }

    tmpl.Execute(w, data)
}

func ProductDetailHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    var product models.Product
    result := database.DB.First(&product, id)
    if result.Error != nil {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    tmpl, err := template.ParseFiles("templates/layout.html", "templates/product_detail.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userID, loggedIn := utils.GetUserID(r)
    data := map[string]interface{}{
        "Product":  product,
        "LoggedIn": loggedIn,
        "UserID":   userID,
    }

    tmpl.Execute(w, data)
}