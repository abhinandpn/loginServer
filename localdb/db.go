package localdb

import (
	"time"
)

var Home = "home.html"
var HomePath = "pages/home.html"
var Login = "login.html"
var LoginPath = "pages/login.html"
var Register = "register.html"
var RegisterPath = "pages/register.html"

var ErrorPage = "errorPage.html"
var ErrorPagePath = "pages/errorPage.html"

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

func (s Session) IsSessionExpire() bool {
	return s.Expire.Before(time.Now())
}

// session db
var SessionsDB = make(map[string]Session)
