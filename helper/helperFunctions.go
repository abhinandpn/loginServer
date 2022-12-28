package helper

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/nikhilnarayanan623/loginServer/localdb"
)

// these funtions are used to help other function
func CreateTemplate(name string, path string) *template.Template {

	tmpl, err := template.New(name).ParseFiles(path)

	if CheckError(err, "temlpate "+name) { //to check error
		//if found error parse error page
		errTmpl, _ := template.New(localdb.ErrorPage).ParseFiles(localdb.ErrorPagePath)
		return errTmpl
	}

	return tmpl
}

func CheckError(err error, name string) bool {
	if err != nil {
		fmt.Println("Error found at ", name)
		return true
	}

	return false
}

func SessionAndCookie(r *http.Request) (localdb.Session, bool) {

	if cookieVal, ok1 := GetCookieVal(r); ok1 { //get coookie val if cookie not get return false

		if session, ok2 := localdb.SessionsDB[cookieVal]; ok2 { //get session using cookie value other return flase and empty session

			if !session.IsSessionExpired() { //check session is expired is expired then delete otherwise return that session

				return session, true
			}

			//delte sessoin if session is expired

			delete(localdb.SessionsDB, cookieVal)
		}
	}

	return localdb.Session{}, false //return nill session
}

// get cookie if need
func GetCookieVal(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session") // session is the cookie name

	if CheckError(err, " getting Cokkie") {
		fmt.Println("session and cokkie func error to get cookie")
		return "", false
	}
	localdb.CookieId = cookie.Value
	return cookie.Value, true
}
