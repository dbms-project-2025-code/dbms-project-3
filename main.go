package main

import (
	"crypto/rand"
	"database/sql"
	"dbms/db"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
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

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	h := http.HandlerFunc(welcomeHandler)
	for _, m := range middleware {
		h = m(h)
	}

	mux.HandleFunc("/login", alumni_login_page)
	mux.HandleFunc("POST /logout", logout_handler)
	mux.HandleFunc("/", select_user_page)
	mux.HandleFunc("/donation", donation_page)
	mux.HandleFunc("/academic_history", academic_history_page)
	mux.HandleFunc("/profile/{roll_no}", view_profile_page)
	mux.HandleFunc("/profile", profile_page)
	mux.HandleFunc("POST /delete_employment", delete_employment)
	mux.HandleFunc("POST /save_employment", save_employment)
	mux.HandleFunc("/home", home_page) //alumni
	mux.HandleFunc("GET /directory", directory_page)
	mux.HandleFunc("/staff_login", staff_login_page)
	mux.HandleFunc("/notices", notices_page)
	mux.HandleFunc("POST /rsvp", rsvp_handler)
	mux.HandleFunc("/index", index_page)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Server is running")
}
func alumni_login_page(w http.ResponseWriter, r *http.Request) {
	//users := map[string]string{"Hello": "Templates check"}
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		if db.Authenticate(username, password, "alumni") {
			setSessionID(username, w, r)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			tpl.ExecuteTemplate(w, "login.html", nil) //map[string]any{"incorrect_pass": true})
		}
	} else {
		tpl.ExecuteTemplate(w, "login.html", nil)
	}
}

func logout_handler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // use Secure if on HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	err = db.DeleteSession(cookie.Value)
	if err != nil {
		fmt.Fprintln(w, "ERROR OCCURRED")
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func staff_login_page(w http.ResponseWriter, r *http.Request) {
	//users := map[string]string{"Hello": "Templates check"}
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		if db.Authenticate(username, password, "staff") {
			setSessionID(username, w, r)
			http.Redirect(w, r, "/index", http.StatusSeeOther)
		} else {
			tpl.ExecuteTemplate(w, "staff_login.html", nil) //map[string]any{"incorrect_pass": true})
		}
	} else {
		tpl.ExecuteTemplate(w, "staff_login.html", nil)
	}
}
func index_page(w http.ResponseWriter, r *http.Request) {
	if username, err := checkSession(r); err == nil {
		if _, err := db.GetAlumni(username); err == nil {
			//fmt.Fprintln(w, user_data.Data)
		} else {
			log.Println(err)
			fmt.Fprintln(w, "INDEX PAGE")
		}
	} else {
		log.Println(err)
		fmt.Fprintln(w, "INDEX PAGE")
	}
}

func select_user_page(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "select_user.html", nil)
}

func home_page(w http.ResponseWriter, r *http.Request) {
	if username, err := checkSession(r); err == nil {
		if user_data, err := db.GetAlumni(username); err == nil {
			tpl.ExecuteTemplate(w, "homepage.html", user_data)
		} else {
			log.Println(err)
			fmt.Fprintln(w, "NO USER FOUND")
		}
	} else {
		log.Println(err)
		fmt.Fprintln(w, "SESSION EXPIRED")
	}
}

func academic_history_page(w http.ResponseWriter, r *http.Request) {
	if username, err := checkSession(r); err == nil {
		if list_marks, err := db.GetAcademicHistory(username); err == nil {
			tpl.ExecuteTemplate(w, "academic.html", list_marks)
		} else {
			log.Println(err)
			fmt.Fprintln(w, "NO USER FOUND")
		}
	} else {
		log.Println(err)
		fmt.Fprintln(w, "SESSION EXPIRED")
	}
}

func directory_page(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	department := r.URL.Query().Get("department")
	batch := r.URL.Query().Get("batch")

	alumni, err := db.GetSearchDirectory(query, department, batch)
	log.Println("error", err)
	log.Println("results", alumni)
	tpl.ExecuteTemplate(w, "directory.html", map[string]any{"alumni": alumni})
}

