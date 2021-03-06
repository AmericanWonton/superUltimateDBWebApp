package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

var theContext context.Context
var mongoURI string //Connection string loaded

func connectDB() *mongo.Client {
	//Setup Mongo connection to Atlas Cluster
	theClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Errored getting mongo client: %v\n", err)
		log.Fatal(err)
	}
	theContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = theClient.Connect(theContext)
	if err != nil {
		fmt.Printf("Errored getting mongo client context: %v\n", err)
		log.Fatal(err)
	}
	//Double check to see if we've connected to the database
	err = theClient.Ping(theContext, readpref.Primary())
	if err != nil {
		fmt.Printf("Errored pinging MongoDB: %v\n", err)
		log.Fatal(err)
	}
	//List all available databases
	/*
		databases, err := theClient.ListDatabaseNames(theContext, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(databases)
	*/

	return theClient
}

func insertUsers(w http.ResponseWriter, req *http.Request) {
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("DEBUG: Error reading byte slice from Mongo Array: \n%v\n", err)
		fmt.Println(err)
	}

	//test Check Mongo Client
	err = mongoClient.Ping(theContext, readpref.Primary())
	if err != nil {
		fmt.Printf("Errored pinging MongoDB: %v\n", err)
		log.Fatal(err)
	}
	//Marshal it into our type
	var postedUsers TheUsers
	json.Unmarshal(bs, &postedUsers)
	fmt.Println(postedUsers)
	//Collect Data for Mongo
	user_collection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
	collectedUsers := []interface{}{}
	for x := 0; x < len(postedUsers.Users); x++ {
		collectedUsers = append(collectedUsers, postedUsers.Users[x])
	}
	//Insert Our Data
	insertManyResult, err := user_collection.InsertMany(context.TODO(), collectedUsers)
	if err != nil {
		fmt.Printf("Error inserting results: \n%v\n", err)
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
}

func updateUser(updatedUser AUser) bool {
	success := true
	theTimeNow := time.Now()
	updatedUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
	ic_collection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
	theFilter := bson.M{
		"userid": bson.M{
			"$eq": updatedUser.UserID, // check if bool field has value of 'false'
		},
	}
	updatedDocument := bson.M{
		"$set": bson.M{
			"username":    updatedUser.UserName,
			"password":    updatedUser.Password,
			"first":       updatedUser.First,
			"last":        updatedUser.Last,
			"role":        updatedUser.Role,
			"userid":      updatedUser.UserID,
			"datecreated": updatedUser.DateCreated,
			"dateupdated": updatedUser.DateUpdated,
			"hotdogs":     updatedUser.Hotdogs,
			"hamburgers":  updatedUser.Hamburgers,
		},
	}
	_, err := ic_collection.UpdateOne(
		theContext,
		theFilter,
		updatedDocument,
	)
	if err != nil {
		fmt.Printf("There was an error replacing one of our User documents: %v\n", err.Error())
		success = false
		log.Printf("%v\n", err)
	}
	return success
}

//This should return a single User, given a UserID from someone's Post
func getUserMongo(userID int) (AUser, bool, string) {
	var theUserReturned AUser //Initialize User to be returned after Mongo query
	successfulFind := true
	returnedErr := "" //Declare Error to be returned

	//Query for the User, given the userID for the User
	ic_collection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
	theFilter := bson.M{
		"userid": bson.M{
			"$eq": userID, // check if bool field has value of 'false'
		},
	}
	findOptions := options.FindOne()
	findUser := ic_collection.FindOne(theContext, theFilter, findOptions)
	if findUser.Err() != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			stringNum := strconv.Itoa(userID) //int to string conversion
			returnedErr = "For " + stringNum + ", no User was returned: " + err.Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
			successfulFind = false
		} else {
			stringNum := strconv.Itoa(userID) //int to string conversion
			returnedErr = "For " + stringNum + ", there was a Mongo Error: " + err.Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
			successfulFind = false
		}
		return theUserReturned, successfulFind, returnedErr //Return with errors
	} else {
		err := findUser.Decode(&theUserReturned)
		if err != nil {
			stringNum := strconv.Itoa(userID) //int to string conversion
			returnedErr = "For " + stringNum +
				", there was an error decoding document from Mongo: " + err.Error()
			fmt.Println(returnedErr)
			logWriter(returnedErr)
			successfulFind = false
			return theUserReturned, successfulFind, returnedErr
		} else {
			stringNum := strconv.Itoa(userID) //int to string conversion
			returnedErr = "For " + stringNum +
				", User should be successfully decoded."
			fmt.Println(returnedErr)
			logWriter(returnedErr)
			successfulFind = true
			return theUserReturned, successfulFind, returnedErr
		}
	}
}

