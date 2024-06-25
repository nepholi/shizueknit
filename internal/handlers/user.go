package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
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

    account := r.FormValue("account")
    email := r.FormValue("email")
    password := r.FormValue("password")

	if err := isUserRegistrationValidate(account, email, password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    user := models.User{
        Account: account,
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

    account := r.FormValue("account")
    password := r.FormValue("password")

    var user models.User
    result := database.DB.Where("account = ?", account).First(&user)
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

// 帳號輸入欄位驗證
func validateAccountField(account string) bool {
    // 限制只能輸入英文及數字，大小寫皆可，上限為30字元
    pattern := `^[a-zA-Z0-9]{1,30}$`
    matched, _ := regexp.MatchString(pattern, account)
    return matched
}

// 電子郵件輸入欄位驗證
func validateEmailField(email string) bool {
    // 檢查電子郵件格式是否正確
    pattern := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}

// 密碼輸入欄位驗證
func validatePasswordField(password string) bool {
    // 限制密碼最少8字元
    return len(password) >= 8
}

// 檢查電子郵件是否已被註冊
func isEmailRegistered(email string) bool {
    // 在資料庫中查詢是否有相同的電子郵件
	var user models.User
	fmt.Println("isEmailRegistered:",email)
	u := database.DB.Where("email = ?", email).First(&user)
	fmt.Printf("%v",u)
    return u!=nil
}

// 處理用戶註冊
func isUserRegistrationValidate(account, email, password string) error {
	fmt.Println("isUserRegistrationValidate")
    if !validateAccountField(account) {
        return errors.New("帳號格式不正確")
    }
    if !validateEmailField(email) {
        return errors.New("電子郵件格式不正確")
    }
    if !validatePasswordField(password) {
        return errors.New("密碼必須至少8個字元")
    }
    if isEmailRegistered(email) {
        return errors.New("電子郵件已被註冊，請檢查")
    }
    return nil
}