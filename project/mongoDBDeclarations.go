package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

var theContext context.Context

func connectDB() *mongo.Client {
	//Setup Mongo connection to Atlas Cluster
	theClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://bigjohnny:figleafs@superdbcluster-shard-00-00.kswud.mongodb.net:27017,superdbcluster-shard-00-01.kswud.mongodb.net:27017,superdbcluster-shard-00-02.kswud.mongodb.net:27017/superdbtest1?ssl=true&replicaSet=atlas-pvjlol-shard-0&authSource=admin&retryWrites=true&w=majority"))
	if err != nil {
		fmt.Printf("Errored getting mongo client: %v\n", err)
		log.Fatal(err)
	}
	theContext, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
	databases, err := theClient.ListDatabaseNames(theContext, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

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
		Condiments:  []string{postedHotDog.Condiment},
		Calories:    postedHotDog.Calories,
		Name:        postedHotDog.Name,
		FoodID:      randomIDCreation(),
		UserID:      postedHotDog.UserID,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}

	//Collect Data for Mongo
	user_collection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
	collectedUsers := []interface{}{mongoHotDogInsert}
	//Insert Our Data
	insertManyResult, err := user_collection.InsertMany(context.TODO(), collectedUsers)
	if err != nil {
		fmt.Printf("Error inserting results: \n%v\n", err)
		fmt.Fprint(w, failureMessage)
		log.Fatal(err)
	} else {
		fmt.Fprint(w, successMessage)
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
		Condiments:  []string{postedHamburger.Condiment},
		Calories:    postedHamburger.Calories,
		Name:        postedHamburger.Name,
		FoodID:      randomIDCreation(),
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
		fmt.Printf("DEBUG: Updating hotdog at id: %v\n", thefoodUpdate.FoodID)
		theTimeNow := time.Now()
		var hotDogUpdate Hotdog = thefoodUpdate.TheHotDog
		updatedHotDogMongo := MongoHotDog{
			HotDogType:  hotDogUpdate.HotDogType,
			Condiments:  []string{hotDogUpdate.Condiment},
			Calories:    hotDogUpdate.Calories,
			Name:        hotDogUpdate.Name,
			UserID:      hotDogUpdate.UserID,
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Add updatedHotDog to Document collection for Hotdogs
		ic_collection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
		filter := bson.D{{"UserID", updatedHotDogMongo.UserID}}                     //Here's our filter to look for
		update := bson.D{                                                           //Here is our data to update
			{"$set", bson.D{
				{"HotDogType", updatedHotDogMongo.HotDogType},
				{"Condiments", updatedHotDogMongo.Condiments},
				{"Calories", updatedHotDogMongo.Calories},
				{"Name", updatedHotDogMongo.Name},
				{"DateUpdated", updatedHotDogMongo.DateUpdated},
			}},
		}

		updateResult, err := ic_collection.UpdateMany(context.TODO(), filter, update)
		if err != nil {
			fmt.Fprintln(w, 3) //Failure Response Response
			log.Fatal(err)
		} else {
			//Our new UpdateResult
			fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
			fmt.Fprintln(w, 1) //Success Response
		}
	} else if thefoodUpdate.FoodType == "hamburger" {
		fmt.Printf("DEBUG: Updating Hamburger at id: %v\n", thefoodUpdate.FoodID)
		theTimeNow := time.Now()
		var hamburgerUpdate Hamburger = thefoodUpdate.TheHamburger
		updatedHamburgerMongo := MongoHamburger{
			BurgerType:  hamburgerUpdate.BurgerType,
			Condiments:  []string{hamburgerUpdate.Condiment},
			Calories:    hamburgerUpdate.Calories,
			Name:        hamburgerUpdate.Name,
			UserID:      hamburgerUpdate.UserID,
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Add updatedHotDog to Document collection for Hotdogs
		ic_collection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
		filter := bson.D{{"UserID", updatedHamburgerMongo.UserID}}                     //Here's our filter to look for
		update := bson.D{                                                              //Here is our data to update
			{"$set", bson.D{
				{"BurgerType", updatedHamburgerMongo.BurgerType},
				{"Condiments", updatedHamburgerMongo.Condiments},
				{"Calories", updatedHamburgerMongo.Calories},
				{"Name", updatedHamburgerMongo.Name},
				{"DateUpdated", updatedHamburgerMongo.DateUpdated},
			}},
		}

		updateResult, err := ic_collection.UpdateMany(context.TODO(), filter, update)
		if err != nil {
			fmt.Fprintln(w, 3) //Failure Response Response
			log.Fatal(err)
		} else {
			//Our new UpdateResult
			fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
			fmt.Fprintln(w, 1) //Success Response
		}
	} else {
		fmt.Fprintln(w, 3)
	}
}

//DEBUG: Work in Progress
func foodDeleteMongo(w http.ResponseWriter, req *http.Request) {
	//Our Food deletion struct
	type foodDeletion struct {
		FoodType     string    `json:"FoodType"`
		TheHamburger Hamburger `json:"TheHamburger"`
		TheHotDog    Hotdog    `json:"TheHotDog"`
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
		/* FIRST DELETE FROM HOTDOG COLLECTION*/
		hotdogCollection := mongoClient.Database("superdbtest1").Collection("hotdogs") //Here's our collection
		deletes := []bson.M{
			{"UserID": theFoodDeletion.TheHotDog.UserID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"HotDogType": bson.M{
			"$eq": foodSlurs[j],
		}}, bson.M{"Condiments": bson.M{
			"$eq": foodSlurs[j],
		}}, bson.M{"Name": bson.M{
			"$eq": foodSlurs[j],
		}})

		// create the slice of write models
		var writes []mongo.WriteModel

		for _, del := range deletes {
			model := mongo.NewDeleteManyModel().SetFilter(del)
			writes = append(writes, model)
		}

		// run bulk write
		res, err := hotdogCollection.BulkWrite(context.TODO(), writes)
		if err != nil {
			logWriter("Error writing Mongo Delete Statement")
			logWriter("\n")
			logWriter(err.Error())
			log.Fatal(err)
		}
		//Print Results
		fmt.Printf("Deleted the following documents: %v\n", res.DeletedCount)
		logWriter("Deleted the following documents: " + string(res.DeletedCount) + "\n")

		/* NOW DELETE FROM USER COLLECITON */

		fmt.Fprintln(w, 1)
	} else if theFoodDeletion.FoodType == "hamburger" {
		hamburgerCollection := mongoClient.Database("superdbtest1").Collection("hamburgers") //Here's our collection
		deletes := []bson.M{
			{"UserID": theFoodDeletion.TheHotDog.UserID},
		} //Here's our filter to look for

		// create the slice of write models
		var writes []mongo.WriteModel

		for _, del := range deletes {
			model := mongo.NewDeleteManyModel().SetFilter(del)
			writes = append(writes, model)
		}

		// run bulk write
		res, err := hamburgerCollection.BulkWrite(context.TODO(), writes)
		if err != nil {
			logWriter("Error writing Mongo Delete Statement")
			logWriter("\n")
			logWriter(err.Error())
			log.Fatal(err)
		}
		//Print Results
		fmt.Printf("Deleted the following documents: %v\n", res.DeletedCount)
		logWriter("Deleted the following documents: " + string(res.DeletedCount) + "\n")

		fmt.Fprintln(w, 2)
	} else {
		fmt.Fprintln(w, 3)
	}
}

//This should give a random id value to both food groups
func randomIDCreation() int {
	fmt.Printf("DEBUG: Creating Random ID for User/Food\n")
	finalID := 0        //The final, unique ID to return to the food/user
	randInt := 0        //The random integer added onto ID
	randIntString := "" //The integer built through a string...
	min, max := 0, 9    //The min and Max value for our randInt
	foundID := false
	for foundID == false {
		for i := 0; i < 8; i++ {
			randInt = rand.Intn(max-min) + min
			randIntString = randIntString + strconv.Itoa(randInt)
		}
		//Once we have a string of numbers, we can convert it back to an integer
		theID, err := strconv.Atoi(randIntString)
		if err != nil {
			fmt.Printf("We got an error converting a string back to a number, %v\n", err)
			fmt.Println(err)
		}
		//Search all our collections to see if this UserID is unique
		canExit := true
		//User collection
		userCollection := mongoClient.Database("superdbtest1").Collection("users") //Here's our collection
		var testAUser AUser
		theErr := userCollection.FindOne(context.TODO(), bson.M{"userid": theID}).Decode(&testAUser)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User and our ID is clean.\n")
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit = false
				log.Fatal(theErr)
			}
		}
		if testAUser.UserID == theID {
			canExit = false
		} else {
			canExit = true
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
		theErr = hotdogCollection.FindOne(context.TODO(), theFilter).Decode(&testHotdog)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User and our ID is clean.\n")
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit = false
				log.Fatal(theErr)
			}
		}
		//Check to see if the ID was found for a hotdog database
		if testAUser.UserID == theID {
			canExit = false
		} else {
			canExit = true
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
		theErr = hamburgerCollection.FindOne(context.TODO(), theFilter2).Decode(&testBurger)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				fmt.Printf("It's all good, this document wasn't found for User and our ID is clean.\n")
			} else {
				fmt.Printf("DEBUG: We have another error for finding a unique UserID: \n%v\n", theErr)
				canExit = false
				log.Fatal(theErr)
			}
		}
		//Check to see if the ID was found for a hotdog database
		if testAUser.UserID == theID {
			canExit = false
		} else {
			canExit = true
		}
		//Final check to see if we can exit this loop
		if canExit == true {
			finalID = theID
			foundID = true
		} else {
			foundID = false
		}
	}

	return finalID
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