func insertHotDogs(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting hotdog records in Mongo.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Here is our byte slice as a string for JSON: \n\n%v\n", string(bs))
	//Marshal it into our type
	var postedHotDogs MongoHotDogs
	json.Unmarshal(bs, &postedHotDogs)

	hotdog_collection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
	collectedHDogs := []interface{}{}
	for x := 0; x < len(postedHotDogs.Hotdogs); x++ {
		collectedHDogs = append(collectedHDogs, postedHotDogs.Hotdogs[x])
	}
	//Insert Our Data
	insertManyResult, err := hotdog_collection.InsertMany(context.TODO(), collectedHDogs)
	if err != nil {
		fmt.Printf("Error inserting results: \n%v\n", err)
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
}

func insertHotDogMongo(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting hotdog record in Mongo.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Here is our byte slice as a string for JSON: \n\n%v\n", string(bs))
	//Marshal it into our type
	var postedHotDog Hotdog
	json.Unmarshal(bs, &postedHotDog)

	//Protections for the hotdog name
	if strings.Compare(postedHotDog.HotDogType, "DEBUGTYPE") == 0 {
		postedHotDog.HotDogType = "NONE"
	}
	//Change into Mongo Hotdog
	theTimeNow := time.Now()
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
	insertManyResult, err := user_collection.InsertMany(theContext, collectedUsers)
	//Define the data to return
	type returnData struct {
		SuccessMsg     string      `json:"SuccessMsg"`
		ReturnedHotDog MongoHotDog `json:"ReturnedHotDog"`
		SuccessBool    bool        `json:"SuccessBool"`
	}
	if err != nil {
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
}

func insertHamburgers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting Hamburger records in Mongo.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Here is our byte slice as a string for JSON: \n\n%v\n", string(bs))
	//Marshal it into our type
	var postedHamburgers MongoHamburgers
	json.Unmarshal(bs, &postedHamburgers)

	hamburger_collection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
	collectedHamburgers := []interface{}{}
	for x := 0; x < len(postedHamburgers.Hamburgers); x++ {
		collectedHamburgers = append(collectedHamburgers, postedHamburgers.Hamburgers[x])
	}
	//Insert Our Data
	insertManyResult, err := hamburger_collection.InsertMany(context.TODO(), collectedHamburgers)
	if err != nil {
		fmt.Printf("Error inserting results: \n%v\n", err)
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
}

func insertHamburgerMongo(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting Hamburger record in Mongo.")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Here is our byte slice as a string for JSON: \n\n%v\n", string(bs))
	//Marshal it into our type
	var postedHamburger Hamburger
	json.Unmarshal(bs, &postedHamburger)

	//Protections for the hotdog name
	if strings.Compare(postedHamburger.BurgerType, "DEBUGTYPE") == 0 {
		postedHamburger.BurgerType = "NONE"
	}
	//Change into Mongo Hotdog
	theTimeNow := time.Now()
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
	insertManyResult, err := user_collection.InsertMany(context.TODO(), collectedUsers)
	if err != nil {
		fmt.Printf("Error inserting results: \n%v\n", err)
		fmt.Fprint(w, failureMessage)
		log.Fatal(err)
	} else {
		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
		fmt.Fprint(w, successMessage)
	}
}

func foodUpdateMongo(w http.ResponseWriter, req *http.Request) {
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

	//Determine if this is a hotdog or hamburger update
	if thefoodUpdate.FoodType == "hotdog" {
		theTimeNow := time.Now()
		var hotDogUpdate Hotdog = thefoodUpdate.TheHotDog
		updatedHotDogMongo := MongoHotDog{
			HotDogType:  hotDogUpdate.HotDogType,
			Condiments:  turnFoodArray(hotDogUpdate.Condiment),
			Calories:    hotDogUpdate.Calories,
			Name:        hotDogUpdate.Name,
			FoodID:      thefoodUpdate.FoodID,
			UserID:      hotDogUpdate.UserID,
			PhotoID:     hotDogUpdate.PhotoID,
			PhotoSrc:    hotDogUpdate.PhotoSrc,
			DateCreated: hotDogUpdate.DateCreated,
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		/* CHECK TO SEE IF THIS Hamburger INSERTION IS CLEAN */
		canPost := true
		if containsLanguage(updatedHotDogMongo.HotDogType) {
			canPost = false
		} else if containsLanguage(updatedHotDogMongo.Name) {
			canPost = false
		} else {
			canPost = true
		}
		for j := 0; j < len(updatedHotDogMongo.Condiments); j++ {
			if containsLanguage(updatedHotDogMongo.Condiments[j]) {
				canPost = false
			}
		}
		if canPost == true {
			//Add updatedHotDog to Document collection for Hotdogs
			ic_collection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
			theFilter := bson.M{
				"foodid": bson.M{
					"$eq": thefoodUpdate.FoodID, // check if bool field has value of 'false'
				},
				"userid": bson.M{
					"$eq": updatedHotDogMongo.UserID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"hotdogtype":  updatedHotDogMongo.HotDogType,
					"condiments":  updatedHotDogMongo.Condiments,
					"calories":    updatedHotDogMongo.Calories,
					"name":        updatedHotDogMongo.Name,
					"foodid":      thefoodUpdate.FoodID,
					"userid":      updatedHotDogMongo.UserID,
					"photoid":     updatedHotDogMongo.PhotoID,
					"photosrc":    updatedHotDogMongo.PhotoSrc,
					"datecreated": updatedHotDogMongo.DateCreated,
					"dateupdated": updatedHotDogMongo.DateUpdated,
				},
			}
			updateResult, err := ic_collection.UpdateOne(theContext, theFilter, updatedDocument)
			if err != nil {
				fmt.Printf("Error updating the hotdog: %v\n\n", err.Error())
				fmt.Fprintln(w, 3) //Failure Response Response
			} else {
				//Our new UpdateResult
				fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
				/* NOW UPDATE FROM USER COLLECITON */
				userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
				theFilter := bson.M{
					"userid": bson.M{
						"$eq": thefoodUpdate.TheHotDog.UserID, // check if bool field has value of 'false'
					},
				}
				var foundUser AUser
				theTimeNow := time.Now()
				foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				theErr := userCollection.FindOne(theContext, theFilter).Decode(&foundUser)
				if theErr != nil {
					if strings.Contains(theErr.Error(), "no documents in result") {
						fmt.Printf("It's all good, this document wasn't found for User,(%v) and our ID is clean.\n",
							thefoodUpdate.TheHotDog.UserID)
					} else {
						fmt.Printf("DEBUG: We have another error for finding a UserID: %v \n%v\n",
							thefoodUpdate.TheHotDog.UserID, theErr)
					}
				} else {
					/* UPDATE USER HOTDOGS */
					fmt.Printf("Finding a hotdog,(%v), to update for this User:\n %v\n", thefoodUpdate.FoodID,
						foundUser)
					newHDogSlice := []MongoHotDog{}
					for i := 0; i < len(foundUser.Hotdogs.Hotdogs); i++ {
						if foundUser.Hotdogs.Hotdogs[i].FoodID == thefoodUpdate.FoodID {
							fmt.Printf("Not adding this food, using new Hotdog instead.\n")
							newHDogSlice = append(newHDogSlice, updatedHotDogMongo)
						} else {
							newHDogSlice = append(newHDogSlice, foundUser.Hotdogs.Hotdogs[i])
						}
					}
					foundUser.Hotdogs.Hotdogs = newHDogSlice
					fmt.Printf("We are sending the User data to update: %v\n", foundUser)
					updateUser(foundUser)

					fmt.Fprintln(w, 1) //Success Response
				}
			}
		} else {
			fmt.Printf("Food contains derogatory terms...no update occuring.\n")
			fmt.Fprintln(w, 4) //Success Response
		}
	} else if thefoodUpdate.FoodType == "hamburger" {
		theTimeNow := time.Now()
		var hamburgerUpdate Hamburger = thefoodUpdate.TheHamburger
		updatedHamburgerMongo := MongoHamburger{
			BurgerType:  hamburgerUpdate.BurgerType,
			Condiments:  turnFoodArray(hamburgerUpdate.Condiment),
			Calories:    hamburgerUpdate.Calories,
			Name:        hamburgerUpdate.Name,
			FoodID:      thefoodUpdate.FoodID,
			UserID:      hamburgerUpdate.UserID,
			PhotoID:     hamburgerUpdate.PhotoID,
			PhotoSrc:    hamburgerUpdate.PhotoSrc,
			DateCreated: hamburgerUpdate.DateCreated,
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		/* CHECK TO SEE IF THIS Hamburger INSERTION IS CLEAN */
		canPost := true
		if containsLanguage(updatedHamburgerMongo.BurgerType) {
			canPost = false
		} else if containsLanguage(updatedHamburgerMongo.Name) {
			canPost = false
		} else {
			canPost = true
		}
		for j := 0; j < len(updatedHamburgerMongo.Condiments); j++ {
			if containsLanguage(updatedHamburgerMongo.Condiments[j]) {
				canPost = false
			}
		}
		if canPost == true {
			//Add updatedHotDog to Document collection for Hotdogs
			ic_collection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
			theFilter := bson.M{
				"foodid": bson.M{
					"$eq": thefoodUpdate.FoodID, // check if bool field has value of 'false'
				},
				"userid": bson.M{
					"$eq": updatedHamburgerMongo.UserID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"burgertype":  updatedHamburgerMongo.BurgerType,
					"condiments":  updatedHamburgerMongo.Condiments,
					"calories":    updatedHamburgerMongo.Calories,
					"name":        updatedHamburgerMongo.Name,
					"foodid":      thefoodUpdate.FoodID,
					"userid":      updatedHamburgerMongo.UserID,
					"photoid":     updatedHamburgerMongo.PhotoID,
					"photosrc":    updatedHamburgerMongo.PhotoSrc,
					"datecreated": updatedHamburgerMongo.DateCreated,
					"dateupdated": updatedHamburgerMongo.DateUpdated,
				},
			}

			updateResult, err := ic_collection.UpdateOne(context.TODO(), theFilter, updatedDocument)
			if err != nil {
				fmt.Printf("There was an error updating hamburgers: %v\n\n", err)
				fmt.Fprintln(w, 3) //Failure Response Response
			} else {
				//Our new UpdateResult
				fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

				/* NOW UPDATE FROM USER COLLECITON */
				userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
				theFilter := bson.M{
					"userid": bson.M{
						"$eq": thefoodUpdate.TheHamburger.UserID, // check if bool field has value of 'false'
					},
				}
				var foundUser AUser
				theTimeNow := time.Now()
				foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				theErr := userCollection.FindOne(theContext, theFilter).Decode(&foundUser)
				if theErr != nil {
					if strings.Contains(theErr.Error(), "no documents in result") {
						fmt.Printf("It's all good, this document wasn't found for User,(%v) and our ID is clean.\n",
							thefoodUpdate.TheHotDog.UserID)
					} else {
						fmt.Printf("DEBUG: We have another error for finding a unique UserID: %v \n%v\n",
							thefoodUpdate.TheHamburger.UserID, theErr)
					}
				} else {
					fmt.Printf("Finding a hamburger,(%v), to update for this User:\n %v\n", thefoodUpdate.FoodID,
						foundUser)
					/* UPDATE USER HAMBURGERS */
					newHamburgerSlice := []MongoHamburger{}
					for i := 0; i < len(foundUser.Hamburgers.Hamburgers); i++ {
						if foundUser.Hamburgers.Hamburgers[i].FoodID == thefoodUpdate.FoodID {
							fmt.Printf("DEBUG: We've found the food to update, skipping it and appending the new hamburger\n")
							newHamburgerSlice = append(newHamburgerSlice, updatedHamburgerMongo)
						} else {
							newHamburgerSlice = append(newHamburgerSlice, foundUser.Hamburgers.Hamburgers[i])
						}
					}
					foundUser.Hamburgers.Hamburgers = newHamburgerSlice
					fmt.Printf("DEBUG: Sending updated hamburgers to User for updating:\n %v\n", foundUser)
					updateUser(foundUser)
					fmt.Fprintln(w, 1) //Success Response
				}
			}
		} else {
			fmt.Printf("Food contains derogatory terms...no update occuring.\n")
			fmt.Fprintln(w, 4) //Success Response
		}
	} else {
		fmt.Printf("Unexpected JSON from update Food function: %v\n", thefoodUpdate.FoodType)
		fmt.Fprintln(w, 3)
	}
}

//A simple version of updating food in Mongo
func mongoUpdateFood(whichFood string, hamburger MongoHamburger, hotdog MongoHotDog, fileUpdate string) bool {
	goodUpdate := true

	//Check to see which food we're updating
	if strings.Contains(strings.ToUpper(whichFood), "HOTDOG") {
		canPost := true
		if containsLanguage(hotdog.HotDogType) {
			canPost = false
		} else if containsLanguage(hotdog.Name) {
			canPost = false
		} else {
			canPost = true
		}
		for j := 0; j < len(hotdog.Condiments); j++ {
			if containsLanguage(hotdog.Condiments[j]) {
				canPost = false
			}
		}
		if canPost == true {
			//Declare hotdog for usage
			//Find the food item
			hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
			//Query Mongo for all hotdogs for a User
			theFilter := bson.M{"foodid": hotdog.FoodID}
			findOptions := options.Find()
			curHDog, err := hotdogCollection.Find(theContext, theFilter, findOptions)
			if err != nil {
				if strings.Contains(err.Error(), "no documents in result") {
					goodUpdate = false
					fmt.Printf("No documents were returned for hotdogs in MongoDB: %v\n", err.Error())
					logWriter("No documents were returned in MongoDB for Hotdogs: " + err.Error())
				} else {
					goodUpdate = false
					fmt.Printf("There was an error returning hotdogs for this User, %v: %v\n", hotdog.FoodID, err.Error())
					logWriter("There was an error returning hotdogs in Mongo for this User " + err.Error())
				}
			}
			//Loop over query results and fill hotdogs array
			var returnedDog MongoHotDog
			for curHDog.Next(theContext) {
				// create a value into which the single document can be decoded
				err := curHDog.Decode(&returnedDog)
				if err != nil {
					fmt.Printf("Error decoding hotdogs in MongoDB for this User, %v: %v\n", hotdog.UserID, err.Error())
					logWriter("Error decoding hotdogs in MongoDB: " + err.Error())
				}
			}
			// Close the cursor once finished
			curHDog.Close(theContext)

			//update food
			theTimeNow := time.Now()
			ic_collection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
			theFilterTwo := bson.M{
				"foodid": bson.M{
					"$eq": hotdog.FoodID, // check if bool field has value of 'false'
				},
			}

			updatedDocument := bson.M{}

			//Update document based on file submission
			if strings.Contains(strings.ToLower(fileUpdate), strings.ToLower("no_photo")) {
				fmt.Printf("DEBUG: fileUpdate in Mongo food update is: %v\n", fileUpdate)
				updatedDocument = bson.M{
					"$set": bson.M{
						"hotdogtype":  hotdog.HotDogType,
						"condiments":  hotdog.Condiments,
						"calories":    hotdog.Calories,
						"name":        hotdog.Name,
						"foodid":      returnedDog.FoodID,
						"userid":      returnedDog.UserID,
						"photoid":     returnedDog.PhotoID,
						"photosrc":    returnedDog.PhotoSrc,
						"datecreated": returnedDog.DateCreated,
						"dateupdated": theTimeNow.Format("2006-01-02 15:04:05"),
					},
				}
			} else if strings.Contains(strings.ToLower(fileUpdate), strings.ToLower("file_update")) {
				fmt.Printf("DEBUG: fileUpdate in Mongo food update is: %v\n", fileUpdate)
				updatedDocument = bson.M{
					"$set": bson.M{
						"hotdogtype":  hotdog.HotDogType,
						"condiments":  hotdog.Condiments,
						"calories":    hotdog.Calories,
						"name":        hotdog.Name,
						"foodid":      hotdog.FoodID,
						"userid":      returnedDog.UserID,
						"photoid":     hotdog.PhotoID,
						"photosrc":    hotdog.PhotoSrc,
						"datecreated": returnedDog.DateCreated,
						"dateupdated": theTimeNow.Format("2006-01-02 15:04:05"),
					},
				}
			} else {
				fmt.Printf("DEBUG: fileUpdate in Mongo food update is: %v\n", fileUpdate)
				goodUpdate = false
				errMsg := "Incorrect fileUpdate in mongoUpdateFood: " + fileUpdate
				fmt.Println(errMsg)
				logWriter(errMsg)
				return goodUpdate
			}

			fmt.Printf("DEBUG: Here is our bson.M in mongo food update: %v\n", updatedDocument)
			_, err2 := ic_collection.UpdateOne(
				theContext,
				theFilterTwo,
				updatedDocument,
			)
			if err2 != nil {
				goodUpdate = false
				errMsg := "Failed to update food in mongoUpdateFood: " + err2.Error()
				fmt.Println(errMsg)
				logWriter(errMsg)
			} else {
				/* NOW UPDATE FROM USER COLLECITON */
				userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
				theFilter := bson.M{
					"userid": bson.M{
						"$eq": hotdog.UserID, // check if bool field has value of 'false'
					},
				}
				var foundUser AUser
				theTimeNow := time.Now()
				foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				theErr := userCollection.FindOne(theContext, theFilter).Decode(&foundUser)
				if theErr != nil {
					if strings.Contains(theErr.Error(), "no documents in result") {
						goodUpdate = false
						fmt.Printf("It's all good, this document wasn't found for User,(%v) and our ID is clean.\n",
							hotdog.UserID)
					} else {
						goodUpdate = false
						fmt.Printf("DEBUG: We have another error for finding a UserID: %v \n%v\n",
							hotdog.UserID, theErr)
					}
				} else {
					/* UPDATE USER FOOD */
					updatedHotDogMongo := MongoHotDog{
						HotDogType:  hotdog.HotDogType,
						Condiments:  hotdog.Condiments,
						Calories:    hotdog.Calories,
						Name:        hotdog.Name,
						FoodID:      hotdog.FoodID,
						UserID:      hotdog.UserID,
						PhotoID:     hotdog.PhotoID,
						PhotoSrc:    hotdog.PhotoSrc,
						DateCreated: returnedDog.DateCreated,
						DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
					}
					newHDogSlice := []MongoHotDog{}
					for i := 0; i < len(foundUser.Hotdogs.Hotdogs); i++ {
						if foundUser.Hotdogs.Hotdogs[i].FoodID == hotdog.FoodID {
							newHDogSlice = append(newHDogSlice, updatedHotDogMongo)
						} else {
							newHDogSlice = append(newHDogSlice, foundUser.Hotdogs.Hotdogs[i])
						}
					}
					foundUser.Hotdogs.Hotdogs = newHDogSlice
					updateUser(foundUser)
				}
			}
		} else {
			goodUpdate = false
			errMsg := "Incorrect whichFood variable entered in mongoUpdateFood: " + whichFood
			logWriter(errMsg)
			fmt.Println(errMsg)
		}
	} else if strings.Contains(strings.ToUpper(whichFood), "HAMBURGER") {
		canPost := true
		if containsLanguage(hamburger.BurgerType) {
			canPost = false
		} else if containsLanguage(hamburger.Name) {
			canPost = false
		} else {
			canPost = true
		}
		for j := 0; j < len(hamburger.Condiments); j++ {
			if containsLanguage(hamburger.Condiments[j]) {
				canPost = false
			}
		}
		if canPost == true {
			//Declare hotdog for usage
			//Find the food item
			hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
			//Query Mongo for all hotdogs for a User
			theFilter := bson.M{"foodid": hamburger.FoodID}
			findOptions := options.Find()
			curHam, err := hamburgerCollection.Find(theContext, theFilter, findOptions)
			if err != nil {
				if strings.Contains(err.Error(), "no documents in result") {
					goodUpdate = false
					fmt.Printf("No documents were returned for hamburgers in MongoDB: %v\n", err.Error())
					logWriter("No documents were returned in MongoDB for Hotdogs: " + err.Error())
				} else {
					goodUpdate = false
					fmt.Printf("There was an error returning hamburgers for this User, %v: %v\n", hotdog.FoodID, err.Error())
					logWriter("There was an error returning hamburgers in Mongo for this User " + err.Error())
				}
			}
			//Loop over query results and fill hotdogs array
			var returnedHamburger MongoHamburger
			for curHam.Next(theContext) {
				// create a value into which the single document can be decoded
				err := curHam.Decode(&returnedHamburger)
				if err != nil {
					fmt.Printf("Error decoding hamburgers in MongoDB for this User, %v: %v\n", hotdog.UserID, err.Error())
					logWriter("Error decoding hamburgers in MongoDB: " + err.Error())
				}
			}
			// Close the cursor once finished
			curHam.Close(theContext)

			//update food
			theTimeNow := time.Now()
			ic_collection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
			theFilterTwo := bson.M{
				"foodid": bson.M{
					"$eq": hamburger.FoodID, // check if bool field has value of 'false'
				},
			}

			updatedDocument := bson.M{}
			//Update document based on file submission
			if strings.Contains(strings.ToLower(fileUpdate), strings.ToLower("no_photo")) {
				updatedDocument = bson.M{
					"$set": bson.M{
						"burgertype":  hamburger.BurgerType,
						"condiments":  hamburger.Condiments,
						"calories":    hamburger.Calories,
						"name":        hamburger.Name,
						"foodid":      returnedHamburger.FoodID,
						"userid":      returnedHamburger.UserID,
						"photoid":     returnedHamburger.PhotoID,
						"photosrc":    returnedHamburger.PhotoSrc,
						"datecreated": returnedHamburger.DateCreated,
						"dateupdated": theTimeNow.Format("2006-01-02 15:04:05"),
					},
				}
			} else if strings.Contains(strings.ToLower(fileUpdate), strings.ToLower("file_update")) {
				updatedDocument = bson.M{
					"$set": bson.M{
						"burgertype":  hamburger.BurgerType,
						"condiments":  hamburger.Condiments,
						"calories":    hamburger.Calories,
						"name":        hamburger.Name,
						"foodid":      returnedHamburger.FoodID,
						"userid":      returnedHamburger.UserID,
						"photoid":     hamburger.PhotoID,
						"photosrc":    hamburger.PhotoSrc,
						"datecreated": returnedHamburger.DateCreated,
						"dateupdated": theTimeNow.Format("2006-01-02 15:04:05"),
					},
				}
			} else {
				goodUpdate = false
				errMsg := "Incorrect fileUpdate in mongoUpdateFood: " + fileUpdate
				fmt.Println(errMsg)
				logWriter(errMsg)
				return goodUpdate
			}

			_, err2 := ic_collection.UpdateOne(
				theContext,
				theFilterTwo,
				updatedDocument,
			)
			if err2 != nil {
				goodUpdate = false
				errMsg := "Failed to update food in mongoUpdateFood: " + err2.Error()
				fmt.Println(errMsg)
				logWriter(errMsg)
			} else {
				/* NOW UPDATE FROM USER COLLECITON */
				userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
				theFilter := bson.M{
					"userid": bson.M{
						"$eq": hamburger.UserID, // check if bool field has value of 'false'
					},
				}
				var foundUser AUser
				theTimeNow := time.Now()
				foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				theErr := userCollection.FindOne(theContext, theFilter).Decode(&foundUser)
				if theErr != nil {
					if strings.Contains(theErr.Error(), "no documents in result") {
						goodUpdate = false
						fmt.Printf("It's all good, this document wasn't found for User,(%v) and our ID is clean.\n",
							hamburger.UserID)
					} else {
						goodUpdate = false
						fmt.Printf("DEBUG: We have another error for finding a UserID: %v \n%v\n",
							hamburger.UserID, theErr)
					}
				} else {
					/* UPDATE USER FOOD */
					updatedHamMongo := MongoHamburger{
						BurgerType:  hamburger.BurgerType,
						Condiments:  hamburger.Condiments,
						Calories:    hamburger.Calories,
						Name:        hamburger.Name,
						FoodID:      hamburger.FoodID,
						UserID:      hamburger.UserID,
						PhotoID:     hamburger.PhotoID,
						PhotoSrc:    hamburger.PhotoSrc,
						DateCreated: returnedHamburger.DateCreated,
						DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
					}
					newHamSlice := []MongoHamburger{}
					for i := 0; i < len(foundUser.Hamburgers.Hamburgers); i++ {
						if foundUser.Hotdogs.Hotdogs[i].FoodID == hotdog.FoodID {
							newHamSlice = append(newHamSlice, updatedHamMongo)
						} else {
							newHamSlice = append(newHamSlice, foundUser.Hamburgers.Hamburgers[i])
						}
					}
					foundUser.Hamburgers.Hamburgers = newHamSlice
					updateUser(foundUser)
				}
			}
		} else {
			goodUpdate = false
			errMsg := "Incorrect whichFood variable entered in mongoUpdateFood: " + whichFood
			logWriter(errMsg)
			fmt.Println(errMsg)
		}
	} else {
		goodUpdate = false
		errMsg := "Incorrect whichFood variable entered in sqlUpdateFood: " + whichFood
		logWriter(errMsg)
		fmt.Println(errMsg)
	}

	return goodUpdate
}

//DEBUG: Work in Progress
func foodDeleteMongo(w http.ResponseWriter, req *http.Request) {
	//Our Food deletion struct
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
	if theFoodDeletion.FoodType == "hotdog" {
		fmt.Printf("Deleting the following hotdog,(%v), from the following User: %v\n", theFoodDeletion.FoodID,
			theFoodDeletion.UserID)
		/* FIRST DELETE FROM HOTDOG COLLECTION*/
		hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
		deletes := []bson.M{
			{"UserID": theFoodDeletion.UserID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"UserID": bson.M{
			"$eq": theFoodDeletion.UserID,
		}}, bson.M{"foodid": bson.M{
			"$eq": theFoodDeletion.FoodID,
		}},
		)

		// create the slice of write models
		var writes []mongo.WriteModel

		for _, del := range deletes {
			model := mongo.NewDeleteManyModel().SetFilter(del)
			writes = append(writes, model)
		}

		// run bulk write
		res, err := hotdogCollection.BulkWrite(theContext, writes)
		if err != nil {
			logWriter("Error writing Mongo Delete Statement")
			logWriter("\n")
			logWriter(err.Error())
		}
		//Print Results
		fmt.Printf("Deleted the following documents: %v\n", res.DeletedCount)
		logWriter("Deleted the following documents: " + string(res.DeletedCount) + "\n")

		/* NOW DELETE FROM USER COLLECITON */
		userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
		theFilter := bson.M{
			"userid": bson.M{
				"$eq": theFoodDeletion.UserID, // check if bool field has value of 'false'
			},
		}
		var foundUser AUser
		theTimeNow := time.Now()
		foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		theErr := userCollection.FindOne(theContext, theFilter).Decode(&foundUser)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User,(%v) and our ID is clean.\n",
					theFoodDeletion.UserID)
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: %v \n%v\n",
					theFoodDeletion.UserID, theErr)
			}
		} else {
			fmt.Printf("DEBUG: We found the User and we'll delete the Hotdogs: %v\n", foundUser.Hotdogs.Hotdogs)
			//Remove Hotdog
			newHDogSlice := []MongoHotDog{}
			if len(foundUser.Hotdogs.Hotdogs) == 1 {
				newHDogSlice = nil
			} else {
				for j := 0; j < len(foundUser.Hotdogs.Hotdogs); j++ {
					fmt.Printf("DEBUG: Here is the %v foodID: %v\n", j, foundUser.Hotdogs.Hotdogs[j].FoodID)
					if foundUser.Hotdogs.Hotdogs[j].FoodID == theFoodDeletion.FoodID {
						fmt.Printf("DEBUG: We ignoring this hotdog and adding the others to our slice.\n")
					} else {
						newHDogSlice = append(newHDogSlice, foundUser.Hotdogs.Hotdogs[j])
					}
				}
			}
			//Update User
			foundUser.Hotdogs.Hotdogs = newHDogSlice
			fmt.Printf("Giving new User data for deletion:\n %v\n", foundUser)
			updateUser(foundUser)

			fmt.Fprintln(w, 1)
		}
	} else if theFoodDeletion.FoodType == "hamburger" {
		fmt.Printf("Deleting the following hamburger,(%v), from the following User: %v\n", theFoodDeletion.FoodID,
			theFoodDeletion.UserID)
		hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
		deletes := []bson.M{
			{"UserID": theFoodDeletion.UserID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"UserID": bson.M{
			"$eq": theFoodDeletion.UserID,
		}}, bson.M{"foodid": bson.M{
			"$eq": theFoodDeletion.FoodID,
		}},
		)

		// create the slice of write models
		var writes []mongo.WriteModel

		for _, del := range deletes {
			model := mongo.NewDeleteManyModel().SetFilter(del)
			writes = append(writes, model)
		}

		// run bulk write
		res, err := hamburgerCollection.BulkWrite(theContext, writes)
		if err != nil {
			logWriter("Error writing Mongo Delete Statement")
			logWriter("\n")
			logWriter(err.Error())
		}
		//Print Results
		fmt.Printf("Deleted the following documents: %v\n", res.DeletedCount)
		logWriter("Deleted the following documents: " + string(res.DeletedCount) + "\n")

		/* NOW DELETE FROM USER COLLECITON */
		userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
		theFilter := bson.M{
			"userid": bson.M{
				"$eq": theFoodDeletion.UserID, // check if bool field has value of 'false'
			},
		}
		var foundUser AUser
		theTimeNow := time.Now()
		foundUser.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		theErr := userCollection.FindOne(theContext, theFilter).Decode(&foundUser)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User,(%v) and our ID is clean.\n",
					theFoodDeletion.UserID)
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: %v \n%v\n",
					theFoodDeletion.UserID, theErr)
			}
		} else {
			//Remove Hamburger
			hamburgerSlice := []MongoHamburger{}
			if len(foundUser.Hamburgers.Hamburgers) == 1 {
				hamburgerSlice = nil
			} else {
				for j := 0; j < len(foundUser.Hamburgers.Hamburgers); j++ {
					fmt.Printf("DEBUG: Here is the %v foodID: %v\n", j, foundUser.Hamburgers.Hamburgers[j].FoodID)
					if foundUser.Hamburgers.Hamburgers[j].FoodID == theFoodDeletion.FoodID {
						fmt.Printf("DEBUG: We will not include this Hamburger in the slice.\n")
					} else {
						hamburgerSlice = append(hamburgerSlice, foundUser.Hamburgers.Hamburgers[j])
					}
				}
			}
			//Update User
			foundUser.Hamburgers.Hamburgers = hamburgerSlice
			fmt.Printf("DEBUG: Updating the User with the deleted Hamburger(s):\n%v\n", foundUser)
			updateUser(foundUser)

			fmt.Fprintln(w, 2)
		}
	} else {
		fmt.Fprintln(w, 3)
	}
}

