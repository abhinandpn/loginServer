package functions

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/nikhilnaryanan623/loginServer/localdb"
)

// varables

var Home = "home.html"
var HomePath = "pages/home.html"
var Login = "login.html"
var LoginPath = "pages/login.html"
var Register = "register.html"
var RegisterPath = "pages/register.html"

var cookieId string

//varables

// if error get when parsing templte use this pages
var ErrorPage = "errorPage.html"
var ErrorPagePath = "pages/errorPage.html"

//message struct for login form

type Messages struct {
	Color   string
	Message string
}

// this have the message what we want to show in login
var loginMessage = Messages{}

// struct to store all eroors what user made when user registering using boolean value
// if a fiels is flase means its no error otherwise there is an error
type regFormErrors struct {
	ErrorName  bool
	ErrorEmail bool
	ErrorPass  bool
}

var regError = regFormErrors{}

// these functions are hndler functions

func RegisterPage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("regPage")

	if _, ok := sessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	regTmpl := createTemplate(Register, RegisterPath)

	regTmpl.Execute(w, regError)
}

// register submit
func RegisterSubmit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("regSub")

	if _, ok := sessionAndCookie(r); ok {
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
		regError.ErrorName = true
	}
	if fEmail == "" {
		regError.ErrorEmail = true
	}
	if fPass1 == "" || fPass2 == "" || fPass1 != fPass2 {
		regError.ErrorPass = true
	}

	//check the entered user is alredy exist

	if _, ok := localdb.DataBase[fEmail]; ok {

		loginMessage.Color = "text-success"
		loginMessage.Message = "You are already a User"
		clearRegError()
		http.Redirect(w, r, "/", http.StatusSeeOther) //render login page to show the message
		return
	}

	//check anyting in regError is have a error as true
	//then render same page with showing this errors
	if regError.ErrorEmail || regError.ErrorName || regError.ErrorPass {
		fmt.Println(regError, "reg error")
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
	loginMessage.Color = "text-success"
	loginMessage.Message = "Successfully Registered Login Please"

	clearRegError()
	//redirect to login page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func clearRegError() {
	//clear all error in regError
	regError.ErrorEmail = false
	regError.ErrorName = false
	regError.ErrorPass = false
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginPage")

	if _, ok := sessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	//session is not avaliable then render login page

	tmpl := createTemplate(Login, LoginPath)
	tmpl.Execute(w, loginMessage)

	//clear all login error messages
	loginMessage.Color = ""
	loginMessage.Message = ""
}

func LoginSubmit(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Login submit start")

	if _, ok := sessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	userEmail := r.FormValue("email")
	userPass := r.FormValue("pass")

	//check user entered email
	if userEmail == "" {
		//setting
		loginMessage.Color = "text-danger"
		loginMessage.Message = "Enter Email Properly"
		//after setting message call login handler to render login page
		LoginPage(w, r)
		return
	}

	// check this email contains in our mapDB
	singleUser, ok := localdb.DataBase[userEmail]

	//check user is exist or not then check password
	if !ok {
		loginMessage.Color = "text-danger"
		loginMessage.Message = "You are not a registered User! you can register"

		LoginPage(w, r)
		return
	} else if userPass != singleUser.Pass { // user exist password not match
		loginMessage.Color = "text-danger"
		loginMessage.Message = "Incorrect Password"

		LoginPage(w, r)
		return
	}
	//create session
	userName := singleUser.Name

	sessionToken := uuid.NewString() //create a new random session token

	//sessionToken := "token"
	sessionTime := time.Now().Add(2 * time.Hour) //expire time current time plus two minute

	newSession := localdb.Session{
		Name:   userName,
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
	//time.Sleep(2 * time.Second)

	HomePage(w, r)

}

func HomePage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("home page ")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	session, ok := sessionAndCookie(r)

	if !ok { //if session no availabe

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	homeTmpl := createTemplate(Home, HomePath)

	homeTmpl.Execute(w, session)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout page ")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	_, ok := sessionAndCookie(r)

	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//get cookie and delte session using that cookie
	// if cokkieVal, ok := getCookieVal(r); ok {
	// 	delete(localdb.SessionsDB, cokkieVal)
	// }

	delete(localdb.SessionsDB, cookieId)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	//LoginPage(w, r)
}

func ErrorHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("error page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if _, ok := sessionAndCookie(r); ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// these funtions are used to help other function
func createTemplate(name string, path string) *template.Template {

	tmpl, err := template.New(name).ParseFiles(path)

	if checkError(err, "tempate "+name) { //to check error
		//if found error parse error page
		errTmpl, _ := template.New(ErrorPage).ParseFiles(ErrorPagePath)
		return errTmpl
	}

	return tmpl
}

func checkError(err error, name string) bool {
	if err != nil {
		fmt.Println("Error found at ", name)
		return true
	}

	return false
}

func sessionAndCookie(r *http.Request) (localdb.Session, bool) {

	if cookieVal, ok1 := getCookieVal(r); ok1 {

		if session, ok2 := localdb.SessionsDB[cookieVal]; ok2 {
			return session, true
		}
	}

	return localdb.Session{}, false //return nill session
}

// get cookie if need
func getCookieVal(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session") // session is the cookie name

	if checkError(err, "Error at getting Cokkie") {
		fmt.Println("session and cokkie func error to get cookie")
		return "", false
	}
	cookieId = cookie.Value
	return cookie.Value, true
}
