package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const successMessage string = "Successful Insert"
const failureMessage string = "Unsuccessful Insert"

//POST hotdog, Mainpage
func insertHotDog(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting hotdog record.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var postedHotDog Hotdog
	json.Unmarshal(bs, &postedHotDog)

	//Protections for the hotdog name
	if strings.Compare(postedHotDog.HotDogType, "DEBUGTYPE") == 0 {
		postedHotDog.HotDogType = "NONE"
	}
	//Fill empty values
	if len(postedHotDog.DateCreated) < 1 {
		theTimeNow := time.Now()
		postedHotDog.DateCreated = theTimeNow.Format("2006-01-02 15:04:05")
	}
	if len(postedHotDog.DateUpdated) < 1 {
		theTimeNow := time.Now()
		postedHotDog.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
	}

	stmt, err := db.Prepare("INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

	r, err := stmt.Exec(postedHotDog.HotDogType, postedHotDog.Condiment, postedHotDog.Calories, postedHotDog.Name, postedHotDog.UserID,
		postedHotDog.FoodID, postedHotDog.DateCreated, postedHotDog.DateUpdated)
	check(err)

	n, err := r.RowsAffected()
	check(err)
	stmt.Close() //Close the SQL

	fmt.Printf("DEBUG: %v rows effected.\n", n)

	//Define the data to return
	type returnData struct {
		SuccessMsg     string `json:"SuccessMsg"`
		ReturnedHotDog Hotdog `json:"ReturnedHotDog"`
		SuccessBool    bool   `json:"SuccessBool"`
	}

	if err != nil {
		theReturnData := returnData{
			SuccessMsg:     failureMessage,
			ReturnedHotDog: postedHotDog,
			SuccessBool:    false,
		}
		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling.")
			logWriter("There's an error marshalling.")
		}
		fmt.Fprintf(w, string(dataJSON))
	} else {
		hDogMarshaled, err := json.Marshal(postedHotDog)
		if err != nil {
			fmt.Printf("Error with %v\n", hDogMarshaled)
		}
		hDogSuccessMSG := successMessage + string(hDogMarshaled)

		theReturnData := returnData{
			SuccessMsg:     hDogSuccessMSG,
			ReturnedHotDog: postedHotDog,
			SuccessBool:    true,
		}

		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling.")
			logWriter("There's an error marshalling.")
		}

		fmt.Fprintf(w, string(dataJSON))
	}
}

//GET HOTDOGS
func getHotDog(w http.ResponseWriter, req *http.Request) {
	//Get the string map of our variables from the request
	fmt.Println("Finding hotdog singular")
	//Collect JSON from Postman or wherever
	reqBody, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("Here's our body: \n%v\n", reqBody)
	//Marshal it into our type
	var postedHotDog Hotdog
	var hotDogSlice []Hotdog
	json.Unmarshal(reqBody, &postedHotDog)
	stmt := "SELECT * FROM hot_dogs WHERE NAME = ?"
	rows, err := db.Query(stmt, postedHotDog.Name)
	check(err)
	var id int64
	var dogType string
	var condiment string
	var calories int
	var hotdogName string
	var userID int
	count := 0
	for rows.Next() {
		err = rows.Scan(&id, &dogType, &condiment, &calories, &hotdogName, &userID)
		check(err)
		//Add the hotdog to the slice list
		returnedHotDog := Hotdog{
			HotDogType: dogType,
			Condiment:  condiment,
			Calories:   calories,
			Name:       hotdogName,
			UserID:     userID,
		}
		hotDogSlice = append(hotDogSlice, returnedHotDog)
		count++
	}
	rows.Close()
	//If nothing returned from the rows
	if count == 0 {
		fmt.Fprint(w, "Nothing returned for this query.")
		return
	} else {
		//Marshal our return message to JSON
		hDogsMarshaled, err := json.Marshal(hotDogSlice)
		if err != nil {
			fmt.Printf("Error with %v\n", hDogsMarshaled)
			fmt.Fprint(w, "Error returned for this marshalling.")
		}
		fmt.Fprint(w, hDogsMarshaled)
	}
}

