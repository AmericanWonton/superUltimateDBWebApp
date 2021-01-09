package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// WARNING! THIS CODE CONTAINS DEROGATORY TERMS, RACIAL/ETHNIC/SEXUAL SLURS,
// AND OTHER OFFENSIVE CONTENT. THE PURPOSE IS TO REMOVE THIS CONTENT OFF OF
// MY PLATFORM. IF ANY OF THIS CONTENT OFFENDS YOU, I APOLOGIZE; PLEASE STAY OFF
// OF THIS PAGE!!!

var allUsernames []string
var usernameMap map[string]bool

/* DEFINED SLURS */
var slurs []string = []string{"penis", "vagina", "dick", "cunt", "asshole", "fag", "faggot",
	"nigglet", "nigger", "beaner", "wetback", "wet back", "chink", "tranny", "bitch", "slut",
	"whore", "fuck", "damn",
	"shit", "piss", "cum", "jizz"}

func containsLanguage(theText string) bool {
	hasLanguage := false
	textLower := strings.ToLower(theText)
	for i := 0; i < len(slurs); i++ {
		if strings.Contains(textLower, slurs[i]) {
			hasLanguage = true
			return hasLanguage
		}
	}
	return hasLanguage
}

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
	} else if containsLanguage(sbs) {
		fmt.Fprint(w, "ContainsLanguage")
	} else {
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

//Begins Sending Email to User and creates a User for database entry
func signUpUserUpdated(w http.ResponseWriter, req *http.Request) {
	// process Ajax ping
	if req.Method == http.MethodPost {
		//Collect JSON from Postman or wherever
		//Get the byte slice from the request body ajax
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
		}

		//Declare DataType from Ajax
		type UserData struct {
			TheUser User   `json:"TheUser"`
			Email   string `json:"Email"`
		}

		//Marshal the user data into our type
		var dataPosted UserData
		json.Unmarshal(bs, &dataPosted)
		//Set the User info
		var postedUser User = dataPosted.TheUser
		// get form values
		username := postedUser.UserName
		password := postedUser.Password
		firstname := postedUser.First
		lastname := postedUser.Last
		role := postedUser.Role
		email := dataPosted.Email
		/* ATTEMPT TO SEND EMAIL...IF IT FAILS, DO NOT CREATE USER */
		goodEmailSend := signUpUserEmail(email, role, firstname, lastname)
		if goodEmailSend == true {
			// create session
			uuidWithHyphen := uuid.New().String()
			newCookie := &http.Cookie{
				Name:  "session",
				Value: uuidWithHyphen,
			}
			newCookie.MaxAge = sessionLength
			http.SetCookie(w, newCookie)
			dbSessions[newCookie.Value] = theSession{username, time.Now()}
			// store user in dbUsers
			//Make User and USERID
			theID := randomIDCreation()

			//Add User to the SQL Database
			bsString := []byte(password)                  //Encode Password
			encodedString := hex.EncodeToString(bsString) //Encode Password Pt2
			theTimeNow := time.Now()
			var insertedUser User = User{
				UserName:    username,
				Password:    encodedString,
				First:       firstname,
				Last:        lastname,
				Role:        role,
				UserID:      theID,
				DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
				DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
			}
			jsonValue, _ := json.Marshal(insertedUser)
			response, err := http.Post("http://"+serverAddress+"/insertUser", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(data))
			}

			//Add User to MongoDB
			var insertionUser AUser = AUser{
				UserName:    username,
				Password:    encodedString,
				First:       firstname,
				Last:        lastname,
				Role:        role,
				UserID:      theID,
				DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
				DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
				Hotdogs:     MongoHotDogs{},
				Hamburgers:  MongoHamburgers{},
			}
			insertionUsers := TheUsers{
				Users: []AUser{insertionUser},
			}
			jsonValue2, _ := json.Marshal(insertionUsers)
			response2, err := http.Post("http://"+serverAddress+"/insertUsers", "application/json", bytes.NewBuffer(jsonValue2))
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response2.Body)
				fmt.Println(string(data))
			}
			//DEBUG, don't know if we need below
			var theUser = User{username, encodedString, firstname, lastname, role, theID, insertionUser.DateCreated,
				insertionUser.DateUpdated}
			dbUsers[username] = theUser
			//Alert Ajax with success

			type successMSG struct {
				Message     string `json:"Message"`
				SuccessNum  int    `json:"SuccessNum"`
				RedirectURL string `json:"RedirectURL"`
			}
			msgSuccess := successMSG{
				Message:     "Added the new account!",
				SuccessNum:  0,
				RedirectURL: "http://" + serverAddress,
			}

			theJSONMessage, err := json.Marshal(msgSuccess)
			if err != nil {
				fmt.Println(err)
				logWriter(err.Error())
			}

			fmt.Printf("DEBUG: Writing back successful User response.\n")
			fmt.Fprint(w, string(theJSONMessage))
		} else {
			//Alert Ajax with failure
			type successMSG struct {
				Message    string `json:"Message"`
				SuccessNum int    `json:"SuccessNum"`
			}
			msgSuccess := successMSG{
				Message:    "Failed to send email and create User.",
				SuccessNum: 1,
			}
			theJSONMessage, err := json.Marshal(msgSuccess)
			if err != nil {
				fmt.Println(err)
				logWriter(err.Error())
			}
			fmt.Fprint(w, string(theJSONMessage))
		}
	}
}
