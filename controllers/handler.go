package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/loginServer/helper"
	"github.com/nikhilnarayanan623/loginServer/localdb"
)

// these functions are hndler functions

func RegisterPage(w http.ResponseWriter, r *http.Request) {

	if _, ok := helper.SessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	regTmpl := helper.CreateTemplate(localdb.Register, localdb.RegisterPath)
	regTmpl.Execute(w, localdb.RegError)
}

// register submit
func RegisterSubmit(w http.ResponseWriter, r *http.Request) {

	if _, ok := helper.SessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	//get values from fom
	fName := r.PostFormValue("name")   //form name
	fEmail := r.PostFormValue("email") //form mail
	fPass1 := r.PostFormValue("fpass") //form first pass
	fPass2 := r.PostFormValue("spass") // form second pass

	// validate the user value
	if fName == "" {
		localdb.RegError.ErrorName = true
	}
	if fEmail == "" {
		localdb.RegError.ErrorEmail = true
	}
	if fPass1 == "" || fPass2 == "" || fPass1 != fPass2 {
		localdb.RegError.ErrorPass = true
	}

	//check the entered user is alredy exist
	if _, ok := localdb.DataBase[fEmail]; ok {

		localdb.LoginMessage.Color = "text-success"
		localdb.LoginMessage.Message = "You are already a User"

		http.Redirect(w, r, "/", http.StatusSeeOther) //render login page to show the message
		return
	}

	//check anyting in regError is have a error as true
	//then render same page with showing this errors
	if localdb.RegError.ErrorEmail || localdb.RegError.ErrorName || localdb.RegError.ErrorPass {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	//if all validation is completed add value to localDB and show the login page
	localdb.DataBase[fEmail] = localdb.UserDetals{
		Name:  fName,
		Email: fEmail,
		Pass:  fPass1,
	}

	//set login text class and message
	localdb.LoginMessage.Color = "text-success"
	localdb.LoginMessage.Message = "Successfully Registered Login Please"

	clearRegError()
	//redirect to login page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func clearRegError() {
	localdb.RegError.ErrorEmail = false
	localdb.RegError.ErrorName = false
	localdb.RegError.ErrorPass = false
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginPage")

	if _, ok := helper.SessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	//session is not avaliable then render login page
	clearRegError() //to clear register form errors

	tmpl := helper.CreateTemplate(localdb.Login, localdb.LoginPath)
	tmpl.Execute(w, localdb.LoginMessage)

	//clear all login error messages
	localdb.LoginMessage.Color = ""
	localdb.LoginMessage.Message = ""
}

func LoginSubmit(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Login submit start")

	if _, ok := helper.SessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	userEmail := r.FormValue("email")
	userPass := r.FormValue("pass")

	//check user entered email
	if userEmail == "" {
		//setting
		localdb.LoginMessage.Color = "text-danger"
		localdb.LoginMessage.Message = "Enter Email Properly"
		//after setting message call login handler to render login page
		LoginPage(w, r)
		return
	}

	// check this email contains in our mapDB
	singleUser, ok := localdb.DataBase[userEmail]

	//check user is exist or not then check password
	if !ok {
		localdb.LoginMessage.Color = "text-danger"
		localdb.LoginMessage.Message = "You are not a registered User! you can register"

		LoginPage(w, r)
		return
	} else if userPass != singleUser.Pass { // user exist password not match
		localdb.LoginMessage.Color = "text-danger"
		localdb.LoginMessage.Message = "Incorrect Password"

		LoginPage(w, r)
		return
	}
	//create session
	sessionToken := uuid.NewString() //create a new random session token

	//sessionToken := "token"
	sessionTime := time.Now().Add(3 * time.Minute) //expire time current time plus two minute

	newSession := localdb.Session{
		Name:   singleUser.Name,
		Expire: sessionTime,
	}

	//add this sessoin to session database
	localdb.SessionsDB[sessionToken] = newSession

	//set cookie
	newCookie := &http.Cookie{
		Name:    "session",
		Value:   sessionToken,
		Expires: sessionTime,
	}
	http.SetCookie(w, newCookie)

	HomePage(w, r)
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	session, ok := helper.SessionAndCookie(r)

	if !ok { //if session no availabe

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	homeTmpl := helper.CreateTemplate(localdb.Home, localdb.HomePath)
	homeTmpl.Execute(w, session)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout page ")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	_, ok := helper.SessionAndCookie(r)

	if ok {
		delete(localdb.SessionsDB, localdb.CookieId)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// this function check the user is entering invalid url from login page or home according to that function redirect to that page
func ErrorHandleFunc(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if _, ok := helper.SessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
