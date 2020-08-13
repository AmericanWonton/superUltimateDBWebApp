package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Here's our User struct
type User struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
	First    string `json:"First"`
	Last     string `json:"Last"`
	Role     string `json:"Role"`
	UserID   int    `json:"UserID"`
}

//Below is our struct for Hotdogs/Hamburgers
type Hotdog struct {
	HotDogType string `json:"HotDogType"`
	Condiment  string `json:"Condiment"`
	Calories   int    `json:"Calories"`
	Name       string `json:"Name"`
	UserID     int    `json:"UserID"` //User WHOMST this hotDog belongs to
}

type Hamburger struct {
	BurgerType string `json:"BurgerType"`
	Condiment  string `json:"Condiment"`
	Calories   int    `json:"Calories"`
	Name       string `json:"Name"`
	UserID     int    `json:"UserID"` //User WHOMST this hotDog belongs to
}

//Data to be created and sent
type SendData struct {
	TheUsers      []User
	TheHotdogs    []Hotdog
	TheHamburgers []Hamburger
}

//Here's our session struct
type session struct {
	username     string
	lastActivity time.Time
}

//mySQL database declarations
var db *sql.DB
var err error

const sessionLength int = 180 //Length of sessions
const min int = 1
const max int = 3

func handleRequests() {
	fmt.Println("We handling Requests")
	myRouter := mux.NewRouter().StrictSlash(true)
	//Database Stuff
	myRouter.HandleFunc("/insertHotDog", insertHotDog).Methods("POST")       //Post a hotdog!
	myRouter.HandleFunc("/insertHamburger", insertHamburger).Methods("POST") //Post a hamburger!
	myRouter.HandleFunc("/insertUser", insertUser).Methods("POST")           //Post a User!
	myRouter.HandleFunc("/getUsers", getUsers).Methods("GET")                //Get a Users!

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	//open SQL connection
	db, err = sql.Open("mysql",
		"joek1:fartghookthestrong69@tcp(food-database.cd8ujtto1hfj.us-east-2.rds.amazonaws.com)/food-database-schema?charset=utf8")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	handleRequests()

	var dataToSend SendData
	dataToSend.TheUsers = createSendDataUser()
	dataToSend.TheHotdogs = createSendDataHDog(dataToSend.TheUsers)
	dataToSend.TheHamburgers = createSendDataHam(dataToSend.TheUsers)

	fmt.Printf("Here is our data to send: \n\n%v\n\n", dataToSend)
}

//Check errors in our mySQL errors
func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func createSendDataUser() []User {
	fmt.Println("creating sendUserData")
	ourUsers := []User{}
	//Create 5 Users and return them to the data to send
	for i := 0; i < 5; i++ {
		//Get Random Name Example
		url := "https://randomuser.me/api/?nat=us"
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Cookie", "__cfduid=d1f4b1fdbc2cdb4c1d8a02a18a24cc4641590522430")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		//Assign the JSON to a person
		var randPerson RandomPerson
		json.Unmarshal(body, &randPerson)
		role := ""
		switch randomRoleAssign := rand.Intn(max-min) + min; randomRoleAssign {
		case 1:
			role = "user"
		case 2:
			role = "admin"
		case 3:
			role = "IT"
		}
		ourID := createUserID()
		firstName := randPerson.Results[0].Name.First
		lastName := randPerson.Results[0].Name.Last
		username := randomUsername(firstName, lastName)
		password := randomPassword()
		fmt.Printf("Here is our randPerson Firstname: \n%v\n", firstName)
		fmt.Printf("Here is our randPerson Lasname: \n%v\n", lastName)
		fmt.Printf("Here is our random Role: %v\n", role)
		fmt.Printf("Here is our random userID: %v\n", ourID)
		fmt.Printf("Here is our random password: %v\n", password)
		fmt.Printf("Here is our random username: %v\n", username)

		ourUser := User{
			UserName: username,
			Password: password,
			First:    firstName,
			Last:     lastName,
			Role:     role,
			UserID:   ourID,
		}

		fmt.Println(ourUser)
		ourUsers = append(ourUsers, ourUser)
		//Add the User to our database
		jsonValue, _ := json.Marshal(ourUser)
		response, err := http.Post("http://3.135.9.238:8080/insertUser", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	}

	return ourUsers
}

func createSendDataHDog(ourUsers []User) []Hotdog {
	fmt.Println("creating hot dog Data")
	ourDogs := []Hotdog{}
	someTypes := []string{"Atlantian", "Space", "Sanitized", "Sweaty",
		"Sexy", "Senior", "Recyclable"}
	someCondiments := []string{"Fish", "Carbon", "HydroAlcahol", "Pit Stains", "Makeup", "Reading Glasses", "Glass"}
	selectionmin, selectionmax := 0, 6 //The min and Max value for our randInt
	caloriesmin, caloriesmax := 1, 2000
	someNames := []string{"The Atlantisdog", "The Trekdog", "The Cleandog", "The Hothotdog", "The Sexydog",
		"The Oldiesdog", "The Greendog"}
	//Generate our random hotdogs and give them to our database
	for j := 0; j < len(ourUsers); j++ {
		//Build a random Hotdog
		aHotdog := Hotdog{}
		theID := ourUsers[j].UserID
		theRandNum := rand.Intn(selectionmax-selectionmin) + selectionmin
		aHotdog.HotDogType = someTypes[theRandNum]
		aHotdog.Condiment = someCondiments[theRandNum]
		aHotdog.Calories = rand.Intn(caloriesmax-caloriesmin) + caloriesmin
		aHotdog.Name = someNames[theRandNum]
		aHotdog.UserID = theID
		//Append it to ourDogs
		ourDogs = append(ourDogs, aHotdog)
		//Add the hotdog to the database
		//Add the User to our database
		jsonValue, _ := json.Marshal(aHotdog)
		response, err := http.Post("http://3.135.9.238:8080/insertHotDog", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	}

	return ourDogs
}

func createSendDataHam(ourUsers []User) []Hamburger {
	fmt.Println("creating Hamburger Data")
	ourHams := []Hamburger{}
	someTypes := []string{"Wet", "Wall", "Dance", "Sick",
		"Cellular", "Memorable", "Recyclable"}
	someCondiments := []string{"H20", "90 degree angle", "Boogie", "Vomit", "Galaxy S8", "Childhood Memories", "Glass"}
	selectionmin, selectionmax := 0, 6 //The min and Max value for our randInt
	caloriesmin, caloriesmax := 1, 2000
	someNames := []string{"The OceanBurger", "The UprightBurger", "The GroovyBurger", "The PestilenceBurger", "The PhoneBurger",
		"The NostalgiaBurger", "The GreenBurger"}
	//Generate our random hotdogs and give them to our database
	for j := 0; j < len(ourUsers); j++ {
		//Build a random Hotdog
		aBurger := Hamburger{}
		theID := ourUsers[j].UserID
		theRandNum := rand.Intn(selectionmax-selectionmin) + selectionmin
		aBurger.BurgerType = someTypes[theRandNum]
		aBurger.Condiment = someCondiments[theRandNum]
		aBurger.Calories = rand.Intn(caloriesmax-caloriesmin) + caloriesmin
		aBurger.Name = someNames[theRandNum]
		aBurger.UserID = theID
		//Append it to ourDogs
		ourHams = append(ourHams, aBurger)
		//Add the hotdog to the database
		//Add the User to our database
		jsonValue, _ := json.Marshal(aBurger)
		response, err := http.Post("http://3.135.9.238:8080/insertHamburger", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	}

	return ourHams
}
