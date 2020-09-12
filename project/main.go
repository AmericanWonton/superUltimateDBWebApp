package main

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

//Here's our User struct
type User struct {
	UserName    string `json:"UserName"`
	Password    string `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
	First       string `json:"First"`
	Last        string `json:"Last"`
	Role        string `json:"Role"`
	UserID      int    `json:"UserID"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

/* Mongo No-SQL Variable Declarations */
type AUser struct { //Using this for Mongo
	UserName    string          `json:"UserName"`
	Password    string          `json:"Password"` //This was formally a []byte but we are changing our code to fit the database better
	First       string          `json:"First"`
	Last        string          `json:"Last"`
	Role        string          `json:"Role"`
	UserID      int             `json:"UserID"`
	DateCreated string          `json:"DateCreated"`
	DateUpdated string          `json:"DateUpdated"`
	Hotdogs     MongoHotDogs    `json:"Hotdogs"`
	Hamburgers  MongoHamburgers `json:"Hamburgers"`
}

type TheUsers struct { //Using this for Mongo
	Users []AUser `json:"Users"`
}

type MongoHotDog struct {
	HotDogType  string   `json:"HotDogType"`
	Condiments  []string `json:"Condiments"`
	Calories    int      `json:"Calories"`
	Name        string   `json:"Name"`
	FoodID      int      `json:"FoodID"`
	UserID      int      `json:"UserID"` //User WHOMST this hotDog belongs to
	PhotoID     int      `json:"PhotoID"`
	PhotoSrc    string   `json:"PhotoSrc"`
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

type MongoHotDogs struct {
	Hotdogs []MongoHotDog `json:"Hotdogs"`
}

type MongoHamburger struct {
	BurgerType  string   `json:"BurgerType"`
	Condiments  []string `json:"Condiments"`
	Calories    int      `json:"Calories"`
	Name        string   `json:"Name"`
	FoodID      int      `json:"FoodID"`
	UserID      int      `json:"UserID"` //User WHOMST this hotDog belongs to
	PhotoID     int      `json:"PhotoID"`
	PhotoSrc    string   `json:"PhotoSrc"`
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

type MongoHamburgers struct {
	Hamburgers []MongoHamburger `json:"Hamburgers"`
}

//Below is our struct for Hotdogs/Hamburgers(standard SQL)
type Hotdog struct {
	HotDogType  string `json:"HotDogType"`
	Condiment   string `json:"Condiment"`
	Calories    int    `json:"Calories"`
	Name        string `json:"Name"`
	UserID      int    `json:"UserID"` //User WHOMST this hotDog belongs to
	FoodID      int    `json:"FoodID"`
	PhotoID     int    `json:"PhotoID"`
	PhotoSrc    string `json:"PhotoSrc"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

type Hamburger struct {
	BurgerType  string `json:"BurgerType"`
	Condiment   string `json:"Condiment"`
	Calories    int    `json:"Calories"`
	Name        string `json:"Name"`
	UserID      int    `json:"UserID"` //User WHOMST this hotDog belongs to
	FoodID      int    `json:"FoodID"`
	PhotoID     int    `json:"PhotoID"`
	PhotoSrc    string `json:"PhotoSrc"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

//Here is our photo struct
type UserPhoto struct {
	UserID      int    `json:"UserID"`
	FoodID      int    `json:"FoodID"`
	PhotoID     int    `json:"PhotoID"`
	PhotoName   string `json:"PhotoName"`
	FileType    string `json:"FileType"`
	Size        int64  `json:"Size"`
	PhotoHash   string `json:"PhotoHash"`
	Link        string `json:"Link"`
	FoodType    string `json:"FoodType"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

//Here is our ViewData struct
type ViewData struct {
	User     User   `json:"User"`
	UserName string `json:"UserName"`
	Role     string `json:"Role"`
}

//Here's our session struct
type theSession struct {
	username     string
	lastActivity time.Time
}

//Session Database info
var dbUsers = map[string]User{}          // user ID, user
var dbSessions = map[string]theSession{} // session ID, session
var dbSessionsCleaned time.Time

//mySQL database declarations
var db *sql.DB
var err error

//Mongo DB Declarations
var mongoClient *mongo.Client

//Here is our waitgroup
var wg sync.WaitGroup

const sessionLength int = 180 //Length of sessions

/* TEMPLATE DEFINITION BEGINNING */
var template1 *template.Template

func logWriter(logMessage string) {
	//Logging info
	/*
		fmt.Println("Writing log files.")
		logFile, err := os.OpenFile("/tmp/superdblogs/superDBAppLog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		defer logFile.Close()

		if err != nil {
			//log.Fatalln("Failed opening file")
			fmt.Println("Failed opening file")
		}

		log.SetOutput(logFile)

		log.Println(logMessage)
	*/
	fmt.Printf("DEBUG: Writing the following message in logwriter\n\n%v\n", logMessage)
}

/* FUNCMAP DEFINITION */
/* DEBUG, I'M NOT SURE IF WE NEED THESE RETURN ROLE USERS */
func (u User) ReturnRoleUser(theUser string) bool {
	if strings.Compare(theUser, "user") == 0 {
		fmt.Printf("DEBUG: WE ARE IN RETURN TRUE USER")
		return true
	} else {
		return false
	}
}

func (u User) ReturnRoleAdmin(theAdmin string) bool {
	if strings.Compare(theAdmin, "admin") == 0 {
		fmt.Printf("DEBUG: WE ARE IN RETURN TRUE ADMIN")
		return true
	} else {
		return false
	}
}

func (u User) ReturnRoleIT(theIT string) bool {
	if strings.Compare(theIT, "IT") == 0 {
		return true
	} else {
		return false
	}
}

var funcMap = template.FuncMap{
	"upperCase":       strings.ToUpper, //upperCase is a key we can call inside of the template html file
	"ReturnRoleUser":  User.ReturnRoleUser,
	"ReturnRoleAdmin": User.ReturnRoleAdmin,
	"ReturnRoleIT":    User.ReturnRoleIT,
}

//Loading our templates in for ParseGlob: https://github.com/gobuffalo/packr/issues/16
var templatesBox = packr.New("templates", "./static")

func MustLoadBoxedTemplate(b *packr.Box) *template.Template {
	t := template.Must(template.New("").Funcs(funcMap), err)
	//t := template.New("")
	err := b.Walk(func(p string, f packr.File) error {
		if p == "" {
			return nil
		}
		var err error
		var csz int64
		if finfo, err := f.FileInfo(); err != nil {
			return err
		} else {
			// skip directory path
			if finfo.IsDir() {
				return nil
			}
			csz = finfo.Size()
		}

		// skip all files except .html
		if !strings.Contains(p, ".html") && !strings.Contains(p, ".gohtml") {
			//fmt.Printf("We are skipping this filename: %v\n", p)
			return nil
		}

		// Normalize template name
		n := p
		if strings.HasPrefix(p, "\\") || strings.HasPrefix(p, "/") {
			n = n[1:] // don't want template name to start with / ie. /index.html
		}
		// replace windows path seperator \ to normalized /
		n = strings.Replace(n, "\\", "/", -1)

		var h = make([]byte, 0, csz)

		if h, err = b.Find(p); err != nil {
			return err
		}

		if _, err = t.New(n).Parse(string(h)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic("error loading template")
	}
	return t
}

//Parse our templates
func init() {
	template1 = MustLoadBoxedTemplate(templatesBox)
	//AmazonCredentialRead
	getCreds()
}

// Handle Errors
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Home page
func homePage(w http.ResponseWriter, r *http.Request) {
	aUser := getUser(w, r) //Get the User, if they exist

	//If a User posts a form to log in!
	//Search for Users in Database, send JSON version of User
	if r.Method == http.MethodPost {
		//Get Form Values
		username := r.FormValue("username")
		password := r.FormValue("password")
		//Query database for those username and password
		rows, err := db.Query(`SELECT * FROM users WHERE USERNAME = ?;`, username)
		check(err)
		defer rows.Close()
		//Count to see if password is found or not
		var returnedTableID int = 0
		var returnedUsername string = ""
		var returnedPassword string = ""
		var returnedFName string = ""
		var returnedLName string = ""
		var returnedRole string = ""
		var returnedUserID int = 0
		var returnedDateCreated string = ""
		var returnedDateUpdated string = ""

		for rows.Next() {
			//assign variable
			err = rows.Scan(&returnedTableID, &returnedUsername, &returnedPassword, &returnedFName, &returnedLName, &returnedRole, &returnedUserID,
				&returnedDateCreated, &returnedDateUpdated)
			check(err)
		}
		//Count to see if password/Username returned at all
		if (strings.Compare(returnedUsername, "") == 0) || (strings.Compare(returnedPassword, "") == 0) {
			fmt.Printf("Username, %v and %v, and Password, %v and %v not Found!\n", returnedUsername, "", returnedPassword, "")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return //DEBUG Not sure if this is needed or wanted
		} else {
			//Check to see if Username returned
			if strings.Compare(username, returnedUsername) == 0 {
				//Checking to see if password matches as well
				theReturnedByte, err := hex.DecodeString(returnedPassword)
				if err != nil {
					log.Fatal(err)
				}
				if strings.Compare(string(theReturnedByte), password) != 0 {
					//Password not found/not hashed correctly
					fmt.Printf("The hashed strings aren't compatable: %v %v\n", string(theReturnedByte), password)
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				} else {
					//Username matched, password matched good stuff
					//User logged in, directing them to the mainpage
					//Going to main page, passing values
					theUser := User{username, returnedPassword, returnedFName, returnedLName, returnedRole, returnedUserID,
						returnedDateCreated, returnedDateUpdated}
					dbUsers[username] = theUser
					// create session
					uuidWithHyphen := uuid.New().String()

					cookie := &http.Cookie{
						Name:  "session",
						Value: uuidWithHyphen,
					}
					cookie.MaxAge = sessionLength
					http.SetCookie(w, cookie)
					dbSessions[cookie.Value] = theSession{username, time.Now()}
					http.Redirect(w, r, "/mainPage", http.StatusSeeOther)
					return
				}
			} else {
				//Passwords do not match
				fmt.Printf("Username, %v and %v or password, %v, did not match!\n", username, returnedUsername,
					returnedPassword)
			}
		}
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "templates/index.gohtml", aUser)
	HandleError(w, err1)
}

//signUp
func signUp(w http.ResponseWriter, req *http.Request) {
	//See if user is already logged in
	if alreadyLoggedIn(w, req) {
		//If already logged in, put them back at the main menu
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	err1 := template1.ExecuteTemplate(w, "templates/signup.gohtml", nil)
	HandleError(w, err1)

	fmt.Printf("Signup Endpoint Hit\n")
}

func signUpUserUpdated(w http.ResponseWriter, req *http.Request) {
	// process Ajax ping
	if req.Method == http.MethodPost {
		fmt.Println("DEBUG: We got the Ajax and are inserting the User.")
		//Collect JSON from Postman or wherever
		//Get the byte slice from the request body ajax
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}

		//Marshal it into our type
		var postedUser User
		json.Unmarshal(bs, &postedUser)
		// get form values
		username := postedUser.UserName
		password := postedUser.Password
		firstname := postedUser.First
		lastname := postedUser.Last
		role := postedUser.Role
		// create session
		uuidWithHyphen := uuid.New().String()
		newCookie := &http.Cookie{
			Name:  "session",
			Value: uuidWithHyphen,
		}
		newCookie.MaxAge = sessionLength
		http.SetCookie(w, newCookie)
		dbSessions[newCookie.Value] = theSession{username, time.Now()}
		// store user in dbUsers
		//Make User and USERID
		theID := randomIDCreation()

		fmt.Println("DEBUG: Adding User data to SQL database")
		//Add User to the SQL Database
		bsString := []byte(password)                  //Encode Password
		encodedString := hex.EncodeToString(bsString) //Encode Password Pt2
		theTimeNow := time.Now()
		var insertedUser User = User{
			UserName:    username,
			Password:    encodedString,
			First:       firstname,
			Last:        lastname,
			Role:        role,
			UserID:      theID,
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		jsonValue, _ := json.Marshal(insertedUser)
		response, err := http.Post("http://3.135.9.238:80/insertUser", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}

		//Add User to MongoDB
		fmt.Printf("DEBUG: Adding User to MongoDB\n")
		var insertionUser AUser = AUser{
			UserName:    username,
			Password:    encodedString,
			First:       firstname,
			Last:        lastname,
			Role:        role,
			UserID:      theID,
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
			Hotdogs:     MongoHotDogs{},
			Hamburgers:  MongoHamburgers{},
		}
		insertionUsers := TheUsers{
			Users: []AUser{insertionUser},
		}
		jsonValue2, _ := json.Marshal(insertionUsers)
		response2, err := http.Post("http://3.135.9.238:80/insertUsers", "application/json", bytes.NewBuffer(jsonValue2))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response2.Body)
			fmt.Println(string(data))
		}

		//DEBUG, don't know if we need below
		var theUser = User{username, encodedString, firstname, lastname, role, theID, insertionUser.DateCreated,
			insertionUser.DateUpdated}
		dbUsers[username] = theUser
		type successMSG struct {
			Message    string `json:"Message"`
			SuccessNum int    `json:"SuccessNum"`
		}
		msgSuccess := successMSG{
			Message:    "Added the new account!",
			SuccessNum: 0,
		}
		theJSONMessage, err := json.Marshal(msgSuccess)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(theJSONMessage))
	}
}

//mainPage
func mainPage(w http.ResponseWriter, req *http.Request) {
	//if User is already logged in, bring them to the mainPage!
	aUser := getUser(w, req) //Get the User, if they exist
	aUserRole := aUser.Role
	vd := ViewData{aUser, aUser.UserName, aUserRole}
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	err1 := template1.ExecuteTemplate(w, "templates/mainpage.gohtml", vd)
	HandleError(w, err1)
}

//Handles all requests coming in
func handleRequests() {
	fmt.Printf("DEBUG: Handling Requests...\n")
	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/signup", signUp)
	myRouter.HandleFunc("/mainPage", mainPage)
	myRouter.HandleFunc("/signUpUserUpdated", signUpUserUpdated)
	//SQL Database Stuff
	myRouter.HandleFunc("/deleteFood", deleteFood).Methods("POST")
	myRouter.HandleFunc("/updateFood", updateFood).Methods("POST")           //Update a certain food item
	myRouter.HandleFunc("/insertHotDog", insertHotDog).Methods("POST")       //Post a hotdog!
	myRouter.HandleFunc("/insertHamburger", insertHamburger).Methods("POST") //Post a hamburger!
	myRouter.HandleFunc("/getAllFoodUser", getAllFoodUser).Methods("POST")   //Get all foods for a User ID
	myRouter.HandleFunc("/getHotDog", getHotDog).Methods("GET")              //Get a SINGULAR hotdog
	myRouter.HandleFunc("/insertUser", insertUser).Methods("POST")           //Post a User!
	myRouter.HandleFunc("/getUsers", getUsers).Methods("GET")                //Get a Users!
	myRouter.HandleFunc("/updateUsers", updateUsers).Methods("POST")         //Get a Users!
	myRouter.HandleFunc("/deleteUsers", deleteUsers).Methods("POST")         //DELETE a Users!
	//Mongo No-SQL Stuff
	myRouter.HandleFunc("/insertUsers", insertUsers).Methods("POST")                   //Post a User!
	myRouter.HandleFunc("/insertHotDogs", insertHotDogs).Methods("POST")               //Post Hotdogs!
	myRouter.HandleFunc("/insertHamburgers", insertHamburgers).Methods("POST")         //Post Hamburgers!
	myRouter.HandleFunc("/insertHotDogMongo", insertHotDogMongo).Methods("POST")       //Post Hamburgers!
	myRouter.HandleFunc("/insertHamburgerMongo", insertHamburgerMongo).Methods("POST") //Post Hamburgers!
	myRouter.HandleFunc("/foodUpdateMongo", foodUpdateMongo).Methods("POST")           //Post Food Update!
	myRouter.HandleFunc("/getAllFoodMongo", getAllFoodMongo).Methods("POST")           //Post All Foods to get!
	myRouter.HandleFunc("/randomIDCreationAPI", randomIDCreationAPI).Methods("POST")   //Get Random IDS
	myRouter.HandleFunc("/foodDeleteMongo", foodDeleteMongo).Methods("POST")           //Delete some Foods
	//Database Insertion stuff
	myRouter.HandleFunc("/hotDogInsertWebPage", hotDogInsertWebPage).Methods("POST")       //Post Hotdogs
	myRouter.HandleFunc("/hamburgerInsertWebPage", hamburgerInsertWebPage).Methods("POST") //Post Hamburgers
	//File Handling Stuff
	myRouter.HandleFunc("/fileInsert", fileInsert).Methods("POST") //Insert a file
	myRouter.HandleFunc("/checkSRC", checkSRC).Methods("POST")     //Check if directory exists
	//Validation Stuff
	myRouter.HandleFunc("/checkUsername", checkUsername) //Check Username
	myRouter.HandleFunc("/loadUsernames", loadUsernames) //Loads in Usernames
	//Middleware logging
	myRouter.Handle("/", loggingMiddleware(http.HandlerFunc(logHandler)))
	//Serve our static files
	myRouter.Handle("/", http.FileServer(templatesBox))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(templatesBox)))
	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed
	//Write initial entry to log
	//logWriter("Deployed superUltimateDBWebApp app.")
	//open SQL connection
	db, err = sql.Open("mysql",
		"joek1:fartghookthestrong69@tcp(food-database.cd8ujtto1hfj.us-east-2.rds.amazonaws.com)/food-database-schema?charset=utf8")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	//Mongo Connect
	mongoClient = connectDB()
	//Handle Requests
	handleRequests()
	defer mongoClient.Disconnect(theContext) //Disconnect in 10 seconds if you can't connect
}

//Check errors in our mySQL errors
func check(err error) {
	if err != nil {
		fmt.Printf("We got an error somewhere, printing it out: %v\n", err.Error())
	}
}

//Some stuff for logging
func logHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Package main, son")
	fmt.Fprint(w, "package main, son.")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logrus.Infof("uri: %v\n", req.RequestURI)
		next.ServeHTTP(w, req)
	})
}

/* DEBUG ZONE */