//This should give a random id value to both food groups
func randomIDCreation() int {
	finalID := 0        //The final, unique ID to return to the food/user
	randInt := 0        //The random integer added onto ID
	randIntString := "" //The integer built through a string...
	min, max := 0, 9    //The min and Max value for our randInt
	foundID := false
	for foundID == false {
		randInt = 0
		randIntString = ""
		//Create the random number, convert it to string
		for i := 0; i < 8; i++ {
			randInt = rand.Intn(max-min) + min
			randIntString = randIntString + strconv.Itoa(randInt)
		}
		//Once we have a string of numbers, we can convert it back to an integer
		theID, err := strconv.Atoi(randIntString)
		if err != nil {
			fmt.Printf("We got an error converting a string back to a number, %v\n", err)
			fmt.Printf("Here is randInt: %v\n and randIntString: %v\n", randInt, randIntString)
			fmt.Println(err)
			log.Fatal(err)
		}
		//Search all our collections to see if this UserID is unique
		canExit := []bool{true, true, true, true}
		fmt.Printf("DEBUG: We are going to see if this ID is in our food or User DBs: %v\n", theID)
		//User collection
		userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
		var testAUser AUser
		theErr := userCollection.FindOne(theContext, bson.M{"userid": theID}).Decode(&testAUser)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User and our ID is clean.\n")
				canExit[0] = true
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit[0] = false
				log.Fatal(theErr)
			}
		}
		//Check hotdog collection
		hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
		var testHotdog MongoHotDog
		//Give 0 values to determine if these IDs are found
		theFilter := bson.M{
			"$or": []interface{}{
				bson.M{"userid": theID},
				bson.M{"foodid": theID},
			},
		}
		theErr = hotdogCollection.FindOne(theContext, theFilter).Decode(&testHotdog)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User/Hotdog and our ID is clean.\n")
				canExit[1] = true
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit[1] = false
			}
		}
		//Check hamburger collection
		hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
		var testBurger MongoHamburger
		//Give 0 values to determine if these IDs are found
		theFilter2 := bson.M{
			"$or": []interface{}{
				bson.M{"userid": theID},
				bson.M{"foodid": theID},
			},
		}
		theErr = hamburgerCollection.FindOne(theContext, theFilter2).Decode(&testBurger)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				canExit[2] = true
				fmt.Printf("It's all good, this document wasn't found for User/hamburger and our ID is clean.\n")
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit[2] = false
			}
		}
		//Check photo collection
		photoCollection := mongoClient.Database("superdbtest1").Collection("userphotos") //Here's our collection
		var userPhoto UserPhoto
		//Give 0 values to determine if these IDs are found
		theFilter3 := bson.M{
			"$or": []interface{}{
				bson.M{"userid": theID},
				bson.M{"foodid": theID},
				bson.M{"photoid": theID},
			},
		}
		theErr = photoCollection.FindOne(theContext, theFilter3).Decode(&userPhoto)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				canExit[2] = true
				fmt.Printf("It's all good, this document wasn't found for User/hamburger/photo and our ID is clean.\n")
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique PhotoID: \n%v\n", theErr)
				canExit[2] = false
			}
		}
		//Check photo collection
		messageCollection := mongoClient.Database("messageboard").Collection("messages") //Here's our collection
		var aMessage Message
		//Give 0 values to determine if these IDs are found
		theFilter4 := bson.M{
			"$or": []interface{}{
				bson.M{"messageid": theID},
				bson.M{"userid": theID},
			},
		}
		theErr = messageCollection.FindOne(theContext, theFilter4).Decode(&aMessage)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				canExit[3] = true
				fmt.Printf("It's all good, this document wasn't found for User/hamburger/photo/message and our ID is clean.\n")
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique PhotoID: \n%v\n", theErr)
				canExit[3] = false
			}
		}
		//Final check to see if we can exit this loop
		if canExit[0] == true && canExit[1] == true && canExit[2] == true && canExit[3] == true {
			finalID = theID
			foundID = true
		} else {
			foundID = false
		}
	}

	return finalID
}

