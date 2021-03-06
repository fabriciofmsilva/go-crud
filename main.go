package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Names represent a name struct
type Names struct {
	ID    int
	Name  string
	Email string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "crudgo"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("tmpl/*"))

// Index function
func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM names ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	n := Names{}

	res := []Names{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.ID = id
		n.Name = name
		n.Email = email

		res = append(res, n)
	}

	tmpl.ExecuteTemplate(w, "Index", res)

	defer db.Close()
}

// Show function
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nID := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nID)
	if err != nil {
		panic(err.Error())
	}

	n := Names{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.ID = id
		n.Name = name
		n.Email = email
	}

	tmpl.ExecuteTemplate(w, "Show", n)

	defer db.Close()
}

// New function
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Edit function
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	n := Names{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.ID = id
		n.Name = name
		n.Email = email
	}

	tmpl.ExecuteTemplate(w, "Edit", n)

	defer db.Close()
}

// Insert function
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES (?, ?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email)

		log.Println("INSERT: Name: " + name + " | Email: " + email)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

// Update function
func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email, id)

		log.Println("UPDATE: Name: " + name + " | Email: " + email)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

// Delete function
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nID := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(nID)

	log.Println("DELETE")

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:9000")

	// URLs management
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	// Actions
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	// Server start
	http.ListenAndServe(":9000", nil)
}
