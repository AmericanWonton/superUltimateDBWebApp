package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
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
	//Define returned JSON
	type returnData struct {
		SuccessMsg     string      `json:"SuccessMsg"`
		ReturnedHotDog MongoHotDog `json:"ReturnedHotDog"`
		SuccessBool    bool        `json:"SuccessBool"`
	}
	/* CHECK TO SEE IF THIS HOTDOG INSERTION IS CLEAN */
	canPost := true
	if containsLanguage(postedHotDog.HotDogType) {
		canPost = false
	} else if containsLanguage(postedHotDog.Condiment) {
		canPost = false
	} else if containsLanguage(postedHotDog.Name) {
		canPost = false
	} else {
		canPost = true
	}
	if canPost == true {
		//First give this hotdog a random ID
		randomFoodID := randomIDCreation()
		postedHotDog.FoodID = randomFoodID
		fmt.Printf("DEBUG: Here is our randomID now: %v\n", postedHotDog.FoodID)
		//Give the correct time to this hotdog
		theTimeNow := time.Now()
		postedHotDog.DateCreated = theTimeNow.Format("2006-01-02 15:04:05")
		postedHotDog.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")

		//Insert into SQL
		theStatement := "INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID," +
			" FOOD_ID, PHOTO_ID, PHOTO_SRC, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(theStatement)

		r, err := stmt.Exec(postedHotDog.HotDogType, postedHotDog.Condiment, postedHotDog.Calories, postedHotDog.Name, postedHotDog.UserID,
			postedHotDog.FoodID, 0, "",
			postedHotDog.DateCreated, postedHotDog.DateUpdated)
		check(err)

		n, err := r.RowsAffected()
		check(err)
		stmt.Close() //Close the SQL

		fmt.Printf("DEBUG: %v rows effected for SQL.\n", n)

		//Insert into Mongo
		//Declare Mongo Dog
		mongoHotDogInsert := MongoHotDog{
			HotDogType:  postedHotDog.HotDogType,
			Condiments:  turnFoodArray(postedHotDog.Condiment),
			Calories:    postedHotDog.Calories,
			Name:        postedHotDog.Name,
			FoodID:      postedHotDog.FoodID,
			UserID:      postedHotDog.UserID,
			PhotoID:     0,
			PhotoSrc:    "",
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Collect Data for Mongo
		hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
		collectedUsers := []interface{}{mongoHotDogInsert}
		//Insert Our Data
		insertManyResult, err2 := hotdogCollection.InsertMany(theContext, collectedUsers)
		if err2 != nil {
			//Marshal the bad news
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
			//Insert Hotdog data into User array
			user_collection := mongoClient.Database("superdbtest1").Collection("users")
			filterUserID := bson.M{"userid": postedHotDog.UserID}
			var foundUser AUser
			foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
			theErr := user_collection.FindOne(theContext, filterUserID).Decode(&foundUser)
			if theErr != nil {
				if strings.Contains(theErr.Error(), "no documents in result") {
					fmt.Printf("It's all good, this document didn't find this UserID: %v\n", postedHotDog.UserID)
				} else {
					fmt.Printf("DEBUG: We have an error finding User for this hotdogUserID: %v\n%v\n", postedHotDog.UserID,
						theErr)
				}
				//Marshal the bad news
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
				//Start updating User
				fmt.Printf("Found the foundUser: %v\n", foundUser)
				foundUser.Hotdogs.Hotdogs = append(foundUser.Hotdogs.Hotdogs, mongoHotDogInsert)
				successfulUserInsert := updateUser(foundUser) //Update this User with the new Hotdog Array
				if successfulUserInsert == true {
					fmt.Printf("This User's hotdogs was updated successfully: %v\n", foundUser.UserID)
					//Set data for photo insertion
					awsuserID = postedHotDog.UserID
					awsfoodID = postedHotDog.FoodID
					awsphotoName = postedHotDog.Name
					awsfoodType = "HOTDOG"
					awsdateUpdated = postedHotDog.DateCreated
					awsdateCreated = postedHotDog.DateUpdated
					//Marshal return data
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
				} else {
					fmt.Printf("This User's hotdogs were NOT updated successfully: %v\n", foundUser.UserID)
					//Marshal the bad news
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
				}
			}
		}
	} else {
		//Marshal the bad news
		theReturnData := returnData{
			SuccessMsg:     "Your food contained foul language, please re-type",
			ReturnedHotDog: MongoHotDog{},
			SuccessBool:    false,
		}
		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling this hotdog.")
			logWriter("There's an error marshalling.")
		}
		fmt.Fprintf(w, string(dataJSON))
	}
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
	//Marshal it into our type
	var postedHamburger Hamburger
	json.Unmarshal(bs, &postedHamburger)
	//Protections for the Burger name
	if strings.Compare(postedHamburger.BurgerType, "DEBUGTYPE") == 0 {
		postedHamburger.BurgerType = "NONE"
	}
	//Define returned JSON
	type returnData struct {
		SuccessMsg        string         `json:"SuccessMsg"`
		ReturnedHamburger MongoHamburger `json:"ReturnedHamburger"`
		SuccessBool       bool           `json:"SuccessBool"`
	}
	/* CHECK TO SEE IF THIS Hamburger INSERTION IS CLEAN */
	canPost := true
	if containsLanguage(postedHamburger.BurgerType) {
		canPost = false
	} else if containsLanguage(postedHamburger.Condiment) {
		canPost = false
	} else if containsLanguage(postedHamburger.Name) {
		canPost = false
	} else {
		canPost = true
	}
	if canPost == true {
		//First give this Burger a random ID
		randomFoodID := randomIDCreation()
		postedHamburger.FoodID = randomFoodID
		//Give the correct time to this hotdog
		theTimeNow := time.Now()
		postedHamburger.DateCreated = theTimeNow.Format("2006-01-02 15:04:05")
		postedHamburger.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")

		//Insert into SQL
		theStatement := "INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, " +
			"FOOD_ID, PHOTO_ID, PHOTO_SRC,  DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(theStatement)

		r, err := stmt.Exec(postedHamburger.BurgerType, postedHamburger.Condiment, postedHamburger.Calories, postedHamburger.Name, postedHamburger.UserID,
			postedHamburger.FoodID, 0, "",
			postedHamburger.DateCreated, postedHamburger.DateUpdated)
		check(err)

		n, err := r.RowsAffected()
		check(err)
		stmt.Close() //Close the SQL

		fmt.Printf("DEBUG: %v rows effected for SQL.\n", n)

		//Insert into Mongo
		//Declare Mongo Dog
		mongoHamburgerInsert := MongoHamburger{
			BurgerType:  postedHamburger.BurgerType,
			Condiments:  turnFoodArray(postedHamburger.Condiment),
			Calories:    postedHamburger.Calories,
			Name:        postedHamburger.Name,
			FoodID:      postedHamburger.FoodID,
			UserID:      postedHamburger.UserID,
			PhotoID:     0,
			PhotoSrc:    "",
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Collect Data for Mongo
		hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
		collectedUsers := []interface{}{mongoHamburgerInsert}
		//Insert Our Data
		insertManyResult, err2 := hamburgerCollection.InsertMany(theContext, collectedUsers)
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
			//Insert Hamburger data into User array
			userCollection := mongoClient.Database("superdbtest1").Collection("users")
			filterUserID := bson.M{"userid": postedHamburger.UserID}
			var foundUser AUser
			foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
			theErr := userCollection.FindOne(theContext, filterUserID).Decode(&foundUser)
			if theErr != nil {
				if strings.Contains(theErr.Error(), "no documents in result") {
					fmt.Printf("It's all good, this document didn't find this UserID: %v\n", postedHamburger.UserID)
				} else {
					fmt.Printf("DEBUG: We had an error finding a User for this hamburgerUserID: %v\n%v\n", postedHamburger.UserID,
						theErr)
				}
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
				foundUser.Hamburgers.Hamburgers = append(foundUser.Hamburgers.Hamburgers, mongoHamburgerInsert)
				successfulUserInsert := updateUser(foundUser) //Update this User with the new Hotdog Array
				if successfulUserInsert == true {
					fmt.Printf("This User's hamburgers was updated successfully: %v\n", foundUser.UserID)
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
					//UPDATE GLOBAL DATA FOR FILE UPLOAD
					awsuserID = postedHamburger.UserID
					awsfoodID = postedHamburger.FoodID
					awsphotoName = postedHamburger.Name
					awsfoodType = "HAMBURGER"
					awsdateUpdated = postedHamburger.DateCreated
					awsdateCreated = postedHamburger.DateUpdated
					fmt.Fprintf(w, string(dataJSON))
					fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
				} else {
					fmt.Printf("This User's hamburgers were NOT updated successfully: %v\n", foundUser.UserID)
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
				}
			}
		}
	} else {
		//Marshal the bad news
		theReturnData := returnData{
			SuccessMsg:        "Your food contained foul language, please re-type",
			ReturnedHamburger: MongoHamburger{},
			SuccessBool:       false,
		}
		dataJSON, err := json.Marshal(theReturnData)
		if err != nil {
			fmt.Println("There's an error marshalling this hamburger.")
			logWriter("There's an error marshalling.")
		}
		fmt.Fprintf(w, string(dataJSON))
	}
}
