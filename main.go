package main

import (
	"fmt"
	"net/http"

	"github.com/nikhilnaryanan623/loginServer/functions"
	"github.com/nikhilnaryanan623/loginServer/localdb"

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

	http.ListenAndServe(port, r)
}
