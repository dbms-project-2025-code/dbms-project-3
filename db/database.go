package db

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	_ "modernc.org/sqlite"
	"time"
)

type Alumnus struct {
	Roll_no    string
	Name       string
	Email      string
	Phone_no   string
	Batch      string
	Department string
}

type Academic_history struct {
	Semester int
	SGPA     float64
	CGPA     float64
}
type Donation struct {
	Name   string
	Amount float64
}
type Employment struct {
	Id          int
	Roll_no     string
	Starting    int
	Ending      int
	Company     string
	Designation string
	Location    string
}
type Announcemnt struct {
	Id          int
	Type        string
	Title       string
	Description sql.NullString
	Date        sql.NullString
	Location    sql.NullString
	Created     string
	RSVP        sql.NullString
}

var db *sql.DB

func Connect(path string) {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
		//return err
	}
	for _, createTable := range SchemaTables {
		_, err = db.Exec(createTable)

		if err != nil {
			log.Panic(err)
			//return err
		}
	}
	insertTestingData()
	log.Println("Database initialized successfully")
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Database Connection Closed")
}

func insertTestingData() {
	for username, c := range users {
		hash_pass, e := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
		if e != nil {
			log.Panic(e)
		}

		_, err := db.Exec(insertCredentials, username, string(hash_pass), c.User_type)
		if err != nil {
			log.Panic(err)
		}
	}

	for _, insertCommand := range insertData {
		_, err := db.Exec(insertCommand)
		if err != nil {
			log.Panic(err)
		}
	}
}

