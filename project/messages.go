package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//Data sent to index page
type MessageViewData struct {
	TestString  string    `json:"TestString"`
	TheMessages []Message `json:"TheMessages"`
	WhatPage    int       `json:"WhatPage"`
}

//Message displayed on the board
type Message struct {
	MessageID       int       `json:"MessageID"`       //ID of this Message
	UserID          int       `json:"UserID"`          //ID of the owner of this message
	Messages        []Message `json:"Messages"`        //Array of Messages under this one
	IsChild         bool      `json:"IsChild"`         //Is this message childed to another message
	HasChildren     bool      `json:"HasChildren"`     //Whether this message has children to list
	ParentMessageID int       `json:"ParentMessageID"` //The ID of this parent
	UberParentID    int       `json:"UberParentID"`    //The final parent of this parent, IF EQUAL PARENT
	Order           int       `json:"Order"`           //Order the commnet is in with it's reply tree
	RepliesAmount   int       `json:"RepliesAmount"`   //Amount of replies this message has
	TheMessage      string    `json:"TheMessage"`      //The MEssage in the post
	DateCreated     string    `json:"DateCreated"`     //When the message was created
	LastUpdated     string    `json:"LastUpdated"`     //When the message was last updated
}

//All the Messages on the board
type MessageBoard struct {
	MessageBoardID         int             `json:"MessageBoardID"`
	AllMessages            []Message       `json:"AllMessages"`            //All the IDs listed
	AllMessagesMap         map[int]Message `json:"AllMessagesMap"`         //A map of ALL messages
	AllOriginalMessages    []Message       `json:"AllOriginalMessages"`    //All the messages that AREN'T replies
	AllOriginalMessagesMap map[int]Message `json:"AllOriginalMessagesMap"` //Map of original Messages
	LastUpdated            string          `json:"LastUpdated"`            //Last time this messageboard was updated
}

var loadedMessagesMap map[int]Message
var theMessageBoard MessageBoard //The board containing all our messages
/* This is the current amount of results our User is looking at
it changes as the User clicks forwards or backwards for more results */
var currentPageNumber int = 1

//Creates a list of test messages
func createTestMessages() {
	createdTestBoard := isMessageBoardCreated()
	if createdTestBoard == true {
		fmt.Printf("DEBUG: Yo, we go that messageboard already\n")
	} else {
		theTimeNow := time.Now()
		//Fill test map values
		fillerBoard := make(map[int]Message)
		//Make test message board to work with
		testMessageBoard := MessageBoard{
			MessageBoardID:         5555,
			AllMessages:            []Message{},
			AllMessagesMap:         fillerBoard,
			AllOriginalMessages:    []Message{},
			AllOriginalMessagesMap: fillerBoard,
			LastUpdated:            theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Assign testMessageboard to actual message board
		theMessageBoard = testMessageBoard
		//Collect Data for Mongo
		message_collection := mongoClient.Database("messageboard").Collection("messageboard") //Here's our collection
		collectedStuff := []interface{}{theMessageBoard}                                      //Send the 'UberMessage' with updated parent info
		//Insert Our main data to the 'messageboard' table, since this is a parent
		insertManyResult, err := message_collection.InsertMany(context.TODO(), collectedStuff)
		if err != nil {
			fmt.Printf("Error inserting results: \n%v\n", err)
			log.Fatal(err)
		} else {
			fmt.Println("Inserted multiple documents in createTestMessages: ",
				insertManyResult.InsertedIDs) //Data insert results
		}

		fmt.Printf("DEBUG: Test MessageBoard Created.\n")
	}
}

//Returns true if our test message board is already created
func isMessageBoardCreated() bool {
	messageCollection := mongoClient.Database("messageboard").Collection("messageboard") //Here's our collection
	//Query Mongo for all Messages
	theFilter := bson.M{
		"messageboardid": bson.M{
			"$eq": 5555, // check if bool field has value of 'false'
		},
	}
	findOptions := options.Find()
	messageBoard, err := messageCollection.Find(theContext, theFilter, findOptions)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			fmt.Printf("No documents were returned for hotdogs in MongoDB: %v\n", err.Error())
		} else {
			fmt.Printf("There was an error returning hotdogs for all Users, spot1: %v\n", err.Error())
		}
	}
	// create a value into which the single document can be decoded
	var aMessageBoard MessageBoard
	//Loop over query results and fill hotdogs array
	for messageBoard.Next(theContext) {
		err := messageBoard.Decode(&aMessageBoard)
		if err != nil {
			fmt.Printf("Error decoding messageboard in MongoDB for all Users: %v\n", err.Error())
		}
		//Assign our message board to the 'theMessageBoard' to work with
		theMessageBoard = aMessageBoard
	}
	// Close the cursor once finished
	messageBoard.Close(theContext)

	if aMessageBoard.MessageBoardID != 5555 {
		return false
	} else {
		return true
	}
}

