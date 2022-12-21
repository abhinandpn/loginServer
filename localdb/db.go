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

func (s Session) IsSessionExpire() bool {
	return s.Expire.Before(time.Now())
}

// session db
var SessionsDB = make(map[string]Session)
