package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

var theContext context.Context

func connectDB() *mongo.Client {
	//Setup Mongo connection to Atlas Cluster
	theClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://joek:superduperPWord@superdbcluster.kswud.mongodb.net/superdbtest1?retryWrites=true&w=majority"))
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