//Inserts one message into our 'messages' collection
func insertOneNewMessage(newMessage Message) {
	//Send this to the 'message' collection for safekeeping
	messageCollection := mongoClient.Database("messageboard").Collection("messages") //Here's our collection
	collectedStuff := []interface{}{newMessage}
	//Insert Our Data
	insertManyResult, err := messageCollection.InsertMany(context.TODO(), collectedStuff)
	if err != nil {
		fmt.Printf("Error inserting results: \n%v\n", err)
		log.Fatal(err)
	} else {
		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs) //Data insert results
	}
}

//Inserts original message into ALL dbs and our map
func uberInsertNewMessage(newMessage Message) {
	theTimeNow := time.Now()
	//Insert into 'message' table
	insertOneNewMessage(newMessage)
	//Update our server map
	loadedMessagesMap[len(loadedMessagesMap)+1] = newMessage
	//Update our other server map
	theMessageBoard.AllMessages = append(theMessageBoard.AllMessages, newMessage)
	theMessageBoard.AllMessagesMap[newMessage.MessageID] = newMessage
	theMessageBoard.AllOriginalMessages = append(theMessageBoard.AllOriginalMessages, newMessage)
	theMessageBoard.AllOriginalMessagesMap[newMessage.MessageID] = newMessage
	theMessageBoard.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
	//Update Mongo with new map
	updateMongoMessageBoard(theMessageBoard)
}

//Simulates updating a message from the parent downward
func uberUpdate(newestMessage Message, parentMessage Message) {
	theTimeNow := time.Now() //Used for time updates
	if parentMessage.IsChild == false {
		fmt.Printf("DEBUG: We are updating an UberParent\n")
		//This is the uberParent; simply add this to the []Message list
		parentMessage.Messages = append(parentMessage.Messages, newestMessage)
		//parentMessage.Messages = append([]Message{newestMessage}, parentMessage.Messages...)
		parentMessage.RepliesAmount = parentMessage.RepliesAmount + 1
		parentMessage.HasChildren = true
		parentMessage.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		fmt.Printf("In UberUpdate, here is the UberParent: %v\n", parentMessage)
		//Update Message in 'messages' table
		updateMessage(parentMessage)
		//Update the messages tables
		insertOneNewMessage(newestMessage)
		//Update MessageBoard properties
		theMessageBoard.AllMessagesMap[parentMessage.MessageID] = parentMessage          //Update UberParent
		theMessageBoard.AllOriginalMessagesMap[parentMessage.MessageID] = parentMessage  //updateUberParent
		theMessageBoard.AllMessages = append(theMessageBoard.AllMessages, newestMessage) //Add newest message
		theMessageBoard.AllMessagesMap[newestMessage.MessageID] = newestMessage          //add newest message
		theMessageBoard.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//update uberParent
		for j := 0; j < len(theMessageBoard.AllOriginalMessages); j++ {
			if theMessageBoard.AllOriginalMessages[j].MessageID == parentMessage.MessageID {
				theMessageBoard.AllOriginalMessages[j] = parentMessage
				//Update the loadedMessageMap
				loadedMessagesMap[j+1] = parentMessage
				break
			}
		}
		//Update Mongo Collections
		updateMongoMessageBoard(theMessageBoard)
	} else {
		fmt.Printf("DEBUG: Updating a NON-Uber Parent Message with a reply\n")
		//Add newest message to parent message to update it
		parentMessage.HasChildren = true //Had it or before, now this has children
		parentMessage.RepliesAmount = parentMessage.RepliesAmount + 1
		parentMessage.Messages = append(parentMessage.Messages, newestMessage)
		//parentMessage.Messages = append([]Message{newestMessage}, parentMessage.Messages...)
		parentMessage.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//Update Message in 'messages' table
		updateMessage(parentMessage)
		//Insert new Message
		insertOneNewMessage(newestMessage)
		fmt.Println()
		fmt.Printf("DEBUG: Here is the parent before updating uber: %v\n", parentMessage)
		fmt.Println()
		fmt.Println()
		fmt.Printf("DEBUG: Here is the child message: %v\n", newestMessage)
		fmt.Println()
		fmt.Println()
		//Search our message board for the Uber Parent
		uberParentMessage := theMessageBoard.AllOriginalMessagesMap[parentMessage.UberParentID]
		//Update all parent messages recursivley until we finally update the uberParentMessage
		uberParentMessage = updateToUber(uberParentMessage, parentMessage)
		fmt.Printf("DEBUG: Here is our UberParent Update: %v\n", uberParentMessage)
		fmt.Println()
		fmt.Println()
		//Update MessageBoard properties
		theMessageBoard.AllMessagesMap[uberParentMessage.MessageID] = uberParentMessage                //Update UberParent
		theMessageBoard.AllOriginalMessagesMap[uberParentMessage.MessageID] = uberParentMessage        //updateUberParent
		theMessageBoard.AllMessages = append([]Message{newestMessage}, theMessageBoard.AllMessages...) //Add newest message
		theMessageBoard.AllMessagesMap[newestMessage.MessageID] = newestMessage                        //add newest message
		theMessageBoard.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//update uberParent
		for j := 0; j < len(theMessageBoard.AllOriginalMessages); j++ {
			if theMessageBoard.AllOriginalMessages[j].MessageID == uberParentMessage.MessageID {
				theMessageBoard.AllOriginalMessages[j] = uberParentMessage
				//Update the loadedMessageMap
				loadedMessagesMap[j+1] = uberParentMessage
			}
		}

		//Update Mongo Collections
		updateMongoMessageBoard(theMessageBoard)
		//Update Message in 'messages' table
		updateMessage(uberParentMessage)
	}
}

