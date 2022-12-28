package localdb

import (
	"time"
)

// to strore single user private data
type UserDetals struct {
	Name  string
	Email string // email in key but if need email can get from this
	Pass  string
}

// using map key as email and struct ass value
var DataBase = make(map[string]UserDetals)

//session

type Session struct {
	Name   string
	Expire time.Time
}

func (s Session) IsSessionExpired() bool {
	return s.Expire.Before(time.Now())
}

// session db
var SessionsDB = make(map[string]Session)

//

//message struct for login form

type Messages struct {
	Color   string
	Message string
}

// this have the message what we want to show in login
var LoginMessage = Messages{}

// struct to store all eroors what user made when user registering using boolean value
// if a fiels is flase means its no error otherwise there is an error
type regFormErrors struct {
	ErrorName  bool
	ErrorEmail bool
	ErrorPass  bool
}

var RegError = regFormErrors{}

// varables

// type Templatess struct {
// 	Name string
// 	Path string
// }

var Home = "home.html"
var HomePath = "templates/home.html"
var Login = "login.html"
var LoginPath = "templates/login.html"
var Register = "register.html"
var RegisterPath = "templates/register.html"

var CookieId string

//varables

// if error get when parsing templte use this templates
var ErrorPage = "errorPage.html"
var ErrorPagePath = "templates/errorPage.html"
