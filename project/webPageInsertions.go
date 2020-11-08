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

//Define returned JSON
type returnDataHam struct {
	SuccessMsg        string         `json:"SuccessMsg"`
	ReturnedHamburger MongoHamburger `json:"ReturnedHamburger"`
	SuccessBool       bool           `json:"SuccessBool"`
}
type returnDataHot struct {
	SuccessMsg     string      `json:"SuccessMsg"`
	ReturnedHotDog MongoHotDog `json:"ReturnedHotDog"`
	SuccessBool    bool        `json:"SuccessBool"`
}

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
			theReturnData := returnDataHot{
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
				theReturnData := returnDataHot{
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
					theReturnData := returnDataHot{
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
					theReturnData := returnDataHot{
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
		theReturnData := returnDataHot{
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
	fmt.Println("DEBUG: Inserting Burger record in Mongo/SQL.")
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
			theReturnData := returnDataHam{
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
				theReturnData := returnDataHam{
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
					theReturnData := returnDataHam{
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
					fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
					fmt.Fprintf(w, string(dataJSON))
				} else {
					fmt.Printf("This User's hamburgers were NOT updated successfully: %v\n", foundUser.UserID)
					theReturnData := returnDataHam{
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
		theReturnData := returnDataHam{
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

func simpleFoodInsert(whatFood string, hotdog Hotdog, hamburger Hamburger) bool {
	goodFoodInsert := true //This determines if the insert was successful

	//Determine if hamburger or hotdog
	if strings.Contains(strings.ToUpper(whatFood), "HOTDOG") {
		//Protections for the Burger name
		if strings.Compare(hotdog.HotDogType, "DEBUGTYPE") == 0 {
			hotdog.HotDogType = "NONE"
		}
		/* CHECK TO SEE IF THIS Hamburger INSERTION IS CLEAN */
		canPost := true
		if containsLanguage(hotdog.HotDogType) {
			canPost = false
		} else if containsLanguage(hotdog.Condiment) {
			canPost = false
		} else if containsLanguage(hotdog.Name) {
			canPost = false
		} else {
			canPost = true
		}
		if canPost == true {

			//Insert into SQL
			theStatement := "INSERT INTO hot_dogs(TYPE, CONDIMENT, CALORIES, NAME, USER_ID," +
				" FOOD_ID, PHOTO_ID, PHOTO_SRC, DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?,?,?)"
			stmt, err := db.Prepare(theStatement)

			r, err := stmt.Exec(hotdog.HotDogType, hotdog.Condiment, hotdog.Calories, hotdog.Name, hotdog.UserID,
				hotdog.FoodID, hotdog.PhotoID, hotdog.PhotoSrc,
				hotdog.DateCreated, hotdog.DateUpdated)
			check(err)

			_, err2 := r.RowsAffected()
			check(err2)
			stmt.Close() //Close the SQL

			//Insert into Mongo
			//Declare Mongo Dog
			mongoHotDogInsert := MongoHotDog{
				HotDogType:  hotdog.HotDogType,
				Condiments:  turnFoodArray(hotdog.Condiment),
				Calories:    hotdog.Calories,
				Name:        hotdog.Name,
				FoodID:      hotdog.FoodID,
				UserID:      hotdog.UserID,
				PhotoID:     hotdog.PhotoID,
				PhotoSrc:    hotdog.PhotoSrc,
				DateCreated: hotdog.DateCreated,
				DateUpdated: hotdog.DateUpdated,
			}
			//Collect Data for Mongo
			hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
			collectedUsers := []interface{}{mongoHotDogInsert}
			//Insert Our Data
			_, err3 := hotdogCollection.InsertMany(theContext, collectedUsers)
			if err3 != nil {
				goodFoodInsert = false
				errMsg := "There's an error inserting this food into Mongo: " + err3.Error()
				logWriter(errMsg)
			} else {
				//Insert Hamburger data into User array
				userCollection := mongoClient.Database("superdbtest1").Collection("users")
				filterUserID := bson.M{"userid": hotdog.UserID}
				var foundUser AUser
				foundUser.DateUpdated = mongoHotDogInsert.DateCreated
				theErr := userCollection.FindOne(theContext, filterUserID).Decode(&foundUser)
				if theErr != nil {
					if strings.Contains(theErr.Error(), "no documents in result") {
						fmt.Printf("It's all good, this document didn't find this UserID: %v\n", hotdog.UserID)
					} else {
						goodFoodInsert = false
						msgErr := "We had an error finding a User for this food UserID: " + theErr.Error()
						fmt.Println(msgErr)
						logWriter(msgErr)
					}
				} else {
					foundUser.Hotdogs.Hotdogs = append(foundUser.Hotdogs.Hotdogs, mongoHotDogInsert)
					successfulUserInsert := updateUser(foundUser) //Update this User with the new Hotdog Array
					if successfulUserInsert == true {
						fmt.Printf("DEBUG: This User's hotdogs was updated successfully: %v\n", foundUser.UserID)
						//UPDATE GLOBAL DATA FOR FILE UPLOAD
						awsuserID = hotdog.UserID
						awsfoodID = hotdog.FoodID
						awsphotoName = hotdog.Name
						awsfoodType = "HOTDOG"
						awsdateUpdated = hotdog.DateCreated
						awsdateCreated = hotdog.DateUpdated
					} else {
						goodFoodInsert = false
						errMsg := "User food was NOT updated successfully"
						fmt.Println(errMsg)
						logWriter(errMsg)
					}
				}
			}
		} else {
			//Marshal the bad news
			errMsg := "Error with canPost in simpleFoodInsert: " + "False"
			logWriter(errMsg)
		}
	} else if strings.Contains(strings.ToUpper(whatFood), "HAMBURGER") {
		//Protections for the Burger name
		if strings.Compare(hamburger.BurgerType, "DEBUGTYPE") == 0 {
			hamburger.BurgerType = "NONE"
		}
		/* CHECK TO SEE IF THIS Hamburger INSERTION IS CLEAN */
		canPost := true
		if containsLanguage(hamburger.BurgerType) {
			canPost = false
		} else if containsLanguage(hamburger.Condiment) {
			canPost = false
		} else if containsLanguage(hamburger.Name) {
			canPost = false
		} else {
			canPost = true
		}
		if canPost == true {

			//Insert into SQL
			theStatement := "INSERT INTO hamburgers(TYPE, CONDIMENT, CALORIES, NAME, USER_ID, " +
				"FOOD_ID, PHOTO_ID, PHOTO_SRC,  DATE_CREATED, DATE_UPDATED) VALUES(?,?,?,?,?,?,?,?,?,?)"
			stmt, err := db.Prepare(theStatement)

			r, err := stmt.Exec(hamburger.BurgerType, hamburger.Condiment, hamburger.Calories, hamburger.Name,
				hamburger.UserID,
				hamburger.FoodID, hamburger.PhotoID, hamburger.PhotoSrc,
				hamburger.DateCreated, hamburger.DateUpdated)
			check(err)

			_, err2 := r.RowsAffected()
			check(err2)
			stmt.Close() //Close the SQL

			//Insert into Mongo
			//Declare Mongo Dog
			mongoHamburgerInsert := MongoHamburger{
				BurgerType:  hamburger.BurgerType,
				Condiments:  turnFoodArray(hamburger.Condiment),
				Calories:    hamburger.Calories,
				Name:        hamburger.Name,
				FoodID:      hamburger.FoodID,
				UserID:      hamburger.UserID,
				PhotoID:     hamburger.PhotoID,
				PhotoSrc:    hamburger.PhotoSrc,
				DateCreated: hamburger.DateCreated,
				DateUpdated: hamburger.DateUpdated,
			}
			//Collect Data for Mongo
			hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
			collectedUsers := []interface{}{mongoHamburgerInsert}
			//Insert Our Data
			_, err3 := hamburgerCollection.InsertMany(theContext, collectedUsers)
			if err3 != nil {
				goodFoodInsert = false
				errMsg := "There's an error inserting this food into Mongo: " + err3.Error()
				logWriter(errMsg)
			} else {
				//Insert Hamburger data into User array
				userCollection := mongoClient.Database("superdbtest1").Collection("users")
				filterUserID := bson.M{"userid": hamburger.UserID}
				var foundUser AUser
				foundUser.DateUpdated = mongoHamburgerInsert.DateCreated
				theErr := userCollection.FindOne(theContext, filterUserID).Decode(&foundUser)
				if theErr != nil {
					if strings.Contains(theErr.Error(), "no documents in result") {
						fmt.Printf("It's all good, this document didn't find this UserID: %v\n", hamburger.UserID)
					} else {
						goodFoodInsert = false
						msgErr := "We had an error finding a User for this food UserID: " + theErr.Error()
						fmt.Println(msgErr)
						logWriter(msgErr)
					}
				} else {
					foundUser.Hamburgers.Hamburgers = append(foundUser.Hamburgers.Hamburgers, mongoHamburgerInsert)
					successfulUserInsert := updateUser(foundUser) //Update this User with the new Hotdog Array
					if successfulUserInsert == true {
						fmt.Printf("DEBUG: This User's hamburgers was updated successfully: %v\n", foundUser.UserID)
						//UPDATE GLOBAL DATA FOR FILE UPLOAD
						awsuserID = hamburger.UserID
						awsfoodID = hamburger.FoodID
						awsphotoName = hamburger.Name
						awsfoodType = "HAMBURGER"
						awsdateUpdated = hamburger.DateCreated
						awsdateCreated = hamburger.DateUpdated
					} else {
						goodFoodInsert = false
						errMsg := "User food was NOT updated successfully"
						fmt.Println(errMsg)
						logWriter(errMsg)
					}
				}
			}
		} else {
			//Marshal the bad news
			errMsg := "Error with canPost in simpleFoodInsert: " + "False"
			logWriter(errMsg)
		}
	} else {
		goodFoodInsert = false
		errMsg := "Error, incorrect whatFood in simpleFoodInsert: " + whatFood
		logWriter(errMsg)
	}

	return goodFoodInsert
}

//func simplePhotoInsert(whatFood string)
