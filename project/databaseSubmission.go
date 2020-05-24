package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

	stmt, err := db.Prepare("INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES(?,?,?,?,?)")
	defer stmt.Close()

	r, err := stmt.Exec(postedHotDog.HotDogType, postedHotDog.Condiment, postedHotDog.Calories, postedHotDog.Name, postedHotDog.UserID)
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("DEBUG: %v rows effected.\n", n)

	if err != nil {
		fmt.Fprint(w, failureMessage)
	} else {
		hDogMarshaled, err := json.Marshal(postedHotDog)
		if err != nil {
			fmt.Printf("Error with %v\n", hDogMarshaled)
		}
		hDogSuccessMSG := successMessage + string(hDogMarshaled)
		fmt.Fprint(w, hDogSuccessMSG)
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
	defer rows.Close()
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
		sqlStatement = "DELETE FROM hot_dogs WHERE ID=?"
		delDog, err := db.Prepare(sqlStatement)
		check(err)

		r, err := delDog.Exec(theFoodDeletion.FoodID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("%v\n", n)

		fmt.Fprintln(w, 1)
	} else if theFoodDeletion.FoodType == "hamburger" {
		sqlStatement = "DELETE FROM hamburgers WHERE ID=?"
		delDog, err := db.Prepare(sqlStatement)
		check(err)

		r, err := delDog.Exec(theFoodDeletion.FoodID)
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

	fmt.Printf("DEBUG: HERE IS OUR postedHamburger: \n%v\n", postedHamburger)

	stmt, err := db.Prepare("INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES(?,?,?,?,?)")
	defer stmt.Close()

	r, err := stmt.Exec(postedHamburger.BurgerType, postedHamburger.Condiment,
		postedHamburger.Calories, postedHamburger.Name, postedHamburger.UserID)
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Printf("DEBUG: %v rows effected.\n", n)

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

	//Declare variables for hotdog
	var ham_id int
	var ham_type string
	var ham_condiment string
	var ham_calories int
	var ham_name string
	var ham_userID int

	//Declare all our our hotdog/hamburger collections returned
	var hotDogSlice []Hotdog
	var hamburgerSlice []Hamburger
	var hotDogIDSlice []int
	var hamburgIDSlice []int

	//Assemble data to send back
	type data struct {
		SuccessMessage string      `json:"SuccessMessage"`
		TheHotDogs     []Hotdog    `json:"TheHotDogs"`
		TheHamburgers  []Hamburger `json:"TheHamburgers:`
		ID_HotDogs     []int       `json:"ID_HotDogs"`
		ID_Hamburgers  []int       `json:"ID_Hamburgers"`
	}

	//Counter for food returned
	dogCounter := 0
	hamCounter := 0

	//If no User ID is submitted, then just food for ALL Users
	if theUser.UserID == 0 {
		//Get HotDogs
		hrows, err1 := db.Query(`SELECT * FROM hot_dogs ORDER BY ID;`)
		check(err1)
		defer hrows.Close()

		for hrows.Next() {
			err = hrows.Scan(&h_id, &h_dogType, &h_condiment, &h_calories, &h_hotdogName, &h_userID)
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
		hamrows, err2 := db.Query(`SELECT * FROM hamburgers ORDER BY ID`)
		check(err2)
		defer hamrows.Close()

		for hamrows.Next() {
			err = hamrows.Scan(&ham_id, &ham_type, &ham_condiment, &ham_calories, &ham_name, &ham_userID)
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
	} else {
		//Get HotDogs
		hrows, err1 := db.Query(`SELECT * FROM hot_dogs WHERE USER_ID=? ORDER BY ID;`, theUser.UserID)
		check(err1)
		defer hrows.Close()

		for hrows.Next() {
			err = hrows.Scan(&h_id, &h_dogType, &h_condiment, &h_calories, &h_hotdogName, &h_userID)
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
		defer hamrows.Close()

		for hamrows.Next() {
			err = hamrows.Scan(&ham_id, &ham_type, &ham_condiment, &ham_calories, &ham_name, &ham_userID)
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
	}

	//Check to see if we have any data to submit
	sendData := data{
		SuccessMessage: "Success",
		TheHotDogs:     hotDogSlice,
		TheHamburgers:  hamburgerSlice,
		ID_HotDogs:     hotDogIDSlice,
		ID_Hamburgers:  hamburgIDSlice,
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
		fmt.Printf("DEBUG: Updating hotdog at id: %v\n", thefoodUpdate.FoodID)
		var updatedHotdog Hotdog = thefoodUpdate.TheHotDog
		sqlStatement = "UPDATE hot_dogs SET TYPE=?, CONDIMENT=?, CALORIES=?," +
			"NAME=?, USER_ID=? WHERE ID=?"

		stmt, err := db.Prepare(sqlStatement)
		check(err)

		r, err := stmt.Exec(updatedHotdog.HotDogType, updatedHotdog.Condiment,
			updatedHotdog.Calories, updatedHotdog.Name, updatedHotdog.UserID,
			thefoodUpdate.FoodID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("%v\n", n)

		fmt.Fprintln(w, 1)

	} else if thefoodUpdate.FoodType == "hamburger" {
		var updatedHamburger Hamburger = thefoodUpdate.TheHamburger
		sqlStatement = "UPDATE hamburgers SET TYPE=?, CONDIMENT=?, CALORIES=?," +
			"NAME=?, USER_ID=? WHERE ID=?"

		stmt, err := db.Prepare(sqlStatement)
		check(err)

		r, err := stmt.Exec(updatedHamburger.BurgerType, updatedHamburger.Condiment,
			updatedHamburger.Calories, updatedHamburger.Name, updatedHamburger.UserID,
			thefoodUpdate.FoodID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Printf("%v\n", n)

		fmt.Fprintln(w, 2)
	} else {
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
		fmt.Println(err)
	}

	//Marshal it into our type
	var postedUser User
	json.Unmarshal(bs, &postedUser)

	//Add User to the SQL Database
	stmt, err := db.Prepare("INSERT INTO users(USERNAME, PASSWORD, FIRSTNAME, LASTNAME, ROLE, USER_ID) VALUES(?,?,?,?,?,?)")
	defer stmt.Close()

	r, err := stmt.Exec(postedUser.UserName, postedUser.Password, postedUser.First,
		postedUser.Last, postedUser.Role, postedUser.UserID)
	check(err)

	n, err := r.RowsAffected()
	check(err)

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
		defer rows.Close()

		//Get User returned
		for rows.Next() {
			err = rows.Scan(&aUser.DBID, &aUser.UserName, &aUser.Password,
				&aUser.First, &aUser.Last, &aUser.Role, &aUser.UserID)
			check(err)
			theReturnedUsers = append(theReturnedUsers, aUser)
		}
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
		stmt := "UPDATE users SET USERNAME=?, PASSWORD=?, FIRSTNAME=?, LASTNAME=?, ROLE=?" +
			" WHERE USER_ID=?"
		theStmt, err := db.Prepare(stmt)
		check(err)
		defer theStmt.Close()

		r, err := theStmt.Exec(submittedUsers.TheUsers[x].UserName, submittedUsers.TheUsers[x].Password,
			submittedUsers.TheUsers[x].First, submittedUsers.TheUsers[x].Last,
			submittedUsers.TheUsers[x].Role, anID)
		check(err)

		n, err := r.RowsAffected()
		check(err)

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