//Should give a random id value to use for both food groups...good for Mongo AND SQL insertion.
func randomIDCreationAPI(w http.ResponseWriter, req *http.Request) {
	//Our Food deletion struct
	type jsonRecieved struct {
		stringReceived string `json:"stringReceived"`
	}
	fmt.Println("DEBUG: We're creating a randomID thru this API...")
	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var ourJSONInfo jsonRecieved
	json.Unmarshal(bs, &ourJSONInfo)

	type returnJSON struct {
		SuccessMsg     string `json:"SuccessMsg"`
		FoodIDReturned int    `json:"FoodIDReturned"`
		SuccessBool    bool   `json:"SuccessBool"`
	}

	theID := randomIDCreation()
	fmt.Printf("DEBUG: Here is our randomID: %v\n", theID)
	giveJSON := returnJSON{
		SuccessMsg:     "Successful ID Given",
		FoodIDReturned: theID, //Go get unique IDS for 2 DBS
		SuccessBool:    true,
	}
	dataJSON, err := json.Marshal(giveJSON)
	if err != nil {
		fmt.Println("There's an error marshalling.")
		logWriter("There's an error marshalling.")
	}

	fmt.Fprintf(w, string(dataJSON))
}

//This should return foodIDS for a User or ALL Users for hotdogs
func getFoodIDSHDog(userID int) []int {
	var foodIDS []int

	hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //here's our collection
	filter := bson.D{{"foodid", userID}}                                           //Here's our filter to look for
	//Here's how to find and assign multiple Documents using a cursor
	// Pass these options to the Find method
	findOptions := options.Find()
	//Not needed, mostly for debugging: findOptions.SetLimit(2)
	cur, err := hotdogCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		logWriter("Issue finding food for this ID: " + string(userID) + " " + err.Error())
		log.Fatal(err)
	}
	//Loop through results
	for cur.Next(theContext) {
		// create a value into which the single document can be decoded
		var elem MongoHotDog
		err := cur.Decode(&elem)
		if err != nil {
			logWriter("Issue writing the current element for hotdogs: " + err.Error())
			log.Fatal(err)
		}
		foodIDS = append(foodIDS, elem.FoodID)
	}
	if err := cur.Err(); err != nil {
		logWriter("Issue looping through hotdogs: " + err.Error())
		log.Fatal(err)
	}

	return foodIDS
}