func donation_page(w http.ResponseWriter, r *http.Request) {
	if username, err := checkSession(r); err == nil {
		if total_donated, err := db.GetTotalDonation(username); err == nil {
			donations, _ := db.GetPrevDonation()

			args := make(map[string]any)
			args["donations"] = donations
			args["total"] = total_donated

			if r.Method == http.MethodPost {
				if id, err := handle_donation(username, w, r); err == nil {
					args["success"] = true
					args["id"] = id
				} else {
					args["success"] = false
				}
			}

			tpl.ExecuteTemplate(w, "donation.html", args)
		} else {
			log.Println(err)
			fmt.Fprintln(w, "NO USER FOUND")
		}
	} else {
		log.Println(err)
		fmt.Fprintln(w, "SESSION EXPIRED")
	}
}

func handle_donation(username string, w http.ResponseWriter, r *http.Request) (int64, error) {
	r.ParseForm()
	amount := r.FormValue("amount")
	message := r.FormValue("message")

	var m sql.NullString
	if message == "" {
		m.Valid = false
	} else {
		m.String = message
		m.Valid = true
	}
	a, _ := strconv.ParseFloat(amount, 64)
	id, err := db.AddDonation(username, a, m)
	return id, err
}

func profile_page(w http.ResponseWriter, r *http.Request) {
	if username, err := checkSession(r); err == nil {
		if employments, err := db.GetEmploymentHistory(username); err == nil {
			if alum, err := db.GetAlumni(username); err == nil {
				args := make(map[string]any)
				args["employments"] = employments
				args["alumnus"] = alum

				tpl.ExecuteTemplate(w, "profile.html", args)
			} else {
				fmt.Fprintln(w, "NO USER FOUND")
				log.Println(err)

			}
		} else {
			fmt.Fprintln(w, "NO USER FOUND")
			log.Println(err)
		}
	} else {
		log.Println(err)
		fmt.Fprintln(w, "SESSION EXPIRED")
	}

}
func view_profile_page(w http.ResponseWriter, r *http.Request) {
	username, err := checkSession(r)
	if err != nil {
		fmt.Fprintln(w, "SESSION EXPIRED")
		log.Println(err)
		return
	}
	roll_no := r.PathValue("roll_no")

	alum, err := db.GetAlumni(roll_no)

	if err != nil {
		fmt.Fprintln(w, "page doesn't exist")
		log.Println(err)
		return
	}
	if username == alum.Roll_no {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

	employments, _ := db.GetEmploymentHistory(roll_no)
	args := make(map[string]any)
	args["employments"] = employments
	args["alumnus"] = alum

	tpl.ExecuteTemplate(w, "profile_view.html", args)

}

func save_employment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	company := r.FormValue("company")
	start := r.FormValue("startYear")
	end := r.FormValue("endYear")
	designation := r.FormValue("designation")
	location := r.FormValue("location")

	s, _ := strconv.Atoi(start)
	e, _ := strconv.Atoi(end)
	data := db.Employment{
		Company:     company,
		Starting:    s,
		Ending:      e,
		Designation: designation,
		Location:    location,
	}

	if username, err := checkSession(r); err == nil {
		data.Roll_no = username
		if id == "" {
			err := db.AddEmploymentHistory(data)
			if err != nil {
				log.Panicln(err)
			}
		} else {
			data.Id, _ = strconv.Atoi(id)
			err := db.UpdateEmploymentHistory(data)
			if err != nil {
				log.Panicln(err)
			}
		}
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	} else {
		fmt.Fprintln(w, "SESSION EXPIRED")
	}
}

func delete_employment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	i, _ := strconv.Atoi(id)

	err := db.DeleteEmploymentHistory(i)
	if err != nil {
		log.Panicln(err)
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
func notices_page(w http.ResponseWriter, r *http.Request) {
	username, err := checkSession(r)
	if err != nil {
		fmt.Fprintln(w, "SESSION EXPIRED")
		return
	}

	announcements, err := db.GetAllAnnouncements(username)

	if err != nil {
		log.Println(err)
	}
	tpl.ExecuteTemplate(w, "notices.html", announcements)
}

func rsvp_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("eventid")
	response := r.FormValue("response")

	username, err := checkSession(r)
	if err != nil {
		fmt.Fprintln(w, "SESSION EXPIRED")
		return
	}

	eventid, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprint(w, "SOMETHING WENT WRONG")
		log.Println(err)
		return
	}

	err = db.AddRVSP(eventid, username, response)
	if err != nil {
		fmt.Fprint(w, "SOMETHING WENT WRONG")
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/notices", http.StatusSeeOther)
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
