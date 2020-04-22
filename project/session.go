package main

import (
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func getUser(w http.ResponseWriter, req *http.Request) user {
	// get cookie
	cookie, err := req.Cookie("session")
	//If there is no session cookie, create a new session cookie
	if err != nil {
		sID, _ := uuid.NewV4() //Give sID a random number
		cookie = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

	}
	http.SetCookie(w, cookie) //Set the cookie to our grabbed cookie,(or new cookie)

	// if the user exists already, get user
	var theUser user
	if un, ok := dbSessions[cookie.Value]; ok {
		theUser = dbUsers[un]
	}
	return theUser
}

func alreadyLoggedIn(req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false //If there is an error getting the cookie, return false
	}
	/* We assign the cookie.Value in our dbSessions to username.
	If we find that username in dbUsers, then they are logged in and we will
	return true! */
	username := dbSessions[cookie.Value]
	_, ok := dbUsers[username]
	return ok
}

//Needed this to test bitcrypt
func testDebug(w http.ResponseWriter, req *http.Request) {
	password := "apassword"
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(bs)
}