func getFoodIDSHam(userID int) []int {
	var foodIDS []int

	hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //here's our collection
	filter := bson.D{{"foodid", userID}}                                                 //Here's our filter to look for
	//Here's how to find and assign multiple Documents using a cursor
	// Pass these options to the Find method
	findOptions := options.Find()
	//Not needed, mostly for debugging: findOptions.SetLimit(2)
	cur, err := hamburgerCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		logWriter("Issue finding hamburger food for this ID: " + string(userID) + " " + err.Error())
		log.Fatal(err)
	}
	//Loop through results
	for cur.Next(theContext) {
		// create a value into which the single document can be decoded
		var elem MongoHamburger
		err := cur.Decode(&elem)
		if err != nil {
			logWriter("Issue writing the current element for hamburgers: " + err.Error())
			log.Fatal(err)
		}
		foodIDS = append(foodIDS, elem.FoodID)
	}
	if err := cur.Err(); err != nil {
		logWriter("Issue looping through hamburgers: " + err.Error())
		log.Fatal(err)
	}

	return foodIDS
}

//Gets all food when the mainpage is loaded
func getAllFoodMongo(w http.ResponseWriter, req *http.Request) {
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	//Marshal it into our type
	var theUser User
	json.Unmarshal(bs, &theUser)

	//Delcare our array of food returned
	var returnedHotDogs []MongoHotDog
	var returnedHamburgers []MongoHamburger
	var returnedHDogPics []UserPhoto
	var returnedHamPics []UserPhoto
	var returnedHDogLink []string
	var returnedHamLink []string
	//Declare data struct to send back
	type data struct {
		SuccessMessage string           `json:"SuccessMessage"`
		TheHotDogs     []MongoHotDog    `json:"TheHotDogs"`
		TheHamburgers  []MongoHamburger `json:"TheHamburgers:`
		TheHotDogsPics []UserPhoto      `json:"TheHotDogsPics"`
		TheHamPics     []UserPhoto      `json:"TheHamPics"`
		TheHDogLinks   []string         `json:"TheHDogLinks"`
		TheHamLinks    []string         `json:"TheHamLinks"`
		HaveHotDogs    bool             `json:"HaveHotDogs"`
		HaveHamburgers bool             `json:"HaveHamburgers"`
		HaveHDogPics   bool             `json:"HaveHDogPics"`
		HaveHamPics    bool             `json:"HaveHamPics"`
	}

	//Search for UserID or for all
	hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs")       //Here's our collection
	hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
	picCollection := mongoClient.Database("superdbtest1").Collection("userphotos")       //Here's our collection

	if theUser.UserID == 0 {
		//Query Mongo for all hotdogs
		theFilter := bson.M{}
		findOptions := options.Find()
		curHDog, err := hotdogCollection.Find(theContext, theFilter, findOptions)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				fmt.Printf("No documents were returned for hotdogs in MongoDB: %v\n", err.Error())
				logWriter("No documents were returned in MongoDB for Hotdogs: " + err.Error())
			} else {
				fmt.Printf("There was an error returning hotdogs for all Users, spot1: %v\n", err.Error())
				logWriter("There was an error returning hotdogs in Mongo for all Users, spot1: " + err.Error())
			}
		}
		//Loop over query results and fill hotdogs array
		for curHDog.Next(theContext) {
			// create a value into which the single document can be decoded
			var aHotDog MongoHotDog
			err := curHDog.Decode(&aHotDog)
			if err != nil {
				fmt.Printf("Error decoding hotdogs in MongoDB for all Users spot2: %v\n", err.Error())
				logWriter("Error decoding hotdogs in MongoDB for all Users spot2: " + err.Error())
			}
			returnedHotDogs = append(returnedHotDogs, aHotDog)
		}
		// Close the cursor once finished
		curHDog.Close(theContext)

		//Query Mongo for all Hamburgers
		curHam, err := hamburgerCollection.Find(theContext, theFilter, findOptions)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				fmt.Printf("No documents were returned for hamburgeres in MongoDB: %v\n", err.Error())
				logWriter("No documents were returned in MongoDB for hamburgers: " + err.Error())
			} else {
				fmt.Printf("There was an error returning hamburgers for all Users: %v\n", err.Error())
				logWriter("There was an error returning hamburgers in Mongo for all Users: " + err.Error())
			}
		}

		//Loop over query results and fill hotdogs array
		for curHam.Next(theContext) {
			// create a value into which the single document can be decoded
			var aHamburger MongoHamburger
			err := curHam.Decode(&aHamburger)
			if err != nil {
				fmt.Printf("Error decoding hamburgers in MongoDB for all Users: %v\n", err.Error())
				logWriter("Error decoding hamburgers in MongoDB for all Users: " + err.Error())
			}
			returnedHamburgers = append(returnedHamburgers, aHamburger)
			fmt.Printf("DEBUG: Here's the returnedHamburgs: %v\n", returnedHamburgers)
		}

		// Close the cursor once finished
		curHam.Close(theContext)

		//Query Mongo for ALL food pics for ALL Users
		curPics, err := picCollection.Find(theContext, theFilter, findOptions)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				fmt.Printf("No documents were returned for pictures in MongoDB: %v\n", err.Error())
				logWriter("No documents were returned in MongoDB for Pictures: " + err.Error())
			} else {
				fmt.Printf("There was an error returning pictures for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("There was an error returning pictures in Mongo for this User " + err.Error())
			}
		}

		//Loop over results and fill the pictues array
		for curPics.Next(theContext) {
			// create a value into which the single document can be decoded
			var aPic UserPhoto
			err := curPics.Decode(&aPic)
			if err != nil {
				fmt.Printf("Error decoding pictures in MongoDB for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("Error decoding pictures in MongoDB: " + err.Error())
			}
			//Determine if this is a hotdog or hamburger pic then add the appropriate links
			if strings.Contains(strings.ToUpper(aPic.FoodType), "HOTDOG") {
				fmt.Printf("DEBUG: Adding a hotdog to the returned pic list.\n")
				returnedHDogPics = append(returnedHDogPics, aPic)
				filePath := filepath.Join(aPic.Link)
				returnedHDogLink = append(returnedHDogLink, filePath)
			} else if strings.Contains(strings.ToUpper(aPic.FoodType), "HAMBURGER") {
				fmt.Printf("DEBUG: Adding a hamburger to the returned pic list.\n")
				returnedHamPics = append(returnedHamPics, aPic)
				filePath := filepath.Join(aPic.Link)
				returnedHamLink = append(returnedHamLink, filePath)
			} else {
				fmt.Printf("Error assigning picture to hamburger/hotdogs. FoodType is: %v\n", aPic.FoodType)
			}
		}

		curPics.Close(theContext) //Close the cursor when finished

		//Get the photo from Amazon if it's not already submitted
		checkPhoto(returnedHDogPics, returnedHamPics)
	} else {
		//Query Mongo for all hotdogs for a User
		theFilter := bson.M{"userid": theUser.UserID}
		findOptions := options.Find()
		curHDog, err := hotdogCollection.Find(theContext, theFilter, findOptions)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				fmt.Printf("No documents were returned for hotdogs in MongoDB: %v\n", err.Error())
				logWriter("No documents were returned in MongoDB for Hotdogs: " + err.Error())
			} else {
				fmt.Printf("There was an error returning hotdogs for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("There was an error returning hotdogs in Mongo for this User " + err.Error())
			}
		}
		//Loop over query results and fill hotdogs array
		for curHDog.Next(theContext) {
			// create a value into which the single document can be decoded
			var aHotDog MongoHotDog
			err := curHDog.Decode(&aHotDog)
			if err != nil {
				fmt.Printf("Error decoding hotdogs in MongoDB for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("Error decoding hotdogs in MongoDB: " + err.Error())
			}
			returnedHotDogs = append(returnedHotDogs, aHotDog)
		}
		// Close the cursor once finished
		curHDog.Close(theContext)

		//Query Mongo for all Hamburgers
		curHam, err := hamburgerCollection.Find(theContext, theFilter, findOptions)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				fmt.Printf("No documents were returned for hamburgers in MongoDB for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("No documents were returned in MongoDB for hamburgers: " + err.Error())
			} else {
				fmt.Printf("There was an error returning hamburgers for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("There was an error returning hamburgers in Mongo: " + err.Error())
				log.Fatal(err)
			}
		}

		//Loop over query results and fill hotdogs array
		for curHam.Next(theContext) {
			// create a value into which the single document can be decoded
			var aHamburger MongoHamburger
			err := curHam.Decode(&aHamburger)
			if err != nil {
				fmt.Printf("Error decoding hamburgers in MongoDB for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("Error decoding hamburgers in MongoDB: " + err.Error())
			}
			returnedHamburgers = append(returnedHamburgers, aHamburger)
		}

		// Close the cursor once finished
		curHam.Close(theContext)

		//Query Mongo for all food pictures for this User
		curPics, err := picCollection.Find(theContext, theFilter, findOptions)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				fmt.Printf("No documents were returned for pictures in MongoDB: %v\n", err.Error())
				logWriter("No documents were returned in MongoDB for Pictures: " + err.Error())
			} else {
				fmt.Printf("There was an error returning pictures for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("There was an error returning pictures in Mongo for this User " + err.Error())
			}
		}

		//Loop over results and fill the pictues array
		for curPics.Next(theContext) {
			// create a value into which the single document can be decoded
			var aPic UserPhoto
			err := curPics.Decode(&aPic)
			if err != nil {
				fmt.Printf("Error decoding pictures in MongoDB for this User, %v: %v\n", theUser.UserID, err.Error())
				logWriter("Error decoding pictures in MongoDB: " + err.Error())
			}
			//Determine if this is a hotdog or hamburger pic then add the appropriate links
			if strings.Contains(strings.ToUpper(aPic.FoodType), "HOTDOG") {
				fmt.Printf("DEBUG: Adding a hotdog to the returned pic list.\n")
				returnedHDogPics = append(returnedHDogPics, aPic)
				filePath := filepath.Join(aPic.Link)
				returnedHDogLink = append(returnedHDogLink, filePath)
			} else if strings.Contains(strings.ToUpper(aPic.FoodType), "HAMBURGER") {
				fmt.Printf("DEBUG: Adding a hamburger to the returned pic list.\n")
				returnedHamPics = append(returnedHamPics, aPic)
				filePath := filepath.Join(aPic.Link)
				returnedHamLink = append(returnedHamLink, filePath)
			} else {
				fmt.Printf("Error assigning picture to hamburger/hotdogs. FoodType is: %v\n", aPic.FoodType)
			}
		}

		curPics.Close(theContext) //Close the cursor when finished

		//Get the photo from Amazon if it's not already submitted
		checkPhoto(returnedHDogPics, returnedHamPics)
	}
	//Assemble data to return
	sendData := data{
		SuccessMessage: "Success",
		TheHotDogs:     returnedHotDogs,
		TheHamburgers:  returnedHamburgers,
		HaveHotDogs:    true,
		HaveHamburgers: true,
		HaveHDogPics:   true,
		HaveHamPics:    true,
	}

	//Do a wellness check for the data
	if len(sendData.TheHotDogs) <= 0 && len(sendData.TheHamburgers) <= 0 {
		sendData.SuccessMessage = "Failure"
	}
	if len(sendData.TheHotDogs) <= 0 {
		sendData.HaveHotDogs = false //Allow our loops to function properly in JS
	}
	if len(sendData.TheHamburgers) <= 0 {
		sendData.HaveHamburgers = false //Allow our loops to funciton properly in JS
	}
	if len(sendData.TheHDogLinks) <= 0 {
		sendData.HaveHDogPics = false //Allow our loops to funciton properly in JS
	}
	if len(sendData.TheHamLinks) <= 0 {
		sendData.HaveHamPics = false //Allow our loops to funciton properly in JS
	}
	//Marshal data to JSON
	dataJSON, err := json.Marshal(sendData)
	if err != nil {
		fmt.Println("There's an error marshalling.")
		logWriter("There's an error marshalling.")
	}

	fmt.Fprintf(w, string(dataJSON))
}

