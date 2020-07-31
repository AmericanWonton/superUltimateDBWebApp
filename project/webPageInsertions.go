package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func hotDogInsertWebPage(w http.ResponseWriter, req *http.Request) {
	//Declare the Struct
	fmt.Println("Inserting hotdog record in Mongo/SQL.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Here is our byte slice hotdog as a string for JSON: \n\n%v\n", string(bs))
	//Marshal it into our type
	var postedHotDog Hotdog
	json.Unmarshal(bs, &postedHotDog)
	//Protections for the hotdog name
	if strings.Compare(postedHotDog.HotDogType, "DEBUGTYPE") == 0 {
		postedHotDog.HotDogType = "NONE"
	}
	//First give this hotdog a random ID
	fmt.Printf("DEBUG: This is what our hotdog foodID is now: %v\n", postedHotDog.FoodID)
	randomFoodID := randomIDCreation()
	postedHotDog.FoodID = randomFoodID
	fmt.Printf("DEBUG: Here is our randomID now: %v\n", postedHotDog.FoodID)
	//Give the correct time to this hotdog
	theTimeNow := time.Now()
	postedHotDog.DateCreated = theTimeNow.Format("2006-01-02 15:04:05")
	postedHotDog.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")

	//Define returned JSON
	type returnData struct {
		SuccessMsg     string      `json:"SuccessMsg"`
		ReturnedHotDog MongoHotDog `json:"ReturnedHotDog"`
		SuccessBool    bool        `json:"SuccessBool"`
	}

	//Insert into SQL
	stmt, err := db.Prepare("INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

	r, err := stmt.Exec(postedHotDog.HotDogType, postedHotDog.Condiment, postedHotDog.Calories, postedHotDog.Name, postedHotDog.UserID,
		postedHotDog.FoodID, postedHotDog.DateCreated, postedHotDog.DateUpdated)
	check(err)

	n, err := r.RowsAffected()
	check(err)
	stmt.Close() //Close the SQL

	//Insert into Mongo
	//Declare Mongo Dog
	mongoHotDogInsert := MongoHotDog{
		HotDogType:  postedHotDog.HotDogType,
		Condiments:  turnFoodArray(postedHotDog.Condiment),
		Calories:    postedHotDog.Calories,
		Name:        postedHotDog.Name,
		FoodID:      postedHotDog.FoodID,
		UserID:      postedHotDog.UserID,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}
	//Collect Data for Mongo
	user_collection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
	collectedUsers := []interface{}{mongoHotDogInsert}
	//Insert Our Data
	insertManyResult, err2 := user_collection.InsertMany(theContext, collectedUsers)
	if err2 != nil {
		theReturnData := returnData{
			SuccessMsg:     failureMessage,
			ReturnedHotDog: mongoHotDogInsert,
			SuccessBool:    false,
		}
		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling this hotdog.")
			logWriter("There's an error marshalling.")
		}
		fmt.Fprintf(w, string(dataJSON))
	} else {
		theReturnData := returnData{
			SuccessMsg:     successMessage,
			ReturnedHotDog: mongoHotDogInsert,
			SuccessBool:    true,
		}
		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling.")
			logWriter("There's an error marshalling.")
		}
		fmt.Fprintf(w, string(dataJSON))
		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
	}

	fmt.Printf("DEBUG: %v rows effected.\n", n)
}

func hamburgerInsertWebPage(w http.ResponseWriter, req *http.Request) {
	//Declare the Struct
	fmt.Println("Inserting Burger record in Mongo/SQL.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Here is our byte slice Burger as a string for JSON: \n\n%v\n", string(bs))
	//Marshal it into our type
	var postedHamburger Hamburger
	json.Unmarshal(bs, &postedHamburger)
	//Protections for the Burger name
	if strings.Compare(postedHamburger.BurgerType, "DEBUGTYPE") == 0 {
		postedHamburger.BurgerType = "NONE"
	}
	//First give this Burger a random ID
	fmt.Printf("DEBUG: This is what our Burger foodID is now: %v\n", postedHamburger.FoodID)
	randomFoodID := randomIDCreation()
	postedHamburger.FoodID = randomFoodID
	fmt.Printf("DEBUG: Here is our randomID now: %v\n", postedHamburger.FoodID)
	//Give the correct time to this hotdog
	theTimeNow := time.Now()
	postedHamburger.DateCreated = theTimeNow.Format("2006-01-02 15:04:05")
	postedHamburger.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")

	//Define returned JSON
	type returnData struct {
		SuccessMsg        string         `json:"SuccessMsg"`
		ReturnedHamburger MongoHamburger `json:"ReturnedHamburger"`
		SuccessBool       bool           `json:"SuccessBool"`
	}

	//Insert into SQL
	stmt, err := db.Prepare("INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, FOOD_ID, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?)")

	r, err := stmt.Exec(postedHamburger.BurgerType, postedHamburger.Condiment, postedHamburger.Calories, postedHamburger.Name, postedHamburger.UserID,
		postedHamburger.FoodID, postedHamburger.DateCreated, postedHamburger.DateUpdated)
	check(err)

	n, err := r.RowsAffected()
	check(err)
	stmt.Close() //Close the SQL

	//Insert into Mongo
	//Declare Mongo Dog
	mongoHamburgerInsert := MongoHamburger{
		BurgerType:  postedHamburger.BurgerType,
		Condiments:  turnFoodArray(postedHamburger.Condiment),
		Calories:    postedHamburger.Calories,
		Name:        postedHamburger.Name,
		FoodID:      postedHamburger.FoodID,
		UserID:      postedHamburger.UserID,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}
	//Collect Data for Mongo
	user_collection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
	collectedUsers := []interface{}{mongoHamburgerInsert}
	//Insert Our Data
	insertManyResult, err2 := user_collection.InsertMany(theContext, collectedUsers)
	if err2 != nil {
		theReturnData := returnData{
			SuccessMsg:        failureMessage,
			ReturnedHamburger: mongoHamburgerInsert,
			SuccessBool:       false,
		}
		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling this hamburger.")
			logWriter("There's an error marshalling.")
		}
		fmt.Fprintf(w, string(dataJSON))
	} else {
		theReturnData := returnData{
			SuccessMsg:        successMessage,
			ReturnedHamburger: mongoHamburgerInsert,
			SuccessBool:       true,
		}
		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling.")
			logWriter("There's an error marshalling.")
		}
		fmt.Fprintf(w, string(dataJSON))
		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
	}

	fmt.Printf("DEBUG: %v rows effected.\n", n)
}
