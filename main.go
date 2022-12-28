package main

import (
	"fmt"
	"net/http"

	"github.com/nikhilnarayanan623/loginServer/controllers"
	"github.com/nikhilnarayanan623/loginServer/localdb"

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
	r.HandleFunc("/", controllers.LoginPage).Methods("GET")
	r.HandleFunc("/", controllers.LoginSubmit).Methods("POST")
	r.HandleFunc("/home", controllers.HomePage)
	r.HandleFunc("/register", controllers.RegisterPage).Methods("GET")
	r.HandleFunc("/register", controllers.RegisterSubmit).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout)
	//no handler found use these handler
	r.NotFoundHandler = http.HandlerFunc(controllers.ErrorHandleFunc)
	fmt.Println("Server running at ", port)

	http.ListenAndServe(port, r)
}
