package main

import (
	"fmt"
	"log"
	"loginPage/functions"
	"loginPage/localdb"
	"net/http"

	"github.com/gorilla/mux"
)

var port = ":8000"

func main() {

	localdb.DataBase["admin@gmail.com"] = localdb.UserDetals{
		Name:  "Nikhil",
		Email: "admin@gamil.com",
		Pass:  "123",
	}

	r := mux.NewRouter()
	r.HandleFunc("/", functions.LoginPage)
	r.HandleFunc("/submit", functions.LoginSubmit)
	r.HandleFunc("/home", functions.HomePage)
	r.HandleFunc("/register", functions.RegisterPage)
	r.HandleFunc("/registerSubmit", functions.RegisterSubmit)
	r.HandleFunc("/logout", functions.Logout)

	r.NotFoundHandler = http.HandlerFunc(functions.ErrorHandleFunc)

	fmt.Println("Server running at ", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
