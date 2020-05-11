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
