package handlers

import (
	"html/template"
	"net/http"
	"shizueknit/internal/database"
	"shizueknit/internal/models"
	"shizueknit/internal/utils"
)
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/layout.html", "templates/home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userID, loggedIn := utils.GetUserID(r)
    data := map[string]interface{}{
        "LoggedIn": loggedIn,
        "UserID":   userID,
    }

    tmpl.Execute(w, data)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/layout.html", "templates/register.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    user := models.User{
        Username: username,
        Email:    email,
        Password: hashedPassword,
    }

    result := database.DB.Create(&user)
    if result.Error != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }

    err = utils.SetUserSession(w, r, user.ID)
    if err != nil {
        http.Error(w, "Error setting session", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl, err := template.ParseFiles("templates/layout.html", "templates/login.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    var user models.User
    result := database.DB.Where("username = ?", username).First(&user)
    if result.Error != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if !utils.CheckPasswordHash(password, user.Password) {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    err := utils.SetUserSession(w, r, user.ID)
    if err != nil {
        http.Error(w, "Error setting session", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    err := utils.ClearSession(w, r)
    if err != nil {
        http.Error(w, "Error clearing session", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}