//This func updates the Mongo MessageBoard collection to what we have now
func updateMongoMessageBoard(updatedMessageBoard MessageBoard) {
	message_collection := mongoClient.Database("messageboard").Collection("messageboard") //Here's our collection
	theFilter := bson.M{
		"messageboardid": bson.M{
			"$eq": updatedMessageBoard.MessageBoardID, // check if test value is present for reply Message
		},
	}

	updatedDocument := bson.M{}
	updatedDocument = bson.M{
		"$set": bson.M{
			"messageboardid":         updatedMessageBoard.MessageBoardID,
			"allmessages":            updatedMessageBoard.AllMessages,
			"allmessagesmap":         updatedMessageBoard.AllMessagesMap,
			"alloriginalmessages":    updatedMessageBoard.AllOriginalMessages,
			"alloriginalmessagesmap": updatedMessageBoard.AllOriginalMessagesMap,
			"lastupdated":            updatedMessageBoard.LastUpdated,
		},
	}

	stuffUpdated, err2 := message_collection.UpdateOne(theContext, theFilter, updatedDocument)
	if err2 != nil {
		fmt.Printf("We got an error updating this document: %v\n", err2.Error())
	} else {
		fmt.Printf("Here is the update for our messageboard: %v\n", stuffUpdated.MatchedCount)
	}
}

