package logger

import (
	"bytes"
	"fmt"
)

type User struct {
	Grupo, User, Pass [10]byte
	Id                string
}

type Logger struct {
	LoggedIn bool
	Usr      *User
}

func (l *Logger) Login(usr *User) bool {
	if !l.IsLoggedIn() {
		l.Usr = usr
		return true
	}
	fmt.Println("ya existe un usuario registrado")
	return false
}

func (l *Logger) IsLoggedIn() bool {
	return l.Usr != nil
}

func (l *Logger) Logout() bool {
	if !l.IsLoggedIn() {
		l.Usr = nil
	}
	fmt.Println("no habia un usuario loggeado")
	return false
}

func (l *Logger) UserIsRoot() bool {
	return bytes.Equal(l.Usr.User[:], []byte("root"))
}

func (l *Logger) GetUserName() [10]byte {
	return l.Usr.User
}

var Log = Logger{
	LoggedIn: false,
	Usr:      nil,
}
