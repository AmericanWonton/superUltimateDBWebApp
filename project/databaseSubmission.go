package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const successMessage string = "Successful Insert"
const failureMessage string = "Unsuccessful Insert"

//POST mainpage
func insertHotDog(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting hotdog record.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Here is our byte slice: %v\n", bs)
	//Marshal it into our type
	var postedHotDog Hotdog
	json.Unmarshal(bs, &postedHotDog)
	//Debug
	fmt.Printf("DEBUG: Here is our hotdog: \n%v\n", postedHotDog)

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
		fmt.Fprint(w, successMessage)
	}
}

//GET all Food, Mainpage
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
	var h_id int64
	var h_dogType string
	var h_condiment string
	var h_calories int
	var h_hotdogName string
	var h_userID int

	//Declare variables for hotdog
	var ham_id int64
	var ham_type string
	var ham_condiment string
	var ham_calories int
	var ham_name string
	var ham_userID int

	//Declare all our our hotdog/hamburger collections returned
	var hotDogSlice []Hotdog
	var hamburgerSlice []Hamburger

	//Counter for food returned
	dogCounter := 0
	hamCounter := 0

	//Get HotDogs
	hrows, err1 := db.Query(`SELECT * FROM hot_dogs WHERE USER_ID=?;`, theUser.UserID)
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
		dogCounter = dogCounter + 1
	}

	//Get Hamburgers
	hamrows, err2 := db.Query(`SELECT * FROM hamburgers WHERE USER_ID=?`, theUser.UserID)
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
		hamCounter = hamCounter + 1
	}

	//Assemble data to send back
	type data struct {
		SuccessMessage string      `json:"SuccessMessage"`
		TheHotDogs     []Hotdog    `json:"TheHotDogs"`
		TheHamburgers  []Hamburger `json:"TheHamburgers:`
	}
	//Check to see if we have any data to submit
	sendData := data{
		SuccessMessage: "Success",
		TheHotDogs:     hotDogSlice,
		TheHamburgers:  hamburgerSlice,
	}

	if len(sendData.TheHotDogs) <= 0 && len(sendData.TheHamburgers) <= 0 {
		sendData.SuccessMessage = "Failure"
	}

	dataJSON, err := json.Marshal(sendData)
	if err != nil {
		fmt.Println("There's an error marshalling.")
	}

	fmt.Fprintf(w, string(dataJSON))
}