//This func finds the parent FROM THE UBERPARENT to update
func updateToUber(uberParentMessage Message, parentMessage Message) Message {
	theTimeNow := time.Now()      //Used for updating time properties in our parent
	uberParentUpdated := false    //Are we currently on the UberParent Message, matching their ID?
	finalUberMessage := Message{} //The final message with the updated UberParent
	//Loop and update until we find the UberParent and update it into the 'finalUberMessage'
	pastParent := parentMessage
	currentMessage := theMessageBoard.AllMessagesMap[parentMessage.ParentMessageID] //First set the searcher parent to it's OWN parent
	for {
		if uberParentUpdated == true {
			break //UberParent is found and updated, ready to be returned. End this updating search
		} else {
			fmt.Printf("DEBUG: We are currently on this parent: %v\n", currentMessage.TheMessage)
			/* Determine if the current parent is an uberParent */
			if currentMessage.MessageID == uberParentMessage.MessageID {
				//This is the parent update UberParent for us to return and update in Mongo, then break!
				uberParentUpdated = true //Set break value
				/* Step 1: Update the parentMessage in the UberParent */
				for v := 0; v < len(currentMessage.Messages); v++ {
					//Search fo rparent in UberParent Messages then update it
					if currentMessage.Messages[v].MessageID == pastParent.MessageID {
						currentMessage.Messages[v] = pastParent
						break
					}
				}
				/* Step 2: Update the Past Parent in the the messageboard table */
				pastParent.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				theMessageBoard.AllMessagesMap[pastParent.MessageID] = pastParent
				/* DEBUG: We should add a goroutine to update other sections of the MessageBoard
				that ARENT' maps and quickly updated */
				/* Step 2: Update finalUberMessage to return for final updating */
				finalUberMessage = currentMessage
			} else {
				/*
					This is not the UberParent. We will update the currentMessage in all appropriate spots,
					THEN we will assign the currentMessage as the ParentID, then move this currentMessage to
					be the 'pastParent'
				*/
				fmt.Printf("DEBUG: We are currently updating for this message: %v\n", currentMessage.TheMessage)
				//Step 1: Update the currentMessage's Message Array with the updated parentMessage
				for x := 0; x < len(currentMessage.Messages); x++ {
					if currentMessage.Messages[x].MessageID == pastParent.MessageID {
						//We found the pastParent in the currentMessage Messages array. Update it
						pastParent.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
						currentMessage.Messages[x] = pastParent
						break
					}
				}
				//Step 2: Update the messageBoard so we won't have infinite loops
				currentMessage.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				pastParent.LastUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				theMessageBoard.AllMessagesMap[currentMessage.MessageID] = currentMessage
				theMessageBoard.AllMessagesMap[pastParent.MessageID] = pastParent
				/* DEBUG: We should add a goroutine to update other sections of the MessageBoard
				that ARENT' maps and quickly updated */
				/* Step 3: Update the search criteria until we hit the uberParent */
				pastParent = currentMessage
				currentMessage = theMessageBoard.AllMessagesMap[currentMessage.ParentMessageID] //First set the searcher parent to it's OWN parent
			}
		}
	}

	return finalUberMessage //This should be the completed UberMessage
}

//This simply updates on message from a message given in the 'messages' collection
func updateMessage(updatedMessage Message) {
	message_collection := mongoClient.Database("messageboard").Collection("messages") //Here's our collection
	theFilter := bson.M{
		"messageid": bson.M{
			"$eq": updatedMessage.MessageID, // check if test value is present for reply Message
		},
	}

	updatedDocument := bson.M{}
	updatedDocument = bson.M{
		"$set": bson.M{
			"messageid":       updatedMessage.MessageID,
			"userid":          updatedMessage.UserID,
			"messages":        updatedMessage.Messages,
			"ischild":         updatedMessage.IsChild,
			"haschildren":     updatedMessage.HasChildren,
			"parentmessageid": updatedMessage.ParentMessageID,
			"uberparentid":    updatedMessage.UberParentID,
			"order":           updatedMessage.Order,
			"repliesamount":   updatedMessage.RepliesAmount,
			"themessage":      updatedMessage.TheMessage,
			"datecreated":     updatedMessage.DateCreated,
			"lastupdated":     updatedMessage.LastUpdated,
		},
	}

	stuffUpdated, err2 := message_collection.UpdateOne(theContext, theFilter, updatedDocument)
	if err2 != nil {
		fmt.Printf("We got an error updating this document: %v\n", err2.Error())
	} else {
		fmt.Printf("Here is the update for our single message: %v\n", stuffUpdated.ModifiedCount)
	}
}

/* This queries all the test messages that are NOT REPLIES...those should be entered as documents
already from the messages given */
func getAllTestMessages() {
	/* Clear map so when the page refreshes we have the correct values */
	loadedMessagesMap = make(map[int]Message)
	messageCollection := mongoClient.Database("messageboard").Collection("messageboard") //Here's our collection
	//Query Mongo for all Messages
	theFilter := bson.M{}
	findOptions := options.Find()
	messageBoard, err := messageCollection.Find(theContext, theFilter, findOptions)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			fmt.Printf("No documents were returned for hotdogs in MongoDB: %v\n", err.Error())
		} else {
			fmt.Printf("There was an error returning hotdogs for all Users, spot1: %v\n", err.Error())
		}
	}
	// create a value into which the single document can be decoded
	var aMessageBoard MessageBoard
	//Loop over query results and fill hotdogs array
	for messageBoard.Next(theContext) {
		err := messageBoard.Decode(&aMessageBoard)
		if err != nil {
			fmt.Printf("Error decoding hotdogs in MongoDB for all Users spot2: %v\n", err.Error())
		}
		//Assign our message board to the 'theMessageBoard' to work with
		theMessageBoard = aMessageBoard
	}
	// Close the cursor once finished
	messageBoard.Close(theContext)
	//Fill the 'loadedMessageMap' with map values for our page
	for x := 0; x < len(theMessageBoard.AllOriginalMessages); x++ {
		loadedMessagesMap[x+1] = theMessageBoard.AllOriginalMessages[x]
	}
}

