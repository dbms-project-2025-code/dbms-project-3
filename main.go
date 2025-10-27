package main

import (
	"crypto/rand"
	"dbms/db"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

var users = map[string]string{"admin": "password", "di": "2004"} //username : password
var cookies = map[string]string{"cook": "admin"}                 //cookies : username
var data = map[string]string{"admin": "secret"}                  //username : secret data

var tpl = template.Must(template.ParseGlob("./templates/*.html"))

var middleware = []func(http.HandlerFunc) http.HandlerFunc{
	checkSessionMiddleware,
}

func main() {
	db.Connect("./app.db")
	defer db.Close()
	go db.GarbageCollector()

	h := http.HandlerFunc(welcomeHandler)
	for _, m := range middleware {
		h = m(h)
	}

	http.HandleFunc("/login", login_page)
	http.HandleFunc("/", index_page)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Server is running")
}
func login_page(w http.ResponseWriter, r *http.Request) {
	//users := map[string]string{"Hello": "Templates check"}
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		if db.Authenticate(username, password) {
			setSessionID(username, w, r)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			tpl.ExecuteTemplate(w, "login.html", map[string]any{"incorrect_pass": true})
		}
	} else {
		tpl.ExecuteTemplate(w, "login.html", nil)
	}
}

func index_page(w http.ResponseWriter, r *http.Request) {
	if username, err := checkSession(r); err == nil {
		if data, err := db.GetUserData(username); err == nil {
			fmt.Fprintln(w, data)
		} else {
			log.Println(err)
			fmt.Fprintln(w, "INDEX PAGE")
		}
	} else {
		log.Println(err)
		fmt.Fprintln(w, "INDEX PAGE")
	}
}

func checkSession(r *http.Request) (string, error) {
	//Returns username and check flag true means session found
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("No cookie recieved")
		}
		return "", errors.New("Unknown Error With Cookie")
	}
	return db.CheckSession(cookie.Value)
}

func checkSessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "No cookie found", http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if cookies[cookie.Value] == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func setSessionID(username string, w http.ResponseWriter, r *http.Request) {
	session_id := generateSessionID()
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session_id,
		Path:     "/", // Cookie is valid for entire site
		Expires:  time.Now().Add(time.Hour),
		MaxAge:   3600,  // 1 hour (in seconds)
		HttpOnly: true,  // Not accessible to JS
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	db.StoreSession(username, session_id)
	//fmt.Fprintf(w, "Cookie has been set!")
}

func generateSessionID() string {
	id := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		panic("failed to generate session id")
	}

	return base64.RawURLEncoding.EncodeToString(id)
}
