package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Database connection string
const dsn = "root:1234@tcp(127.0.0.1:3306)/userdb?charset=utf8mb4&parseTime=True&loc=Local"

// Struct for a Record
type User struct {
	ID    int
	Name  string
	Email string
}

var db *sql.DB
var tmpl *template.Template

func init() {
	// Initialize templates
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	// Open database connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}
}

func main() {
	http.HandleFunc("/", listUsers)
	http.HandleFunc("/create", createUser)
	http.HandleFunc("/save", saveUser)
	http.HandleFunc("/edit", editUser)
	http.HandleFunc("/update", updateUser)
	http.HandleFunc("/delete", deleteUser)

	// Serve static files (CSS, JS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "Failed to read user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	tmpl.ExecuteTemplate(w, "index.html", users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "create.html", nil)
}

func saveUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")

		_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
		if err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func editUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl.ExecuteTemplate(w, "update.html", user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		_, err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", name, email, id)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
