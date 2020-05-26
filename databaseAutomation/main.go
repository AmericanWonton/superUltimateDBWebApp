package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
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

func main() {
	fmt.Println("We ran code.")

	//open SQL connection
	db, err = sql.Open("mysql",
		"joek1:fartghookthestrong69@tcp(food-database.cd8ujtto1hfj.us-east-2.rds.amazonaws.com)/food-database-schema?charset=utf8")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	createSendData()
}

//Check errors in our mySQL errors
func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func createSendData() {
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
	fmt.Printf("Here is our randPerson Name: \n%v\n", firstName)
	fmt.Printf("Here is our randPerson Name: \n%v\n", lastName)
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
}

func createUserID() int {
	//Make User and USERID
	goodNum := false
	theID := 0

	for goodNum == false {
		//Build the random, unique integer to be assigned to this User
		randInt := 0                 //The random integer added onto ID
		randIntString := ""          //The integer built through a string...
		useridmin, useridmax := 0, 9 //The min and Max value for our randInt
		var idCount int              //A count of how many times our ID is in the database
		for i := 0; i < 8; i++ {
			randInt = rand.Intn(useridmax-useridmin) + useridmin
			randIntString = randIntString + strconv.Itoa(randInt)
		}
		theID, err = strconv.Atoi(randIntString)
		if err != nil {
			fmt.Println(err)
			idCount = 2
		}
		//Check to see if ID is in database
		row, err := db.Query("SELECT user_id FROM users WHERE USER_ID=?;", theID)
		check(err)
		defer row.Close()

		for row.Next() {
			idCount = idCount + 1
		}

		if idCount >= 1 {
			goodNum = false
		} else {
			goodNum = true
		}
	}

	return theID
}
