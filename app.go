package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "postgres"
	dbname   = "accounts"
)

type Users struct {
	Username string
	Password string
}

func index(response http.ResponseWriter, request *http.Request) {
	temp, _ := template.ParseFiles("templates/index.html")
	temp.Execute(response, nil)
}

func register(response http.ResponseWriter, request *http.Request) {
	temp, _ := template.ParseFiles("templates/register.html")
	temp.Execute(response, nil)
}

func confirm(response http.ResponseWriter, request *http.Request) {
	db := connect()
	temp, _ := template.ParseFiles("templates/confirm.html")
	var query string
	user := Users{}
	user.Username = request.FormValue("name")
	user.Password = request.FormValue("pw")
	if uniqueName(db, user.Username) == false || len(user.Username) < 3 {
		db.Close()
		if len(user.Username) < 3 {
			temp, _ := template.ParseFiles("templates/nametooshort.html")
			temp.Execute(response, nil)

		} else {
			temp, _ = template.ParseFiles("templates/notunique.html")
			temp.Execute(response, nil)
		}
		return
	}
	if len(user.Password) < 3 {
		db.Close()
		temp, _ := template.ParseFiles("templates/pwtooshort.html")
		temp.Execute(response, nil)
		return
	}
	query = "INSERT INTO users (username, password, status)"
	query += " VALUES ($1, $2, $3)"
	db.QueryRow(query, user.Username, user.Password, "notapplied")
	defer db.Close()
	temp.Execute(response, user)
}

func uniqueName(db *sql.DB, name string) bool {
	rows, _ := db.Query("select username from users")
	for rows.Next() {
		var username string
		rows.Scan(&username)
		if name == username {
			return false
		}
	}
	return true
}

func connect() *sql.DB {
	var conn string
	conn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully connected to database")
	return db
}

func main() {
	fmt.Println("Hello, go!")

	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/confirm", confirm)
	http.ListenAndServe(":7000", nil)
}
