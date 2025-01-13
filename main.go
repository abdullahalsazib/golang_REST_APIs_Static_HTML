package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var db *gorm.DB
var err error

var PORT = ":8080"

func init() {
	dns := "root:1234@tcp(localhost:3306)/userdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}
	db.AutoMigrate(User{})
	fmt.Print("Connected to database successfully")
}

func main() {
	fmt.Println("Hello, World")
	r := mux.NewRouter()

	// routers
	r.HandleFunc("/api/users", CreateUser).Methods("POST")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(PORT, r))

}

// CreateUsers
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
	}
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Faild to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