func Authenticate(username string, password string, user_type string) bool {

	if user_type != "staff" && user_type != "alumni" {
		log.Panicf("invalid user type: %s", user_type)
	}

	row := db.QueryRow(authenticateUser, username, user_type)
	if row.Err() != nil {
		log.Panic(row.Err())
	}

	var hash_pass string

	if err := row.Scan(&hash_pass); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with the given ID.")
		} else {
			log.Panic(err)
		}
		return false
	}
	log.Println(hash_pass)
	err := bcrypt.CompareHashAndPassword([]byte(hash_pass), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func CheckSession(session_id string) (string, error) {
	query := "Select username, expiry from sessions where session_id = ? ;"
	row := db.QueryRow(query, session_id)
	if row.Err() != nil {
		log.Panic(row.Err())
	}

	var username string
	var expiry time.Time

	if err := row.Scan(&username, &expiry); err != nil {
		if err == sql.ErrNoRows {
			//log.Println("No user found with the given ID.")
		} else {
			log.Panic(err)
		}
		return "", errors.New("No User found with the given ID.")
	}

	if expiry.Before(time.Now()) {
		return "", errors.New("Session Expired")
	}
	return username, nil
}

func StoreSession(username string, session_id string) {
	command := "INSERT OR REPLACE into sessions(session_id, username, expiry) values (?,?,?);"
	_, err := db.Exec(command, session_id, username, time.Now().Add(60*time.Minute))
	if err != nil {
		log.Panic(err)
	}
}

func DeleteSession(cookie string) error {
	_, err := db.Exec(deleteSession, cookie)
	return err
}
func GetAlumni(username string) (Alumnus, error) {
	row := db.QueryRow(getAlumni, username)
	if row.Err() != nil {
		log.Panic(row.Err())
	}

	var u Alumnus
	err := row.Scan(&u.Roll_no, &u.Name, &u.Email, &u.Phone_no, &u.Batch, &u.Department)

	if err != nil {
		if err == sql.ErrNoRows {
			//log.Println("No user found with the given ID.")
		} else {
			log.Println(err)
		}
		return Alumnus{}, errors.New("No User data found")
	}
	return u, nil
}

func GetSearchDirectory(query string, department string, batch string) ([]Alumnus, error) {
	rows, err := db.Query(searchDirectory, query, department, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumni []Alumnus

	for rows.Next() {
		var a Alumnus
		err := rows.Scan(&a.Name, &a.Roll_no, &a.Batch, &a.Department)
		log.Println("one row", a)
		if err != nil {
			return alumni, err
		}
		alumni = append(alumni, a)
	}
	if err = rows.Err(); err != nil {
		return alumni, err
	}
	return alumni, nil
}

func GetAcademicHistory(username string) ([]Academic_history, error) {
	rows, err := db.Query(getAcademicHistory, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list_marks []Academic_history

	for rows.Next() {
		var a Academic_history
		err := rows.Scan(&a.Semester, &a.SGPA)

		if a.Semester == 1 {
			a.CGPA = a.SGPA
		} else {
			a.CGPA = (list_marks[len(list_marks)-1].CGPA*float64(a.Semester-1) + a.SGPA) / float64(a.Semester)
		}

		if err != nil {
			return list_marks, err
		}
		list_marks = append(list_marks, a)
	}
	if err = rows.Err(); err != nil {
		return list_marks, err
	}

	return list_marks, nil
}

func GetPrevDonation() ([]Donation, error) {
	rows, err := db.Query(getPrevDonations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donations []Donation

	for rows.Next() {
		var d Donation
		err := rows.Scan(&d.Name, &d.Amount)

		if err != nil {
			return donations, err
		}
		donations = append(donations, d)
	}
	if err = rows.Err(); err != nil {
		return donations, err
	}

	return donations, nil
}

func GetTotalDonation(username string) (int64, error) {
	row := db.QueryRow(getTotalDonationsByAlum, username)
	if row.Err() != nil {
		log.Panic(row.Err())
	}

	var t int64

	err := row.Scan(&t)

	if err != nil {
		if err == sql.ErrNoRows {
			//log.Println("No user found with the given ID.")
		} else {
			log.Println(err)
		}
		return 0, errors.New("No donations")
	}
	return t, nil
}

func AddDonation(username string, amount float64, message sql.NullString) (int64, error) {
	result, err := db.Exec(addDonation, username, amount, message)

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id, err
}
func GetEmploymentHistory(username string) ([]Employment, error) {
	rows, err := db.Query(getEmploymentHistory, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employments []Employment

	for rows.Next() {
		var e Employment
		err := rows.Scan(&e.Id, &e.Roll_no, &e.Starting, &e.Ending, &e.Company, &e.Designation, &e.Location)

		if err != nil {
			return employments, err
		}
		employments = append(employments, e)
	}
	if err = rows.Err(); err != nil {
		return employments, err
	}

	return employments, nil
}

func AddEmploymentHistory(e Employment) error {
	_, err := db.Exec(addEmploymentHistory, e.Roll_no, e.Starting, e.Ending, e.Company, e.Designation, e.Location)

	return err
}

func UpdateEmploymentHistory(e Employment) error {
	_, err := db.Exec(updateEmploymentHistory, e.Id, e.Roll_no, e.Starting, e.Ending, e.Company, e.Designation, e.Location)

	return err
}

func DeleteEmploymentHistory(id int) error {
	_, err := db.Exec(deleteEmploymentHistory, id)

	return err
}
func GetAllAnnouncements(username string) ([]Announcemnt, error) {
	rows, err := db.Query(getNoticesAndEvents, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var announcements []Announcemnt

	for rows.Next() {
		var t time.Time
		var a Announcemnt
		err := rows.Scan(&a.Id, &a.Title, &a.Description, &a.Date, &a.Location, &a.Type, &a.RSVP, &t)

		a.Created = t.Format("2006-01-02")
		if err != nil {
			return announcements, err
		}
		announcements = append(announcements, a)
	}
	if err = rows.Err(); err != nil {
		return announcements, err
	}

	return announcements, nil
}
func AddRVSP(eventid int, username string, response string) error {
	_, err := db.Exec(addRSVP, eventid, username, response)
	return err
}

func GarbageCollector() {
	command := "DELETE FROM sessions where expiry < ?"
	for {
		time.Sleep(30 * time.Minute)

		result, err := db.Exec(command, time.Now())
		if err != nil {
			log.Panic(err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Rows affected by delete: %d\n", rowsAffected)
	}
}