func getTenResults(whatPageNum int) ([]Message, bool) {
	giveMessages := []Message{}
	topResult := whatPageNum * 10              //Last Result to add to map
	minResult := ((whatPageNum * 10) - 10) + 1 //First result to add to map
	okayResult := true                         //The result returned if we have messages to return

	//Initial check to see if a map exists
	if _, ok := loadedMessagesMap[minResult]; ok {
		//top value exists, get the message in range
		for x := minResult; x <= topResult; x++ {
			//Check to see if top result exists to value of ten; if not, add nothing
			if _, ok := loadedMessagesMap[x]; ok {
				giveMessages = append(giveMessages, loadedMessagesMap[x])
			} else {
				//Do nothing, there's no message here
			}
		}
		okayResult = true
	} else {
		fmt.Printf("DEBUG: Page value does not exist! The Value: %v\n", minResult)
		fmt.Printf("DEBUG: Here is our map currently: \n\n%v\n\n", loadedMessagesMap)
		okayResult = false
	}

	//Reversing order of slice for 'MessageBoard Display' purposes
	giveMessages = reverseSlice(giveMessages)

	for q := 0; q < len(giveMessages); q++ {
		//fmt.Printf("giveMessages results %v: %v\n", q, giveMessages[q])
	}

	return giveMessages, okayResult
}

//This is for reversing the order of a Message array for display
func reverseSlice(orderedSlice []Message) []Message {
	last := len(orderedSlice) - 1
	for i := 0; i < len(orderedSlice)/2; i++ {
		orderedSlice[i], orderedSlice[last-i] = orderedSlice[last-i], orderedSlice[i]
	}

	return orderedSlice
}

/* Called in Ajax from Javascript everytime User clicks left or right or submits a page
with results they'd like to see. If it's successful, it returns a number of JSON formatted Messages
for the page to update with. If not, it returns an error, which can be put in the pageNumber field. */
func evaluateTenResults(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: We submitted a ajax form, now in evaluateTenResults \n")
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Declare datatype from Ajax
	type PageData struct {
		ThePage int `json:"ThePage"`
	}
	//Unmarshal JSON
	var pageDataPosted PageData
	json.Unmarshal(bs, &pageDataPosted)
	//Attempt to get data from loaded message map
	someMessages, goodMessageFind := getTenResults(pageDataPosted.ThePage)
	//Declare data to return
	type ReturnMessage struct {
		Messages   []Message `json:"Messages"`
		ResultMsg  string    `json:"ResultMsg"`
		SuccOrFail int       `json:"SuccOrFail"`
	}
	if goodMessageFind == true {
		//Set the current page number server side in case User refreshes
		currentPageNumber = pageDataPosted.ThePage
		//Return failure message
		theReturnMessage := ReturnMessage{
			Messages:   someMessages,
			ResultMsg:  "Page Found",
			SuccOrFail: 0,
		}
		theJSONMessage, err := json.Marshal(theReturnMessage)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(theJSONMessage))
	} else {
		//Return failure message
		theReturnMessage := ReturnMessage{
			Messages:   someMessages,
			ResultMsg:  "Error finding page...",
			SuccOrFail: 1,
		}
		theJSONMessage, err := json.Marshal(theReturnMessage)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(theJSONMessage))
	}
}

