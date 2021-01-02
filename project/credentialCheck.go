package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func canLogin(w http.ResponseWriter, r *http.Request) {
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Declare DataType from Ajax
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	//Declare Response back to Ajax
	type Response struct {
		ResultNum     int    `json:"ResultNum"`
		ResultMessage string `json:"ResultMessage"`
	}

	responseMessage := Response{}

	//Marshal the user data into our type
	var dataForLogin LoginData
	json.Unmarshal(bs, &dataForLogin)

	//Check to see if the login is legit
	//Query database for those username and password
	rows, err := db.Query(`SELECT * FROM users WHERE USERNAME = ?;`, dataForLogin.Username)
	check(err)
	defer rows.Close()

	//Count to see if password is found or not
	var returnedTableID int = 0
	var returnedUsername string = ""
	var returnedPassword string = ""
	var returnedFName string = ""
	var returnedLName string = ""
	var returnedRole string = ""
	var returnedUserID int = 0
	var returnedDateCreated string = ""
	var returnedDateUpdated string = ""

	for rows.Next() {
		//assign variable
		err = rows.Scan(&returnedTableID, &returnedUsername, &returnedPassword, &returnedFName, &returnedLName, &returnedRole, &returnedUserID,
			&returnedDateCreated, &returnedDateUpdated)
		check(err)
	}

	//Count to see if password/Username returned at all
	if (strings.Compare(returnedUsername, "") == 0) || (strings.Compare(returnedPassword, "") == 0) {
		//Username or password not found, return failure
		responseMessage.ResultMessage = "Username and/or password not found!"
		responseMessage.ResultNum = 1
		theJSONMessage, err := json.Marshal(responseMessage)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
		}
		fmt.Fprint(w, string(theJSONMessage))
		return
	} else {
		//Check to see if Username returned
		if strings.Compare(dataForLogin.Username, returnedUsername) == 0 {
			//Checking to see if password matches as well
			theReturnedByte, err := hex.DecodeString(returnedPassword)
			if err != nil {
				log.Fatal(err)
				logWriter(err.Error())
			}
			if strings.Compare(string(theReturnedByte), dataForLogin.Password) != 0 {
				//Password not found/not hashed correctly
				responseMessage.ResultMessage = "Username and/or password do not match!"
				responseMessage.ResultNum = 1
				theJSONMessage, err := json.Marshal(responseMessage)
				if err != nil {
					fmt.Println(err)
					logWriter(err.Error())
				}
				fmt.Fprint(w, string(theJSONMessage))
				return
			} else {
				//Username matched, password matched good stuff
				//User logged in, directing them to the mainpage
				//Going to next, passing values
				theUser := User{dataForLogin.Username, returnedPassword, returnedFName, returnedLName, returnedRole, returnedUserID,
					returnedDateCreated, returnedDateUpdated}
				dbUsers[dataForLogin.Username] = theUser
				// create session
				uuidWithHyphen := uuid.New().String()

				cookie := &http.Cookie{
					Name:  "session",
					Value: uuidWithHyphen,
				}
				cookie.MaxAge = sessionLength
				http.SetCookie(w, cookie)
				dbSessions[cookie.Value] = theSession{dataForLogin.Username, time.Now()}
				/* values inserted, write back to ajax so we can
				go to the choice page */
				responseMessage.ResultMessage = "User found!"
				responseMessage.ResultNum = 0
				theJSONMessage, err := json.Marshal(responseMessage)
				if err != nil {
					fmt.Println(err)
					logWriter(err.Error())
				}
				fmt.Fprint(w, string(theJSONMessage))
				return
			}
		} else {
			//Username/Password do not match
			responseMessage.ResultMessage = "Username and/or password do not match!"
			responseMessage.ResultNum = 1
			theJSONMessage, err := json.Marshal(responseMessage)
			if err != nil {
				fmt.Println(err)
				logWriter(err.Error())
			}
			fmt.Fprint(w, string(theJSONMessage))
			return
		}
	}
}
