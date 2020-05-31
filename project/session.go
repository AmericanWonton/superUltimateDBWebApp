package main

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func getUser(w http.ResponseWriter, req *http.Request) User {
	// get cookie
	cookie, err := req.Cookie("session")
	//If there is no session cookie, create a new session cookie
	if err != nil {
		sID, err := uuid.NewV4() //Give sID a random number
		cookie = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		if err != nil {
			fmt.Println(err)
		}

	}
	//Set the cookie age to the max length again.
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie) //Set the cookie to our grabbed cookie,(or new cookie)

	// if the user exists already, get user
	var theUser User
	if session, ok := dbSessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
		theUser = dbUsers[session.username]
	}
	return theUser
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false //If there is an error getting the cookie, return false
	}
	//if session is found, we update the session with the newest time since activity!
	session, ok := dbSessions[cookie.Value]
	if ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
	}
	/* Check to see if the Username exists from this Session Username. If not, we return false. */
	_, ok = dbUsers[session.username]
	// refresh session
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie)
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