//DELETE hotdog
func deleteFood(w http.ResponseWriter, req *http.Request) {
	type foodDeletion struct {
		FoodType string `json:"FoodType"`
		FoodID   int    `json:"FoodID"`
		UserID   int    `json:"UserID"`
	}
	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var theFoodDeletion foodDeletion
	json.Unmarshal(bs, &theFoodDeletion)

	//Determine if this is a hotdog or hamburger deletion
	sqlStatement := ""
	if theFoodDeletion.FoodType == "hotdog" {
		sqlStatement = "DELETE FROM hot_dogs WHERE FOOD_ID=? AND USER_ID=?"
		delDog, err := db.Prepare(sqlStatement)
		check(err)

		r, err := delDog.Exec(theFoodDeletion.FoodID, theFoodDeletion.UserID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("%v\n", n)

		fmt.Fprintln(w, 1)
	} else if theFoodDeletion.FoodType == "hamburger" {
		sqlStatement = "DELETE FROM hamburgers WHERE FOOD_ID=? AND USER_ID=?"
		delDog, err := db.Prepare(sqlStatement)
		check(err)

		r, err := delDog.Exec(theFoodDeletion.FoodID, theFoodDeletion.UserID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("%v\n", n)

		fmt.Fprintln(w, 2)
	} else {
		fmt.Fprintln(w, 3)
	}
}

//INSERT HOTDOG
func insertHamburger(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting hamburger record.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var postedHamburger Hamburger
	json.Unmarshal(bs, &postedHamburger)

	//Protections for the hamburger name
	if strings.Compare(postedHamburger.BurgerType, "DEBUGTYPE") == 0 {
		postedHamburger.BurgerType = "NONE"
	}
	//Fill empty values
	if postedHamburger.FoodID == 0 || postedHamburger.FoodID < 8 {
		postedHamburger.FoodID = randomIDCreation()
	}
	if len(postedHamburger.DateCreated) < 1 {
		theTimeNow := time.Now()
		postedHamburger.DateCreated = theTimeNow.Format("2006-01-02 15:04:05")
	}
	if len(postedHamburger.DateUpdated) < 1 {
		theTimeNow := time.Now()
		postedHamburger.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
	}

	stmt, err := db.Prepare("INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

	r, err := stmt.Exec(postedHamburger.BurgerType, postedHamburger.Condiment,
		postedHamburger.Calories, postedHamburger.Name, postedHamburger.UserID, postedHamburger.FoodID, postedHamburger.DateCreated, postedHamburger.DateUpdated)
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("DEBUG: %v rows effected.\n", n)

	stmt.Close() //Close the SQL

	if err != nil {
		fmt.Fprint(w, failureMessage)
	} else {
		hamMarshaled, err := json.Marshal(postedHamburger)
		if err != nil {
			fmt.Printf("Error with %v\n", hamMarshaled)
		}
		hamSuccessMSG := successMessage + string(hamMarshaled)
		fmt.Fprint(w, hamSuccessMSG)
	}
}

//GET all Food
func getAllFoodUser(w http.ResponseWriter, req *http.Request) {
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	//Marshal it into our type
	var theUser User
	json.Unmarshal(bs, &theUser)

	//Declare variables for hotdog
	var h_id int
	var h_dogType string
	var h_condiment string
	var h_calories int
	var h_hotdogName string
	var h_userID int
	var h_foodID int
	var h_dateCreated int
	var h_dateUpdated int

	//Declare variables for hamburger
	var ham_id int
	var ham_type string
	var ham_condiment string
	var ham_calories int
	var ham_name string
	var ham_userID int
	var ham_foodID int
	var ham_dateCreated int
	var ham_dateUpdated int

	//Declare all our our hotdog/hamburger collections returned
	var hotDogSlice []Hotdog
	var hamburgerSlice []Hamburger
	var hotDogIDSlice []int
	var hamburgIDSlice []int
	var hDogMongoIDS []int
	var hamMonogIDS []int

	//Assemble data to send back
	type data struct {
		SuccessMessage string      `json:"SuccessMessage"`
		TheHotDogs     []Hotdog    `json:"TheHotDogs"`
		TheHamburgers  []Hamburger `json:"TheHamburgers:`
		ID_HotDogs     []int       `json:"ID_HotDogs"`
		ID_Hamburgers  []int       `json:"ID_Hamburgers"`
		HDogFoodIDS    []int       `json:"HDogFoodIDS"`
		HamFoodIDS     []int       `json:"HamFoodIDS"`
	}

	//Counter for food returned
	dogCounter := 0
	hamCounter := 0

	//If no User ID is submitted, then just food for ALL Users
	if theUser.UserID == 0 {
		//Get HotDogs
		hrows, err1 := db.Query(`SELECT * FROM hot_dogs ORDER BY ID;`)
		check(err1)

		for hrows.Next() {
			err = hrows.Scan(&h_id, &h_dogType, &h_condiment, &h_calories, &h_hotdogName, &h_userID,
				&h_foodID, &h_dateCreated, &h_dateUpdated)
			check(err) //Check to make sure there was no error doing that above.
			//Add Hotdog to a new Hotdog and add to slice
			var newHotDog Hotdog = Hotdog{
				HotDogType: h_dogType,
				Condiment:  h_condiment,
				Calories:   h_calories,
				Name:       h_hotdogName,
				UserID:     h_userID,
			}
			hotDogSlice = append(hotDogSlice, newHotDog)
			hotDogIDSlice = append(hotDogIDSlice, h_id)

			dogCounter = dogCounter + 1
		}

		//Get Mongo Food IDS for Hot Dogs
		hDogMongoIDS = getFoodIDSHDog(theUser.UserID)

		//Get Hamburgers
		hamrows, err2 := db.Query(`SELECT * FROM hamburgers ORDER BY ID`)
		check(err2)

		for hamrows.Next() {
			err = hamrows.Scan(&ham_id, &ham_type, &ham_condiment, &ham_calories, &ham_name, &ham_userID,
				&ham_foodID, &ham_dateCreated, &ham_dateUpdated)
			check(err) //Check to make sure there was no error doing that above.
			//Add Hamburger to a new Hamburger and add to slice
			var newHamburger Hamburger = Hamburger{
				BurgerType: ham_type,
				Condiment:  ham_condiment,
				Calories:   ham_calories,
				Name:       ham_name,
				UserID:     ham_userID,
			}
			hamburgerSlice = append(hamburgerSlice, newHamburger)
			hamburgIDSlice = append(hamburgIDSlice, ham_id)
			hamCounter = hamCounter + 1
		}

		//Get Mongo Food IDS for Hamburgers
		hamMonogIDS = getFoodIDSHam(theUser.UserID)
		//Close Connections
		hrows.Close()
		hamrows.Close()
	} else {
		//Get HotDogs
		hrows, err1 := db.Query(`SELECT * FROM hot_dogs WHERE USER_ID=? ORDER BY ID;`, theUser.UserID)
		check(err1)

		for hrows.Next() {
			err = hrows.Scan(&h_id, &h_dogType, &h_condiment, &h_calories, &h_hotdogName, &h_userID,
				&h_foodID, &h_dateCreated, &h_dateUpdated)
			check(err) //Check to make sure there was no error doing that above.
			//Add Hotdog to a new Hotdog and add to slice
			var newHotDog Hotdog = Hotdog{
				HotDogType: h_dogType,
				Condiment:  h_condiment,
				Calories:   h_calories,
				Name:       h_hotdogName,
				UserID:     h_userID,
			}
			hotDogSlice = append(hotDogSlice, newHotDog)
			hotDogIDSlice = append(hotDogIDSlice, h_id)

			dogCounter = dogCounter + 1
		}

		//Get Hamburgers
		hamrows, err2 := db.Query(`SELECT * FROM hamburgers WHERE USER_ID=? ORDER BY ID`, theUser.UserID)
		check(err2)

		for hamrows.Next() {
			err = hamrows.Scan(&ham_id, &ham_type, &ham_condiment, &ham_calories, &ham_name, &ham_userID,
				&ham_foodID, &ham_dateCreated, &ham_dateUpdated)
			check(err) //Check to make sure there was no error doing that above.
			//Add Hamburger to a new Hamburger and add to slice
			var newHamburger Hamburger = Hamburger{
				BurgerType: ham_type,
				Condiment:  ham_condiment,
				Calories:   ham_calories,
				Name:       ham_name,
				UserID:     ham_userID,
			}
			hamburgerSlice = append(hamburgerSlice, newHamburger)
			hamburgIDSlice = append(hamburgIDSlice, ham_id)
			hamCounter = hamCounter + 1
		}
		//Close Connections
		hrows.Close()
		hamrows.Close()
	}

	//Check to see if we have any data to submit
	sendData := data{
		SuccessMessage: "Success",
		TheHotDogs:     hotDogSlice,
		TheHamburgers:  hamburgerSlice,
		ID_HotDogs:     hotDogIDSlice,
		ID_Hamburgers:  hamburgIDSlice,
		HDogFoodIDS:    hDogMongoIDS,
		HamFoodIDS:     hamMonogIDS,
	}

	if len(sendData.TheHotDogs) <= 0 && len(sendData.TheHamburgers) <= 0 {
		sendData.SuccessMessage = "Failure"
	} else if len(sendData.ID_HotDogs) == 0 {
		hotDogIDSlice = append(hotDogIDSlice, -1) //This is a code fix for null slices getting passed
		sendData.ID_HotDogs = hotDogIDSlice
		debugHotDog := Hotdog{
			HotDogType: "DEBUGTYPE",
			Condiment:  "DEBUGCONDIMENT",
			Calories:   0,
			Name:       "DEBUGNAME",
			UserID:     0,
		}
		hotDogSlice = append(hotDogSlice, debugHotDog)
		sendData.TheHotDogs = hotDogSlice
	} else if len(sendData.ID_Hamburgers) == 0 {
		hamburgIDSlice = append(hamburgIDSlice, -1) //This is a code fix for null slices getting passed
		sendData.ID_Hamburgers = hamburgIDSlice
		debugHamburger := Hamburger{
			BurgerType: "DEBUGTYPE",
			Condiment:  "DEBUGCONDIMENT",
			Calories:   0,
			Name:       "DEBUGNAME",
			UserID:     0,
		}
		hamburgerSlice = append(hamburgerSlice, debugHamburger)
		sendData.TheHamburgers = hamburgerSlice
	}

	dataJSON, err := json.Marshal(sendData)
	if err != nil {
		fmt.Println("There's an error marshalling.")
	}

	fmt.Fprintf(w, string(dataJSON))
}

//UPDATE FOOD
func updateFood(w http.ResponseWriter, req *http.Request) {
	type foodUpdate struct {
		FoodType     string    `json:"FoodType"`
		FoodID       int       `json:"FoodID"`
		TheHamburger Hamburger `json:"TheHamburger"`
		TheHotDog    Hotdog    `json:"TheHotDog"`
	}
	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	//Marshal it into our type
	var thefoodUpdate foodUpdate
	json.Unmarshal(bs, &thefoodUpdate)

	sqlStatement := ""

	//Determine if this is a hotdog or hamburger update
	if thefoodUpdate.FoodType == "hotdog" {
		var updatedHotdog Hotdog = thefoodUpdate.TheHotDog
		/* CHECK TO SEE IF THE FIELDS ARE OKAY */
		canPost := true
		if containsLanguage(updatedHotdog.HotDogType) {
			canPost = false
		} else if containsLanguage(updatedHotdog.Condiment) {
			canPost = false
		} else if containsLanguage(updatedHotdog.Name) {
			canPost = false
		} else {
			canPost = true
		}
		if canPost == true {
			sqlStatement = "UPDATE hot_dogs SET TYPE=?, CONDIMENT=?, CALORIES=?," +
				"NAME=?, PHOTO_ID=?, PHOTO_SRC=?, DATE_UPDATED=? WHERE FOOD_ID=?"

			stmt, err := db.Prepare(sqlStatement)
			check(err)
			theTimeNow := time.Now()
			r, err := stmt.Exec(updatedHotdog.HotDogType, updatedHotdog.Condiment,
				updatedHotdog.Calories, updatedHotdog.Name, updatedHotdog.PhotoID,
				updatedHotdog.PhotoSrc, updatedHotdog.UserID,
				theTimeNow.Format("2006-01-02 15:04:05"), thefoodUpdate.FoodID)
			check(err)

			n, err := r.RowsAffected()
			check(err)

			fmt.Printf("%v\n", n)

			if n < 1 {
				fmt.Printf("Only %v rows effected, foodUpdate unsuccessful. No foodID found for: %v\n", n,
					thefoodUpdate.FoodID)
				fmt.Fprintln(w, 3)
			} else {
				fmt.Fprintln(w, 1)
			}
		} else {
			fmt.Printf("Language detected in food. Update unsuccessful.\n")
			fmt.Fprintln(w, 4)
		}
	} else if thefoodUpdate.FoodType == "hamburger" {
		var updatedHamburger Hamburger = thefoodUpdate.TheHamburger
		/* CHECK TO SEE IF THE FIELDS ARE OKAY */
		canPost := true
		if containsLanguage(updatedHamburger.BurgerType) {
			canPost = false
		} else if containsLanguage(updatedHamburger.Condiment) {
			canPost = false
		} else if containsLanguage(updatedHamburger.Name) {
			canPost = false
		} else {
			canPost = true
		}
		if canPost == true {
			sqlStatement = "UPDATE hamburgers SET TYPE=?, CONDIMENT=?, CALORIES=?," +
				"NAME=?, PHOTO_ID=?, PHOTO_SRC=?, DATE_UPDATED=? WHERE FOOD_ID=?"

			stmt, err := db.Prepare(sqlStatement)
			check(err)
			theTimeNow := time.Now()
			r, err := stmt.Exec(updatedHamburger.BurgerType, updatedHamburger.Condiment,
				updatedHamburger.Calories, updatedHamburger.Name, updatedHamburger.PhotoID,
				updatedHamburger.PhotoSrc,
				theTimeNow.Format("2006-01-02 15:04:05"), updatedHamburger.FoodID)
			check(err)

			n, err := r.RowsAffected()
			check(err)

			if n < 1 {
				fmt.Printf("Only %v rows effected, foodUpdate unsuccessful. No foodID found for: %v\n", n,
					thefoodUpdate.FoodID)
				fmt.Fprintln(w, 3)
			} else {
				fmt.Printf("%v\n", n)
				fmt.Fprintln(w, 2)
			}
		} else {
			fmt.Printf("Language detected in food. Update unsuccessful.\n")
			fmt.Fprintln(w, 4)
		}
	} else {
		fmt.Printf("No good value sent through JSON to update food: %v\n", thefoodUpdate.FoodID)
		fmt.Fprintln(w, 3)
	}
}

//INSERT USER(s)
func insertUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting User record.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("We go an error reading the JSON: %v\n", err.Error())
		fmt.Println(err)
	}

	//Marshal it into our type
	var postedUser User
	json.Unmarshal(bs, &postedUser)

	//Add User to the SQL Database
	stmt, err := db.Prepare("INSERT INTO users(USERNAME, PASSWORD, FIRSTNAME, LASTNAME, ROLE, USER_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

	r, err := stmt.Exec(postedUser.UserName, postedUser.Password, postedUser.First,
		postedUser.Last, postedUser.Role, postedUser.UserID, postedUser.DateCreated, postedUser.DateUpdated)
	check(err)

	n, err2 := r.RowsAffected()
	check(err2)
	stmt.Close()

	fmt.Printf("Inserted Record: %v\n", n)
}

//GET USER(S)
func getUsers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Getting User record(s).")
	//Define the JSON submitted
	type UserCollection struct {
		UserIDs []int `json:"UserIDs"`
	}
	//Define data returned
	type userReturned struct {
		DBID     int    `json:"DBID"`
		UserName string `json:"UserName"`
		Password string `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
		First    string `json:"First"`
		Last     string `json:"Last"`
		Role     string `json:"Role"`
		UserID   int    `json:"UserID"`
	}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	//Marshal it into our type
	var postedUserIDs UserCollection
	json.Unmarshal(bs, &postedUserIDs)
	theIDS := postedUserIDs.UserIDs //Get all the IDS into an easier to read variable

	//Assign variables to construct our Users into
	var theReturnedUsers []userReturned
	var aUser userReturned
	//Query to get all the User IDs in database from slice of UserIDs
	for x := 0; x < len(theIDS); x++ {
		anID := theIDS[x] //Assign 1 ID to a variable
		//Run Query on that ID
		stmt := "SELECT * FROM users WHERE USER_ID = ?"
		rows, err := db.Query(stmt, anID)
		check(err)

		//Get User returned
		for rows.Next() {
			err = rows.Scan(&aUser.DBID, &aUser.UserName, &aUser.Password,
				&aUser.First, &aUser.Last, &aUser.Role, &aUser.UserID)
			check(err)
			theReturnedUsers = append(theReturnedUsers, aUser)
		}
		rows.Close()
	}

	//Check to see if theReturnedUsers got nothing; else, send the full JSON
	if len(theReturnedUsers) <= 0 {
		fmt.Printf("No Users returned\n")
		theJSONMessage, err := json.Marshal(theReturnedUsers)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(theJSONMessage))
	} else {
		theJSONMessage, err := json.Marshal(theReturnedUsers)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(theJSONMessage))
	}
}

//UPDATE USER(S)
func updateUsers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Updating User record(s).")
	//Define the JSON submitted
	type UserCollection struct {
		TheUsers []User `json:"UserIDs"`
	}

	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	//Marshal it into our type
	var submittedUsers UserCollection
	json.Unmarshal(bs, &submittedUsers)

	returnedString := "" //TheJSon string to return
	//Update the Users from the gotten UserID
	for x := 0; x < len(submittedUsers.TheUsers); x++ {
		anID := submittedUsers.TheUsers[x].UserID //Assign 1 ID to a variable
		//Run Query on that ID
		stmt := "UPDATE users SET USERNAME=?, PASSWORD=?, FIRSTNAME=?, LASTNAME=?, ROLE=?, DATE_UPDATED=?" +
			" WHERE USER_ID=?"
		theStmt, err := db.Prepare(stmt)
		check(err)
		theTimeNow := time.Now()
		r, err := theStmt.Exec(submittedUsers.TheUsers[x].UserName, submittedUsers.TheUsers[x].Password,
			submittedUsers.TheUsers[x].First, submittedUsers.TheUsers[x].Last,
			submittedUsers.TheUsers[x].Role, theTimeNow.Format("2006-01-02 15:04:05"), anID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		theStmt.Close()
		returnedString = returnedString + " " + "updated at ID " + string(anID) + " " + string(n)
	}

	//Return the JSON Message
	theJSONMessage, err := json.Marshal(returnedString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//DELETE USER(S)
func deleteUsers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("DELETING User record(s).")
	//Define the JSON submitted
	type UserCollection struct {
		UserIDs []int `json:"UserIDs"`
	}

	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	//Marshal it into our type
	var postedUserIDs UserCollection
	json.Unmarshal(bs, &postedUserIDs)
	theIDS := postedUserIDs.UserIDs //Get all the IDS into an easier to read variable
	//Query to get all the User IDs in database from slice of UserIDs to delete
	for x := 0; x < len(theIDS); x++ {
		anID := theIDS[x] //Assign 1 ID to a variable
		//Run Query on that ID
		stmt := "DELETE FROM users WHERE USER_ID = ?"
		rows, err := db.Query(stmt, anID)
		check(err)
		defer rows.Close()
	}
	//Return the JSON response
	returnString := "We deleted the following: "
	for y := 0; y < len(theIDS); y++ {
		theStringID := strconv.Itoa(theIDS[y])
		returnString = returnString + " " + theStringID + " "
	}
	theJSONMessage, err := json.Marshal(returnString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/*************** AMAZON PHOTO QUERIES *********************/
//Insert Photos into SQL
func insertUserPhotos(userid int, foodid int, photoid int, photoName string, fileType string, size int64,
	photoHash string, link string, foodType string, dateCreated string, dateUpdated string) bool {
	successfulInsert := true
	fmt.Printf("DEBUG: Inserting photos into SQL.\n")
	theTimeNow := time.Now()
	//Which Type of food?
	if strings.Contains(foodType, "HOTDOG") {
		fmt.Printf("Inserting Hotdog Photo\n")
		theStatement := "INSERT INTO user_photos" +
			"(USER_ID, FOOD_ID, PHOTO_ID, PHOTO_NAME, FILE_TYPE, SIZE, PHOTO_HASH, LINK, FOOD_TYPE, DATE_CREATED, DATE_UPDATED) " +
			"VALUES(?,?,?,?,?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(theStatement)

		r, err := stmt.Exec(userid, foodid, photoid, photoName, fileType,
			size, photoHash, link, foodType, theTimeNow.Format("2006-01-02 15:04:05"),
			theTimeNow.Format("2006-01-02 15:04:05"))
		check(err)

		n, err := r.RowsAffected()
		check(err)
		fmt.Printf("%v rows effected.\n", n)
		stmt.Close() //Close the SQL
		/********* UPDATE SQL WITH PHOTO INFORMATION ************/
		type foodUpdate struct {
			FoodType     string    `json:"FoodType"`
			FoodID       int       `json:"FoodID"`
			TheHamburger Hamburger `json:"TheHamburger"`
			TheHotDog    Hotdog    `json:"TheHotDog"`
		}
		//Get Hotdog from FoodID
		var theHotDog Hotdog
		theStmt := "SELECT TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, " +
			"PHOTO_ID, PHOTO_SRC, DATE_CREATED, DATE_UPDATED FROM hot_dogs WHERE FOOD_ID = ?"
		rows, err := db.Query(theStmt, foodid)
		check(err)
		for rows.Next() {
			err = rows.Scan(&theHotDog.HotDogType, &theHotDog.Condiment, &theHotDog.Calories,
				&theHotDog.Name, &theHotDog.UserID, &theHotDog.FoodID,
				&theHotDog.PhotoID, &theHotDog.PhotoSrc, &theHotDog.DateCreated, &theHotDog.DateUpdated)
			check(err)
			theHotDog.PhotoID = photoid
			newLink := filepath.Join("static", "images", link)
			theLink := urlFixer(newLink)
			theHotDog.PhotoSrc = theLink
			theHotDog.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
			fmt.Printf("DEBUG: Here is our hotdog: %v\n", theHotDog)
		}
		rows.Close()
		hotDogUpdate := foodUpdate{
			FoodType:     "hotdog",
			FoodID:       foodid,
			TheHamburger: Hamburger{},
			TheHotDog:    theHotDog,
		}
		jsonValue, _ := json.Marshal(hotDogUpdate)
		response, err := http.Post("http://localhost:80/updateFood", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			successfulInsert = false
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	} else if strings.Contains(foodType, "HAMBURGER") {
		fmt.Printf("Inserting Hamburger Photo into SQL\n")
		theStatement := "INSERT INTO user_photos " +
			"(USER_ID, FOOD_ID, PHOTO_ID, PHOTO_NAME, FILE_TYPE, SIZE, PHOTO_HASH, LINK, FOOD_TYPE, DATE_CREATED, DATE_UPDATED) " +
			"VALUES(?,?,?,?,?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(theStatement)

		r, err := stmt.Exec(userid, foodid, photoid, photoName, fileType,
			size, photoHash, link, foodType, theTimeNow.Format("2006-01-02 15:04:05"),
			theTimeNow.Format("2006-01-02 15:04:05"))
		check(err)

		n, err := r.RowsAffected()
		check(err)
		fmt.Printf("%v rows effected.\n", n)
		stmt.Close() //Close the SQL
		/********* UPDATE SQL WITH PHOTO INFORMATION ************/
		type foodUpdate struct {
			FoodType     string    `json:"FoodType"`
			FoodID       int       `json:"FoodID"`
			TheHamburger Hamburger `json:"TheHamburger"`
			TheHotDog    Hotdog    `json:"TheHotDog"`
		}
		//Get Hotdog from FoodID
		var theHamb Hamburger
		theStmt := "SELECT TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, " +
			"PHOTO_ID, PHOTO_SRC, DATE_CREATED, DATE_UPDATED FROM hamburgers WHERE FOOD_ID = ?"
		rows, err := db.Query(theStmt, foodid)
		check(err)
		for rows.Next() {
			err = rows.Scan(&theHamb.BurgerType, &theHamb.Condiment, &theHamb.Calories,
				&theHamb.Name, &theHamb.UserID, &theHamb.FoodID,
				&theHamb.PhotoID, &theHamb.PhotoSrc, &theHamb.DateCreated, &theHamb.DateUpdated)
			check(err)
			theHamb.PhotoID = photoid
			newLink := filepath.Join("static", "images", link)
			theLink := urlFixer(newLink)
			theHamb.PhotoSrc = theLink
			theHamb.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
			fmt.Printf("DEBUG: Here is our hotdog: %v\n", theHamb)
		}
		rows.Close()
		hamUpdate := foodUpdate{
			FoodType:     "hamburger",
			FoodID:       foodid,
			TheHamburger: theHamb,
			TheHotDog:    Hotdog{},
		}
		jsonValue, _ := json.Marshal(hamUpdate)
		response, err := http.Post("http://localhost:80/updateFood", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			successfulInsert = false
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	} else {
		fmt.Printf("Wrong food type returned for food photo data insertion: %v\n", foodType)
		successfulInsert = false
	}

	return successfulInsert
}

func getUserPhotos(userID int, theFood string) ([]string, []string) {
	var thePhotoURLS []string
	var theFileNames []string
	if strings.Contains(theFood, "HOTDOG") {
		stmt := "SELECT DISTINCT LINK, PHOTO_NAME FROM user_photos WHERE FOOD_TYPE = ? AND USER_ID = ?"
		rows, err := db.Query(stmt, "HOTDOG", userID)
		check(err)

		var aPhotoURL string
		var aFileName string
		for rows.Next() {
			err = rows.Scan(&aPhotoURL, &aFileName)
			check(err)
			thePhotoURLS = append(thePhotoURLS, aPhotoURL)
			theFileNames = append(theFileNames, aFileName)
		}
		rows.Close()
		fmt.Printf("There were %v Hotdog URLS returned for User, %v:\n%v\n%v\n", len(thePhotoURLS), userID, thePhotoURLS,
			theFileNames)
		return thePhotoURLS, theFileNames
	} else {
		stmt := "SELECT DISTINCT LINK, PHOTO_NAME FROM user_photos WHERE FOOD_TYPE = ? AND USER_ID = ?"
		rows, err := db.Query(stmt, "HAMBURGER", userID)
		check(err)

		var aPhotoURL string
		var aFileName string

		for rows.Next() {
			err = rows.Scan(&aPhotoURL, &aFileName)
			check(err)
			thePhotoURLS = append(thePhotoURLS, aPhotoURL)
			theFileNames = append(theFileNames, aFileName)
		}
		rows.Close()
		fmt.Printf("There were %v Hamburger URLS returned for User, %v:\n%v\n%v\n", len(thePhotoURLS), userID, thePhotoURLS,
			theFileNames)
		return thePhotoURLS, theFileNames
	}
}
