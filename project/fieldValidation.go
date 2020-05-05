package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var allUsernames []string
var usernameMap map[string]bool

//Checks the Usernames after every keystroke
func checkUsername(w http.ResponseWriter, req *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	if len(sbs) <= 0 {
		fmt.Fprint(w, "TooShort")
	} else if len(sbs) > 20 {
		fmt.Fprint(w, "TooLong")
	} else {
		fmt.Println("USERNAME: ", sbs)
		fmt.Fprint(w, usernameMap[sbs])
	}
}

//Loads all our Usernames when the document loads.
func loadUsernames(w http.ResponseWriter, req *http.Request) {
	/* DEBUG NOTE: I SHOULD RE-WRITE THIS TO USE CHANNELS AT SOME POINT */
	usernameMap = make(map[string]bool) //Clear Map for future use on page load
	var grabbedUsername string          //The Username grabbed from the database
	//Query the database for all names
	row, err := db.Query(`SELECT username FROM users;`)
	check(err)
	defer row.Close()
	//Append each name to the next
	for row.Next() {
		err = row.Scan(&grabbedUsername)
		check(err)
		usernameMap[grabbedUsername] = true
	}

	if err != nil {
		fmt.Fprint(w, "false")
	} else {
		fmt.Fprint(w, "true")
	}
}