/* Called in Ajax from Javascript when a User submits a reply to any message/
reply */
func messageReplyAjax(w http.ResponseWriter, r *http.Request) {
	//Initialize struct for taking messages
	type MessageReply struct {
		ParentMessage Message `json:"ParentMessage"`
		ChildMessage  Message `json:"ChildMessage"`
		CurrentPage   int     `json:"CurrentPage"`
	}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var postedMessageReply MessageReply
	json.Unmarshal(bs, &postedMessageReply)

	/* Format the parent message; we grab this becuase it might have been updated
	from Ajax, but the page hasn't refreshed to give us a new parent value we might
	have from a previous refresh */
	formattedParent := theMessageBoard.AllMessagesMap[postedMessageReply.ParentMessage.MessageID]

	//Determine if this parent is an UberParent
	theMessageUberID := 0
	if formattedParent.IsChild == false {
		theMessageUberID = postedMessageReply.ParentMessage.MessageID
	} else {
		theMessageUberID = formattedParent.UberParentID
	}

	theTimeNow := time.Now()

	//Format the newestMessage
	newestMessage := Message{
		MessageID:       randomIDCreation(),
		UserID:          postedMessageReply.ChildMessage.UserID,
		Messages:        []Message{},
		IsChild:         true,
		HasChildren:     false,
		ParentMessageID: formattedParent.MessageID,
		UberParentID:    theMessageUberID,
		Order:           len(formattedParent.Messages) + 1,
		RepliesAmount:   0,
		TheMessage:      postedMessageReply.ChildMessage.TheMessage,
		DateCreated:     theTimeNow.Format("2006-01-02 15:04:05"),
		LastUpdated:     theTimeNow.Format("2006-01-02 15:04:05"),
	}

	//Update the Message
	uberUpdate(newestMessage, formattedParent)

	//Declare return data and inform Ajax
	type ReturnData struct {
		SuccessMsg     string  `json:"SuccessMsg"`
		SuccessBool    bool    `json:"SuccessBool"`
		SuccessInt     int     `json:"SuccessInt"`
		CreatedMessage Message `json:"CreatedMessage"`
		ParentMessage  Message `json:"ParentMessage"`
	}
	theReturnData := ReturnData{
		SuccessMsg:     "You updated the messages",
		SuccessBool:    true,
		SuccessInt:     0,
		CreatedMessage: newestMessage,
		ParentMessage:  postedMessageReply.ParentMessage,
	}
	dataJSON, err := json.Marshal(theReturnData)
	if err != nil {
		fmt.Println("There's an error marshalling this data")
	}
	fmt.Fprintf(w, string(dataJSON))
}

/* Called in ajax from Javascript when a User submits an original message to a thread */
func messageOriginalAjax(w http.ResponseWriter, r *http.Request) {
	//Initialize struct for taking messages
	type OriginalMessage struct {
		TheMessage  string `json:"TheMessage"`
		TheUserID   int    `json:"TheUserID"`
		TheUsername string `json:"TheUsername"`
	}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//Marshal it into our type
	var postedMessage OriginalMessage
	json.Unmarshal(bs, &postedMessage)

	theTimeNow := time.Now()

	//Format the new Original Message
	newestMessage := Message{
		MessageID:       randomIDCreation(),
		UserID:          postedMessage.TheUserID,
		Messages:        []Message{},
		IsChild:         false,
		HasChildren:     false,
		ParentMessageID: 0,
		UberParentID:    0,
		Order:           len(theMessageBoard.AllOriginalMessages) + 1,
		RepliesAmount:   0,
		TheMessage:      postedMessage.TheMessage,
		DateCreated:     theTimeNow.Format("2006-01-02 15:04:05"),
		LastUpdated:     theTimeNow.Format("2006-01-02 15:04:05"),
	}

	//Insert new Message into database and update on server
	uberInsertNewMessage(newestMessage)

	//Declare return data and inform Ajax
	type DataReturn struct {
		SuccessMsg     string  `json:"SuccessMsg"`
		SuccessBool    bool    `json:"SuccessBool"`
		SuccessInt     int     `json:"SuccessInt"`
		CreatedMessage Message `json:"CreatedMessage"`
		ThePageNow     int     `json:"ThePageNow"`
	}
	theDataReturn := DataReturn{
		SuccessMsg:     "You created a new, original message",
		SuccessBool:    true,
		SuccessInt:     0,
		CreatedMessage: newestMessage,
		ThePageNow:     currentPageNumber,
	}
	dataJSON, err := json.Marshal(theDataReturn)
	if err != nil {
		fmt.Println("There's an error marshalling this data")
	}
	fmt.Fprintf(w, string(dataJSON))
}
