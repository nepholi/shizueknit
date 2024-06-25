package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func GetSession(r *http.Request) (*sessions.Session, error) {
    return store.Get(r, "shizueknit-session")
}

func SetUserSession(w http.ResponseWriter, r *http.Request, userID uint) error {
    session, err := GetSession(r)
    if err != nil {
        return err
    }
    session.Values["userID"] = userID
    return session.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
    session, err := GetSession(r)
    if err != nil {
        return err
    }
    session.Options.MaxAge = -1
    return session.Save(r, w)
}

func GetUserID(r *http.Request) (uint, bool) {
    session, err := GetSession(r)
    if err != nil {
        return 0, false
    }
    userID, ok := session.Values["userID"].(uint)
    return userID, ok
}