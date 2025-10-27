package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var db *sql.DB

var users = map[string]string{"admin": "password", "lian": "3340"} //username : password
var cookies = map[string]string{"cook": "admin"}                   //cookies : username
var data = map[string]string{"admin": "secret", "di": "cookies"}   //username : secret data

func Connect(path string) {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
		//return err
	}
	createTables := `
	CREATE TABLE IF NOT EXISTS credentials (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		session_id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		expiry TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS user_data (
		username TEXT PRIMARY KEY,
		data TEXT
	);`
	_, err = db.Exec(createTables)
	if err != nil {
		log.Panic(err)
		//return err
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
	for username, password := range users {
		hash_pass, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if e != nil {
			log.Panic(e)
		}

		command := "insert or replace into credentials(username, password) values (?,?);"
		_, err := db.Exec(command, username, string(hash_pass))
		if err != nil {
			log.Panic(err)
		}
	}

	for username, data := range data {
		command := "insert or ignore into user_data(username, data) values (?,?);"
		_, err := db.Exec(command, username, data)
		if err != nil {
			log.Panic(err)
		}
	}
}

func Authenticate(username string, password string) bool {
	query := "SELECT password from credentials where username = ?"
	row := db.QueryRow(query, username)
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
	//log.Println(ret_pass)
	//log.Println(hash_pass)
	err := bcrypt.CompareHashAndPassword([]byte(hash_pass), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func CheckSession(session_id string) (string, error) {
	query := "Select username, expiry from sessions where session_id = ?;"
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
	_, err := db.Exec(command, session_id, username, time.Now().Add(2*time.Minute))
	if err != nil {
		log.Panic(err)
	}
}

func GetUserData(username string) (string, error) {
	query := "SELECT data from user_data where username = ?"
	row := db.QueryRow(query, username)
	if row.Err() != nil {
		log.Panic(row.Err())
	}

	var data string

	if err := row.Scan(&data); err != nil {
		if err == sql.ErrNoRows {
			//log.Println("No user found with the given ID.")
		} else {
			log.Println(err)
		}
		return "", errors.New("No User data found")
	}
	return data, nil
}

func GarbageCollector() {
	command := "DELETE FROM sessions where expiry < ?"
	for {
		time.Sleep(1 * time.Minute)

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