//This is for sorting the food into one string array for Mongo
func turnFoodArray(foodString string) []string {
	var returnedFood []string

	testArray := strings.Fields(foodString)

	returnedFood = testArray

	return returnedFood
}

//This is for turning the foodArray into a string
func sortFoodArray(foodArray []string) string {
	var returnString string
	for n := 0; n < len(foodArray); n++ {
		if n != len(foodArray)-1 {
			returnString = returnString + foodArray[n] + " "
		} else {
			returnString = returnString + foodArray[n]
		}
	}
	return returnString
}

/********* AMAZON QUERIES **********/

//Insert Photo into DB
func mongoInsertPhoto(userid int, foodid int, photoid int, photoName string, fileType string, size int64,
	photoHash string, link string, foodType string, dateCreated string, dateUpdated string) bool {
	successfulInsert := true
	fmt.Printf("DEBUG: Inserting photos into Mongo.\n")
	theTimeNow := time.Now()

	if strings.Contains(strings.ToUpper(foodType), "HOTDOG") {
		fmt.Printf("Inserting Hotdog Photo into MongoDB\n")
		photoInsertion := UserPhoto{
			UserID:      userid,
			FoodID:      foodid,
			PhotoID:     photoid,
			PhotoName:   photoName,
			FileType:    fileType,
			Size:        size,
			PhotoHash:   photoHash,
			Link:        link,
			FoodType:    foodType,
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Collect Data for Mongo
		photoCollection := mongoClient.Database("superdbtest1").Collection("userphotos") //Here's our collection
		collectedUsers := []interface{}{photoInsertion}
		//Insert Our Data
		insertManyResult, err := photoCollection.InsertMany(theContext, collectedUsers)
		if err != nil {
			fmt.Printf("Error inserting a hotdog photo in Mongo: %v\n", err.Error())
			successfulInsert = false
		} else {
			fmt.Printf("Insertion successful: %v\n", insertManyResult)
			/******** UPDATE HOTDOG AND USER COLLECTIONS *******/
			var aHotDogMongo MongoHotDog // create a value into which the single document can be decoded
			hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs")
			theFilter := bson.M{"foodid": foodid}
			findOptions := options.Find()
			curHDog, err := hotdogCollection.Find(theContext, theFilter, findOptions)
			if err != nil {
				if strings.Contains(err.Error(), "no documents in result") {
					fmt.Printf("No documents were returned for hotdogs in MongoDB: %v\n", err.Error())
					logWriter("No documents were returned in MongoDB for Hotdogs: " + err.Error())
				} else {
					fmt.Printf("There was an error returning hotdogs for this Hotdog, %v: %v\n",
						foodid, err.Error())
					logWriter("There was an error returning hotdogs in Mongo for this User " + err.Error())
				}
			}
			//Loop over query results and fill hotdogs array
			for curHDog.Next(theContext) {
				err := curHDog.Decode(&aHotDogMongo)
				if err != nil {
					fmt.Printf("Error decoding hotdogs in MongoDB for this User, %v: %v\n",
						foodid, err.Error())
					logWriter("Error decoding hotdogs in MongoDB: " + err.Error())
				}
			}
			// Close the cursor once finished
			curHDog.Close(theContext)
			//Update Dogs
			type foodUpdate struct {
				FoodType     string    `json:"FoodType"`
				FoodID       int       `json:"FoodID"`
				TheHamburger Hamburger `json:"TheHamburger"`
				TheHotDog    Hotdog    `json:"TheHotDog"`
			}
			newLink := filepath.Join("amazonimages", link)
			fixedLink := urlFixer(newLink)
			aHotDog := Hotdog{
				HotDogType:  aHotDogMongo.HotDogType,
				Condiment:   sortFoodArray(aHotDogMongo.Condiments),
				Calories:    aHotDogMongo.Calories,
				Name:        aHotDogMongo.Name,
				UserID:      aHotDogMongo.UserID,
				FoodID:      aHotDogMongo.FoodID,
				PhotoID:     photoid,
				PhotoSrc:    fixedLink,
				DateCreated: aHotDogMongo.DateCreated,
				DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
			}
			hDogUpdate := foodUpdate{
				FoodType:     "hotdog",
				FoodID:       foodid,
				TheHamburger: Hamburger{},
				TheHotDog:    aHotDog,
			}
			jsonValue, _ := json.Marshal(hDogUpdate)
			response, err := http.Post("http://"+serverAddress+"/foodUpdateMongo",
				"application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				successfulInsert = false
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(data))
			}
		}
	} else if strings.Contains(strings.ToUpper(foodType), "HAMBURGER") {
		fmt.Printf("Inserting Hamburger Photo into MongoDB\n")
		photoInsertion := UserPhoto{
			UserID:      userid,
			FoodID:      foodid,
			PhotoID:     photoid,
			PhotoName:   photoName,
			FileType:    fileType,
			Size:        size,
			PhotoHash:   photoHash,
			Link:        link,
			FoodType:    foodType,
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Collect Data for Mongo
		photoCollection := mongoClient.Database("superdbtest1").Collection("userphotos") //Here's our collection
		collectedUsers := []interface{}{photoInsertion}
		//Insert Our Data
		insertManyResult, err := photoCollection.InsertMany(theContext, collectedUsers)
		if err != nil {
			fmt.Printf("Error inserting a hamburger photo in Mongo: %v\n", err.Error())
			successfulInsert = false
		} else {
			fmt.Printf("Insertion successful: %v\n", insertManyResult)
			/******** UPDATE HAMBURGER AND USER COLLECTIONS *******/
			var aHamMongo MongoHamburger // create a value into which the single document can be decoded
			hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers")
			theFilter := bson.M{"foodid": foodid}
			findOptions := options.Find()
			curHam, err := hamburgerCollection.Find(theContext, theFilter, findOptions)
			if err != nil {
				if strings.Contains(err.Error(), "no documents in result") {
					fmt.Printf("No documents were returned for hamburgers in MongoDB: %v\n", err.Error())
					logWriter("No documents were returned in MongoDB for hamburgers: " + err.Error())
				} else {
					fmt.Printf("There was an error returning hamburgers for this Hamburger, %v: %v\n",
						foodid, err.Error())
					logWriter("There was an error returning hamburgers in Mongo for this User " + err.Error())
				}
			}
			//Loop over query results and fill hotdogs array
			for curHam.Next(theContext) {
				err := curHam.Decode(&aHamMongo)
				if err != nil {
					fmt.Printf("Error decoding hamburgers in MongoDB for this User, %v: %v\n",
						foodid, err.Error())
					logWriter("Error decoding hamburgers in MongoDB: " + err.Error())
				}
			}
			// Close the cursor once finished
			curHam.Close(theContext)
			//Update Hamburgers
			type foodUpdate struct {
				FoodType     string    `json:"FoodType"`
				FoodID       int       `json:"FoodID"`
				TheHamburger Hamburger `json:"TheHamburger"`
				TheHotDog    Hotdog    `json:"TheHotDog"`
			}
			newLink := filepath.Join("amazonimages", link)
			fixedLink := urlFixer(newLink)
			aHam := Hamburger{
				BurgerType:  aHamMongo.BurgerType,
				Condiment:   sortFoodArray(aHamMongo.Condiments),
				Calories:    aHamMongo.Calories,
				Name:        aHamMongo.Name,
				UserID:      aHamMongo.UserID,
				FoodID:      aHamMongo.FoodID,
				PhotoID:     photoid,
				PhotoSrc:    fixedLink,
				DateCreated: aHamMongo.DateCreated,
				DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
			}
			hamUpdate := foodUpdate{
				FoodType:     "hamburger",
				FoodID:       foodid,
				TheHamburger: aHam,
				TheHotDog:    Hotdog{},
			}
			jsonValue, _ := json.Marshal(hamUpdate)
			response, err := http.Post("http://"+serverAddress+"/foodUpdateMongo",
				"application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				successfulInsert = false
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(data))
			}
		}
	} else {
		fmt.Printf("Wrong food type returned for food photo data insertion: %v\n", foodType)
		successfulInsert = false
	}

	return successfulInsert
}

func mongoPhotoDeletion(foodID int) {
	fmt.Printf("Deleting the photo for this ID: %v\n", foodID)
	/* FIRST DELETE FROM HOTDOG COLLECTION*/
	picCollection := mongoClient.Database("superdbtest1").Collection("userphotos") //Here's our collection
	deletes := []bson.M{
		{"foodid": foodID},
	} //Here's our filter to look for
	deletes = append(deletes, bson.M{"foodid": bson.M{
		"$eq": foodID,
	}},
	)

	// create the slice of write models
	var writes []mongo.WriteModel

	for _, del := range deletes {
		model := mongo.NewDeleteManyModel().SetFilter(del)
		writes = append(writes, model)
	}

	// run bulk write
	res, err := picCollection.BulkWrite(theContext, writes)
	if err != nil {
		logWriter("Error writing Mongo Delete Statement for photo deletion")
		logWriter("\n")
		logWriter(err.Error())
	}
	//Print Results
	fmt.Printf("Deleted the following documents: %v\n", res.DeletedCount)
	logWriter("Deleted the following documents for food photos: " + string(res.DeletedCount) + "\n")
}
