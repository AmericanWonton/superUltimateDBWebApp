package main

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

//Information for serving files/testing
var serverAddress string

/* INFORMATION FOR OUR EMAIL VARIABLES */
var senderAddress string
var senderPWord string

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
	User           User   `json:"User"`
	UserName       string `json:"UserName"`
	Role           string `json:"Role"`
	Port           string `json:"Port"`
	MessageDisplay int    `json:"MessageDisplay"` //This is IF we need a message displayed
	Message        string `json:"Message"`        //This is the message to display
	IsSuccess      int    `json:"IsSuccess"`      //Determines if the message is successful
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

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "superDBAppLog.txt")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
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

//Parse our templates
func init() {
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
	//AmazonCredentialRead
	getCreds()
	OAuthGmailService() //Initialize Gmail Services
}

// Handle Errors passing templates
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
					http.Redirect(w, r, "/choicepage", http.StatusSeeOther)
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
	err1 := template1.ExecuteTemplate(w, "index.gohtml", aUser)
	HandleError(w, err1)
}

//signup
func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: We are in the signup endpoint\n")
	//See if user is already logged in
	if alreadyLoggedIn(w, r) {
		//If already logged in, put them back at the main menu
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err1 := template1.ExecuteTemplate(w, "signup.gohtml", nil)
	HandleError(w, err1)
}

//Begins Sending Email to User and creates a User for database entry
func signUpUserUpdated(w http.ResponseWriter, req *http.Request) {
	// process Ajax ping
	if req.Method == http.MethodPost {
		fmt.Printf("DEBUG: We submitted a ajax form, now in signUpUserUpdated \n")
		//Collect JSON from Postman or wherever
		//Get the byte slice from the request body ajax
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}

		//Declare DataType from Ajax
		type UserData struct {
			TheUser User   `json:"TheUser"`
			Email   string `json:"Email"`
		}

		//Marshal the user data into our type
		var dataPosted UserData
		json.Unmarshal(bs, &dataPosted)
		//Set the User info
		var postedUser User = dataPosted.TheUser
		// get form values
		username := postedUser.UserName
		password := postedUser.Password
		firstname := postedUser.First
		lastname := postedUser.Last
		role := postedUser.Role
		email := dataPosted.Email
		/* ATTEMPT TO SEND EMAIL...IF IT FAILS, DO NOT CREATE USER */
		goodEmailSend := signUpUserEmail(email, role, firstname, lastname)
		if goodEmailSend == true {
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
			response, err := http.Post("http://"+serverAddress+"/insertUser", "application/json", bytes.NewBuffer(jsonValue))
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
			response2, err := http.Post("http://"+serverAddress+"/insertUsers", "application/json", bytes.NewBuffer(jsonValue2))
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
			//Alert Ajax with success

			type successMSG struct {
				Message     string `json:"Message"`
				SuccessNum  int    `json:"SuccessNum"`
				RedirectURL string `json:"RedirectURL"`
			}
			msgSuccess := successMSG{
				Message:     "Added the new account!",
				SuccessNum:  0,
				RedirectURL: "http://" + serverAddress,
			}

			theJSONMessage, err := json.Marshal(msgSuccess)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("DEBUG: Writing back successful User response.\n")
			fmt.Fprint(w, string(theJSONMessage))

			/* EVERY CHECK FOR CREATING A USER IS SUCCESSFUL. REDIRECT TO THE HOMEPAGE */
			fmt.Printf("DEBUG: SHOULD BE REDIRECTING NOW...\n")
		} else {
			fmt.Printf("DEBUG: YOU FAILED TO CREATE USER\n")
			//Alert Ajax with failure
			type successMSG struct {
				Message    string `json:"Message"`
				SuccessNum int    `json:"SuccessNum"`
			}
			msgSuccess := successMSG{
				Message:    "Failed to send email and create User.",
				SuccessNum: 1,
			}
			theJSONMessage, err := json.Marshal(msgSuccess)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(w, string(theJSONMessage))
		}
	}
}

//mainPage (also known as the database page)
func mainPage(w http.ResponseWriter, r *http.Request) {
	//if User is already logged in, bring them to the mainPage!
	aUser := getUser(w, r) //Get the User, if they exist
	aUserRole := aUser.Role
	thePort := os.Getenv("PORT")
	if thePort == "" {
		thePort = "80"
		logWriter("Defautling to this port: " + thePort)
	}
	vd := ViewData{aUser, aUser.UserName, aUserRole, thePort, 1, "Welcome to the Main page!", 0}
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//See if there is a submission for new food or updates/deletes
	if r.Method == http.MethodPost {
		//Determine what kind of form post this is
		theAction := r.FormValue("hiddenFoodAction")
		if strings.Contains(strings.ToLower(theAction), strings.ToLower("food_submit")) {
			//Get form stuff Photo Stuff
			maxSize := int64(1024000) // allow only 1MB of file size
			err := r.ParseMultipartForm(maxSize)
			if err != nil {
				fmt.Printf("Image too large. Max Size: %v\n", maxSize)
				log.Println(err)
				return
			}

			file, fileHeader, err := r.FormFile("newFile") //Insert name of file element here
			if err != nil {
				errMsg := "Could not get uploaded file in food submit Error getting file submission: " + err.Error()
				fmt.Println(errMsg)
				logWriter(errMsg)
			}
			defer file.Close()
			fmt.Printf("DEBUG: Here's Fileheader: %v\n", fileHeader.Filename)
			//Decide if this is hamburger or hotdog submission
			hiddenFoodType := r.FormValue("hiddenFoodType")
			if strings.Contains(strings.ToUpper(hiddenFoodType), "HAMBURGER") {
				foodType := r.FormValue("hamburgType")
				condimentType := r.FormValue("condimentType")
				caloriesType := r.FormValue("caloriesType")
				theCalories, _ := strconv.Atoi(caloriesType)
				nameType := r.FormValue("nameType")
				userIDInput := r.FormValue("userIDInput")
				theUserID, _ := strconv.Atoi(userIDInput)
				hiddenUserNum := r.FormValue("hiddenUserNum")
				hiddenFoodNum := r.FormValue("hiddenFoodNum")
				fmt.Printf("Heres hiddenUser: %v\n Here's HiddenFood: %v\n", hiddenUserNum, hiddenFoodNum)

				//Create Food to send
				theTimeNow := time.Now()
				//curDir, _ := os.Getwd()
				photoDir := filepath.Join("amazonimages", "pictures", userIDInput, strings.ToUpper(hiddenFoodType),
					fileHeader.Filename)
				finalURL := urlFixer(photoDir)
				sendHamburger := Hamburger{
					BurgerType:  foodType,
					Condiment:   condimentType,
					Calories:    theCalories,
					Name:        nameType,
					UserID:      theUserID,
					FoodID:      randomIDCreation(),
					PhotoID:     randomIDCreation(),
					PhotoSrc:    finalURL,
					DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
					DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
				}
				sendHotdog := Hotdog{}
				fmt.Printf("DEBUG, here's our hdog: %v\n and Ham: %v\n", sendHamburger, sendHotdog)
				//InsertFood to dbs and activate AWS variables
				goodInsert := simpleFoodInsert("HAMBURGER", sendHotdog, sendHamburger)
				if goodInsert == true {
					//Put file in directory AND send to Amazon
					_, goodFileInsert, hexName, _ := fileInsert(w, r)
					if goodFileInsert == true {
						//Add the photo details to the database
						//Upload photo details to DB
						extension := filepath.Ext(fileHeader.Filename)
						fileURL := filepath.Join("pictures", userIDInput, awsfoodType, hexName+extension)
						insertedPhoto := insertUserPhotos(awsuserID, awsfoodID, sendHamburger.PhotoID,
							fileHeader.Filename, extension, fileHeader.Size,
							hexName, fileURL, awsfoodType, "", "")
						if insertedPhoto == true {
							succMsg := "Inserted photo information into SQL DB"
							logWriter(succMsg)
							fmt.Println("DEBUG: Inserted photo information into SQL DB.") //Inform User
							//INSERT INTO MONGO
							mongoInsertion := mongoInsertPhoto(awsuserID, awsfoodID, sendHamburger.PhotoID,
								fileHeader.Filename, extension, fileHeader.Size,
								hexName, fileURL, awsfoodType, "", "")
							if mongoInsertion == true {
								succMsgTwo := "Photo information successfully submitted into MongoDB"
								logWriter(succMsgTwo)
								fmt.Println(succMsgTwo)
								//Assemble data to be sent to inform User
								vd.MessageDisplay = 0
								vd.Message = succMsgTwo
								vd.IsSuccess = 0
								err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
								HandleError(w, err1)
							} else {
								errMsg := "Issue inserting photo into MongoDB"
								logWriter(errMsg)
								fmt.Println("DEBUG: Issue inserting photo into MongoDB")
								//Assemble data to be sent to inform User
								vd.MessageDisplay = 0
								vd.Message = errMsg
								vd.IsSuccess = 1
								err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
								HandleError(w, err1)
								return
							}
						} else {
							errMsg := "Issue inserting photo information into SQL DB; insertedPhoto is false"
							logWriter(errMsg)
							fmt.Println(errMsg)
							//Assemble data to be sent to inform User
							vd.MessageDisplay = 0
							vd.Message = errMsg
							vd.IsSuccess = 1
							err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
							HandleError(w, err1)
							return
						}
					} else {
						//Assemble data to be sent to inform User
						vd.MessageDisplay = 0
						vd.Message = "File was unable to be inserted for food"
						vd.IsSuccess = 1
						err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
						HandleError(w, err1)
						return
					}
				} else {
					errMsg := "Error, goodInsert is false"
					logWriter(errMsg)
					fmt.Println(errMsg)
					//Assemble data to be sent to inform User
					vd.MessageDisplay = 0
					vd.Message = errMsg
					vd.IsSuccess = 1
					err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
					HandleError(w, err1)
					return
				}
			} else if strings.Contains(strings.ToUpper(hiddenFoodType), "HOTDOG") {
				foodType := r.FormValue("hDogType")
				condimentType := r.FormValue("condimentType")
				caloriesType := r.FormValue("caloriesType")
				theCalories, _ := strconv.Atoi(caloriesType)
				nameType := r.FormValue("nameType")
				userIDInput := r.FormValue("userIDInput")
				theUserID, _ := strconv.Atoi(userIDInput)
				hiddenUserNum := r.FormValue("hiddenUserNum")
				hiddenFoodNum := r.FormValue("hiddenFoodNum")
				fmt.Printf("DEBUG: Heres hiddenUser: %v\n Here's HiddenFood: %v\n", hiddenUserNum, hiddenFoodNum)

				//Create Food to send
				theTimeNow := time.Now()
				//curDir, _ := os.Getwd()
				photoDir := filepath.Join("amazonimages", "pictures", userIDInput, strings.ToUpper(hiddenFoodType),
					fileHeader.Filename)
				finalURL := urlFixer(photoDir)
				sendHamburger := Hamburger{}
				sendHotdog := Hotdog{
					HotDogType:  foodType,
					Condiment:   condimentType,
					Calories:    theCalories,
					Name:        nameType,
					UserID:      theUserID,
					FoodID:      randomIDCreation(),
					PhotoID:     randomIDCreation(),
					PhotoSrc:    finalURL,
					DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
					DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
				}
				//InsertFood to dbs and activate AWS variables
				goodInsert := simpleFoodInsert("HOTDOG", sendHotdog, sendHamburger)
				if goodInsert == true {
					//Put file in directory AND send to Amazon
					_, goodFileInsert, hexName, _ := fileInsert(w, r)
					if goodFileInsert == true {
						//Add the photo details to the database
						fmt.Printf("DEBUG: Uploading file to SQL database.\n")
						//Upload photo details to DB
						extension := filepath.Ext(fileHeader.Filename)
						fileURL := filepath.Join("pictures", userIDInput, awsfoodType, hexName+extension)
						insertedPhoto := insertUserPhotos(awsuserID, awsfoodID, sendHamburger.PhotoID,
							fileHeader.Filename, extension, fileHeader.Size,
							hexName, fileURL, awsfoodType, "", "")
						if insertedPhoto == true {
							succMsg := "Inserted photo information into SQL DB"
							logWriter(succMsg)
							fmt.Println("DEBUG: Inserted photo information into SQL DB.") //Inform User
							//INSERT INTO MONGO
							mongoInsertion := mongoInsertPhoto(awsuserID, awsfoodID, sendHamburger.PhotoID,
								fileHeader.Filename, extension, fileHeader.Size,
								hexName, fileURL, awsfoodType, "", "")
							if mongoInsertion == true {
								succMsgTwo := "Photo information successfully submitted into MongoDB"
								logWriter(succMsgTwo)
								fmt.Println(succMsgTwo)
								//Assemble data to be sent to inform User
								vd.MessageDisplay = 0
								vd.Message = succMsgTwo
								vd.IsSuccess = 0
								err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
								HandleError(w, err1)
								return
							} else {
								errMsg := "Issue inserting photo into MongoDB"
								logWriter(errMsg)
								fmt.Println("DEBUG: Issue inserting photo into MongoDB")
								//Assemble data to be sent to inform User
								vd.MessageDisplay = 0
								vd.Message = errMsg
								vd.IsSuccess = 1
								err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
								HandleError(w, err1)
								return
							}
						} else {
							errMsg := "Issue inserting photo information into SQL DB; insertedPhoto is false"
							logWriter(errMsg)
							fmt.Println(errMsg)
							//Assemble data to be sent to inform User
							vd.MessageDisplay = 0
							vd.Message = errMsg
							vd.IsSuccess = 1
							err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
							HandleError(w, err1)
							return
						}
					} else {
						//Assemble data to be sent to inform User
						vd.MessageDisplay = 0
						vd.Message = "File was not able to insert for the given food"
						vd.IsSuccess = 1
						err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
						HandleError(w, err1)
						return
					}
				} else {
					errMsg := "Error, goodInsert is false"
					logWriter(errMsg)
					fmt.Println(errMsg)
					//Assemble data to be sent to inform User
					vd.MessageDisplay = 0
					vd.Message = errMsg
					vd.IsSuccess = 1
					err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
					HandleError(w, err1)
					return
				}
			} else {
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.IsSuccess = 1
				vd.Message = "Error, incorrect hiddenFoodType: " + hiddenFoodType
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}
		} else if strings.Contains(strings.ToLower(theAction), strings.ToLower("food_update")) {
			hiddenHasPhoto := r.FormValue("hiddenHasPhoto")
			//Determine if we need to update a photo or not
			if strings.Contains(strings.ToLower(hiddenHasPhoto), "has_photo") {
				fmt.Println("DEBUG: Updtaing food with photo")
				//Parse Form
				maxSize := int64(1024000) // allow only 1MB of file size
				err := r.ParseMultipartForm(maxSize)
				if err != nil {
					fmt.Printf("Image too large. Max Size: %v\n", maxSize)
					log.Println(err)
					//Assemble data to be sent to inform User
					vd.MessageDisplay = 0
					vd.Message = "Image too large."
					vd.IsSuccess = 1
					err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
					HandleError(w, err1)
					return
				}

				file, _, err := r.FormFile("newFile") //Insert name of file element here
				defer file.Close()
				if err != nil {
					errMsg := "Could not get uploaded file in food_update. Error getting file submission: " + err.Error()
					fmt.Println(errMsg)
					logWriter(errMsg)
					//Assemble data to be sent to inform User
					vd.MessageDisplay = 0
					vd.IsSuccess = 1
					vd.Message = "Could not get uploaded file in food_update. Error getting file submission: " + err.Error()
					err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
					HandleError(w, err1)
					return
				}
				//Form file is gotten in fileUpdate
				//Determine what food this is
				foodChoice := r.FormValue("foodChoice")
				if strings.Contains(strings.ToUpper(foodChoice), strings.ToUpper("hotdog")) {
					foodChoice := r.FormValue("foodChoice")
					foodType := r.FormValue("foodType")
					condimentType := r.FormValue("condimentType")
					caloriesType := r.FormValue("calories")
					theCalories, _ := strconv.Atoi(caloriesType)
					foodName := r.FormValue("foodName")
					userIDInput := r.FormValue("hiddenUserID")
					theUserID, _ := strconv.Atoi(userIDInput)
					foodID := r.FormValue("foodID")
					theFoodID, _ := strconv.Atoi(foodID)
					//Update Photo first in directory and Amazon first
					goodUpdate, hexFileName, photoID := fileUpdate(w, r, strings.ToUpper(foodChoice), theFoodID, theUserID)
					if goodUpdate == true {
						//Create updated photo path for food entries
						theDir := urlFixer(filepath.Join("amazonimages", "pictures", userIDInput,
							strings.ToUpper(foodChoice), hexFileName))
						//Declare values to be filled in later
						aHotDog := Hotdog{
							HotDogType:  foodType,
							Condiment:   condimentType,
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     photoID,
							PhotoSrc:    theDir,
							DateCreated: "",
							DateUpdated: "",
						}

						aMongoHotDog := MongoHotDog{
							HotDogType:  foodType,
							Condiments:  turnFoodArray(condimentType),
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     photoID,
							PhotoSrc:    theDir,
							DateCreated: "",
							DateUpdated: "",
						}

						aHamburger := Hamburger{
							BurgerType:  foodType,
							Condiment:   condimentType,
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     0,
							PhotoSrc:    "",
							DateCreated: "",
							DateUpdated: "",
						}

						aMongoHamburger := MongoHamburger{
							BurgerType:  foodType,
							Condiments:  turnFoodArray(condimentType),
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     0,
							PhotoSrc:    "",
							DateCreated: "",
							DateUpdated: "",
						}

						//Update SQL and Mongo
						goodReturn := sqlUpdateFood(strings.ToUpper(foodChoice), aHamburger, aHotDog, "file_update")
						goodReturnMongo := mongoUpdateFood(strings.ToUpper(foodChoice), aMongoHamburger, aMongoHotDog, "file_update")
						//Check to see if updates were successful
						if goodReturn == true && goodReturnMongo == true {
							succMsg := "Food successfully updated in Mongo for foodID: " + foodID
							fmt.Println(succMsg)
							logWriter(succMsg)
							//Assemble data to be sent to inform User
							vd.MessageDisplay = 0
							vd.Message = succMsg
							vd.IsSuccess = 0
							err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
							HandleError(w, err1)
							return
						} else {
							errMsg := "goodReturn/goodReturnMongo are false in main"
							logWriter(errMsg)
							fmt.Println(errMsg)
							//Assemble data to be sent to inform User
							vd.MessageDisplay = 0
							vd.Message = errMsg
							vd.IsSuccess = 1
							err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
							HandleError(w, err1)
							return
						}
					} else {
						errMsg := "Error updating file in main"
						logWriter(errMsg)
						fmt.Println(errMsg)
						//Assemble data to be sent to inform User
						vd.MessageDisplay = 0
						vd.Message = errMsg
						vd.IsSuccess = 1
						err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
						HandleError(w, err1)
						return
					}
				} else if strings.Contains(strings.ToUpper(foodChoice), strings.ToUpper("hamburger")) {
					foodChoice := r.FormValue("foodChoice")
					foodType := r.FormValue("foodType")
					condimentType := r.FormValue("condimentType")
					caloriesType := r.FormValue("calories")
					theCalories, _ := strconv.Atoi(caloriesType)
					foodName := r.FormValue("foodName")
					userIDInput := r.FormValue("hiddenUserID")
					theUserID, _ := strconv.Atoi(userIDInput)
					foodID := r.FormValue("foodID")
					theFoodID, _ := strconv.Atoi(foodID)
					//Update Photo first in directory and Amazon first
					goodUpdate, hexFileName, photoID := fileUpdate(w, r, strings.ToUpper(foodChoice), theFoodID, theUserID)
					if goodUpdate == true {
						//Create updated photo path for food entries
						theDir := urlFixer(filepath.Join("amazonimages", "pictures", userIDInput,
							strings.ToUpper(foodChoice), hexFileName))
						//Declare values to be filled in later
						aHotDog := Hotdog{
							HotDogType:  foodType,
							Condiment:   condimentType,
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     photoID,
							PhotoSrc:    theDir,
							DateCreated: "",
							DateUpdated: "",
						}

						aMongoHotDog := MongoHotDog{
							HotDogType:  foodType,
							Condiments:  turnFoodArray(condimentType),
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     photoID,
							PhotoSrc:    theDir,
							DateCreated: "",
							DateUpdated: "",
						}

						aHamburger := Hamburger{
							BurgerType:  foodType,
							Condiment:   condimentType,
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     photoID,
							PhotoSrc:    theDir,
							DateCreated: "",
							DateUpdated: "",
						}

						aMongoHamburger := MongoHamburger{
							BurgerType:  foodType,
							Condiments:  turnFoodArray(condimentType),
							Calories:    theCalories,
							Name:        foodName,
							UserID:      theUserID,
							FoodID:      theFoodID,
							PhotoID:     photoID,
							PhotoSrc:    theDir,
							DateCreated: "",
							DateUpdated: "",
						}

						//Update SQL and Mongo
						goodReturn := sqlUpdateFood(strings.ToUpper(foodChoice), aHamburger, aHotDog, "file_update")
						goodReturnMongo := mongoUpdateFood(strings.ToUpper(foodChoice), aMongoHamburger, aMongoHotDog, "file_update")
						//Check to see if updates were successful
						if goodReturn == true && goodReturnMongo == true {
							succMsg := "Food successfully updated in Mongo for foodID: " + foodID
							fmt.Println(succMsg)
							logWriter(succMsg)
							//Assemble data to be sent to inform User
							vd.MessageDisplay = 0
							vd.Message = succMsg
							vd.IsSuccess = 0
							err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
							HandleError(w, err1)
							return
						} else {
							errMsg := "goodReturn/goodReturnMongo are false in main"
							logWriter(errMsg)
							fmt.Println(errMsg)
							//Assemble data to be sent to inform User
							vd.MessageDisplay = 0
							vd.Message = errMsg
							vd.IsSuccess = 1
							err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
							HandleError(w, err1)
							return
						}
					} else {
						errMsg := "Error updating file in main"
						logWriter(errMsg)
						fmt.Println(errMsg)
						//Assemble data to be sent to inform User
						vd.MessageDisplay = 0
						vd.Message = errMsg
						vd.IsSuccess = 1
						err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
						HandleError(w, err1)
						return
					}
				} else {
					errMsg := "Incorrect foodChoice in main: " + foodChoice
					logWriter(errMsg)
					fmt.Println(errMsg)
					//Assemble data to be sent to inform User
					vd.MessageDisplay = 0
					vd.Message = errMsg
					vd.IsSuccess = 1
					err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
					HandleError(w, err1)
					return
				}
			} else if strings.Contains(strings.ToLower(hiddenHasPhoto), "no_photo") {
				fmt.Println("DEBUG: Updating food with no photo")
				//Parse Form
				maxSize := int64(1024000) // allow only 1MB of file size
				err := r.ParseForm()
				if err != nil {
					fmt.Printf("DEBUG: Ignoring this because no file should be submitted %v, %v\n", maxSize,
						err.Error())
				}
				//Determine what food this is (no photo update)
				foodChoice := r.FormValue("foodChoice")
				if strings.Contains(strings.ToUpper(foodChoice), strings.ToUpper("hotdog")) {
					foodChoice := r.FormValue("foodChoice")
					foodType := r.FormValue("foodType")
					condimentType := r.FormValue("condimentType")
					caloriesType := r.FormValue("calories")
					theCalories, _ := strconv.Atoi(caloriesType)
					foodName := r.FormValue("foodName")
					userIDInput := r.FormValue("hiddenUserID")
					theUserID, _ := strconv.Atoi(userIDInput)
					foodID := r.FormValue("foodID")
					theFoodID, _ := strconv.Atoi(foodID)
					//No Photo update, just food
					//Create updated photo path for food entries
					theDir := urlFixer(filepath.Join("amazonimages", "pictures", userIDInput, strings.ToUpper(foodChoice), "hexFileName"))
					//Declare values to be filled in later
					aHotDog := Hotdog{
						HotDogType:  foodType,
						Condiment:   condimentType,
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    theDir,
						DateCreated: "",
						DateUpdated: "",
					}

					aMongoHotDog := MongoHotDog{
						HotDogType:  foodType,
						Condiments:  turnFoodArray(condimentType),
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    theDir,
						DateCreated: "",
						DateUpdated: "",
					}

					aHamburger := Hamburger{
						BurgerType:  foodType,
						Condiment:   condimentType,
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    "",
						DateCreated: "",
						DateUpdated: "",
					}

					aMongoHamburger := MongoHamburger{
						BurgerType:  foodType,
						Condiments:  turnFoodArray(condimentType),
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    "",
						DateCreated: "",
						DateUpdated: "",
					}

					//Update SQL and Mongo
					goodReturn := sqlUpdateFood(strings.ToUpper(foodChoice), aHamburger, aHotDog, "no_photo")
					goodReturnMongo := mongoUpdateFood(strings.ToUpper(foodChoice), aMongoHamburger, aMongoHotDog, "no_photo")
					//Check to see if updates were successful
					if goodReturn == true && goodReturnMongo == true {
						succMsg := "Food successfully updated in Mongo for foodID: " + foodID
						fmt.Println(succMsg)
						logWriter(succMsg)
					} else {
						errMsg := "goodReturn/goodReturnMongo are false in main"
						logWriter(errMsg)
						fmt.Println(errMsg)
					}
				} else if strings.Contains(strings.ToUpper(foodChoice), strings.ToUpper("hamburger")) {
					foodChoice := r.FormValue("foodChoice")
					foodType := r.FormValue("foodType")
					condimentType := r.FormValue("condimentType")
					caloriesType := r.FormValue("calories")
					theCalories, _ := strconv.Atoi(caloriesType)
					foodName := r.FormValue("foodName")
					userIDInput := r.FormValue("hiddenUserID")
					theUserID, _ := strconv.Atoi(userIDInput)
					foodID := r.FormValue("foodID")
					theFoodID, _ := strconv.Atoi(foodID)
					//Update Photo first in directory and Amazon first
					//goodUpdate, hexFileName, photoID := fileUpdate(w, r, strings.ToUpper(foodChoice), theFoodID, theUserID)
					//Create updated photo path for food entries
					theDir := urlFixer(filepath.Join("amazonimages", "pictures", userIDInput, strings.ToUpper(foodChoice), "hexFileName"))
					//Declare values to be filled in later
					aHotDog := Hotdog{
						HotDogType:  foodType,
						Condiment:   condimentType,
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    theDir,
						DateCreated: "",
						DateUpdated: "",
					}

					aMongoHotDog := MongoHotDog{
						HotDogType:  foodType,
						Condiments:  turnFoodArray(condimentType),
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    theDir,
						DateCreated: "",
						DateUpdated: "",
					}

					aHamburger := Hamburger{
						BurgerType:  foodType,
						Condiment:   condimentType,
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    theDir,
						DateCreated: "",
						DateUpdated: "",
					}

					aMongoHamburger := MongoHamburger{
						BurgerType:  foodType,
						Condiments:  turnFoodArray(condimentType),
						Calories:    theCalories,
						Name:        foodName,
						UserID:      theUserID,
						FoodID:      theFoodID,
						PhotoID:     0,
						PhotoSrc:    theDir,
						DateCreated: "",
						DateUpdated: "",
					}

					//Update SQL and Mongo
					goodReturn := sqlUpdateFood(strings.ToUpper(foodChoice), aHamburger, aHotDog, "no_photo")
					goodReturnMongo := mongoUpdateFood(strings.ToUpper(foodChoice), aMongoHamburger, aMongoHotDog, "no_photo")
					//Check to see if updates were successful
					if goodReturn == true && goodReturnMongo == true {
						succMsg := "Food successfully updated in Mongo for foodID: " + foodID
						fmt.Println(succMsg)
						logWriter(succMsg)
					} else {
						errMsg := "goodReturn/goodReturnMongo are false in main"
						logWriter(errMsg)
						fmt.Println(errMsg)
					}
				} else {
					errMsg := "Incorrect foodChoice in main: " + foodChoice
					logWriter(errMsg)
					fmt.Println(errMsg)
				}
			} else {
				errMsg := "Incorrect hiddenHasPhoto in main: " + hiddenHasPhoto
				logWriter(errMsg)
				fmt.Println(errMsg)
			}
		} else if strings.Contains(strings.ToLower(theAction), strings.ToLower("test")) {
			//DEBUG This is for testing:
			maxSize := int64(1024000) // allow only 1MB of file size
			err := r.ParseMultipartForm(maxSize)
			if err != nil {
				fmt.Printf("Image too large. Max Size: %v\n", maxSize)
				log.Println(err)
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.Message = "Image too large."
				vd.IsSuccess = 1
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}

			file, fileHeader, err := r.FormFile("newFile") //Insert name of file element here
			if err != nil {
				errMsg := "DEBUG TEST: Could not get uploaded file. Error getting file submission: " + err.Error()
				fmt.Println(errMsg)
				logWriter(errMsg)
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.Message = errMsg
				vd.IsSuccess = 1
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}
			defer file.Close()

			theText := r.FormValue("DEBUGfakeInput")
			hiddenValue := r.FormValue("DEBUGHIDDIENINPUT")
			fmt.Printf("Here is the Text: %v\nHere is the hiddenText: %v\n", theText, hiddenValue)
			//DEBUG Insert this in our images folder
			theDir, _ := os.Getwd()
			thePath := filepath.Join(theDir, "amazonimages", "pictures", "testfolder")
			os.MkdirAll(thePath, 777)

			//Write file on server
			theFileName := fileHeader.Filename
			f, err := os.OpenFile(theFileName, os.O_WRONLY|os.O_CREATE, 0777)

			if err != nil {
				fmt.Printf("Error opening this file to store on server: %v\n", err.Error())
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.Message = "Error opening this file to store on server:" + err.Error()
				vd.IsSuccess = 1
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}
			io.Copy(f, file)
			f.Close()
			defer file.Close()
			//Move file to folder
			thePath2 := filepath.Join(theDir, "amazonimages", "pictures", "testfolder", theFileName)
			readFile, err := os.Open(theFileName)
			if err != nil {
				fmt.Printf("Error opening this file: %v\n", err.Error())
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.Message = "Error opening this file:" + err.Error()
				vd.IsSuccess = 1
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}
			writeToFile, err := os.Create(thePath2)
			if err != nil {
				fmt.Printf("DEBUG: Error creating writeToFile: \n%v\n", err.Error())
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.Message = "Error creating writeToFile: " + err.Error()
				vd.IsSuccess = 1
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}
			//Move file Contents to folder
			_, err = io.Copy(writeToFile, readFile)
			if err != nil {
				fmt.Printf("Error copying the contents of the one image to the other.\n%v\n", err.Error())
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.Message = "Error copying the contents of the one image to the other " + err.Error()
				vd.IsSuccess = 1
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}
			readFile.Close()    //Close File
			writeToFile.Close() //Close File
			//Delete created file
			removeErr := os.Remove(theFileName)
			if removeErr != nil {
				fmt.Printf("Error removing the file: %v\n", removeErr.Error())
				//Assemble data to be sent to inform User
				vd.MessageDisplay = 0
				vd.Message = "Error removing the file: " + err.Error()
				vd.IsSuccess = 1
				err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
				HandleError(w, err1)
				return
			}
		} else {
			errMsg := "Error, incorrect theAction: " + theAction
			logWriter(errMsg)
			fmt.Println(errMsg)
			//Assemble data to be sent to inform User
			vd.MessageDisplay = 0
			vd.Message = errMsg
			vd.IsSuccess = 1
			err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
			HandleError(w, err1)
			return
		}
	} else {
		//Serve the mainpage normally
		err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
		HandleError(w, err1)
		return
	}
}

//Handles the documentation page
func documentation(w http.ResponseWriter, req *http.Request) {
	thePort := os.Getenv("PORT")
	if thePort == "" {
		thePort = "80"
	}

	err1 := template1.ExecuteTemplate(w, "documentation.gohtml", nil)
	HandleError(w, err1)
}

//Handles the documentation page
func contact(w http.ResponseWriter, r *http.Request) {
	thePort := os.Getenv("PORT")
	if thePort == "" {
		thePort = "80"
		fmt.Printf("DEBUG: Defaulting to this port %v\n", thePort)
	}

	if r.Method == http.MethodPost {
		//Handle the email Ajax sent to us
		fmt.Printf("DEBUG: AN EMAIL IS BEING SENT TO ME.\n")
		//Collect JSON from Postman or wherever
		//Get the byte slice from the request body ajax
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		//Marshal the user data into our type
		var dataPosted UserJSON
		json.Unmarshal(bs, &dataPosted)

		successEmail := emailToMe(dataPosted)

		if successEmail == true {
			//Send successful response back
			type successMSG struct {
				Message     string `json:"Message"`
				SuccessNum  int    `json:"SuccessNum"`
				RedirectURL string `json:"RedirectURL"`
			}
			msgSuccess := successMSG{
				Message:     "Added the new account!",
				SuccessNum:  0,
				RedirectURL: "http://" + serverAddress,
			}
			theJSONMessage, err := json.Marshal(msgSuccess)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprint(w, string(theJSONMessage))
		} else {
			type successMSG struct {
				Message     string `json:"Message"`
				SuccessNum  int    `json:"SuccessNum"`
				RedirectURL string `json:"RedirectURL"`
			}
			msgSuccess := successMSG{
				Message:     "Added the new account!",
				SuccessNum:  1,
				RedirectURL: "http://" + serverAddress,
			}

			theJSONMessage, err := json.Marshal(msgSuccess)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprint(w, string(theJSONMessage))
		}

	} else {
		//Serve the template normally
		err1 := template1.ExecuteTemplate(w, "contact.gohtml", nil)
		HandleError(w, err1)
	}
}

//Handles the messageboard page
func messageboard(w http.ResponseWriter, r *http.Request) {
	thePort := os.Getenv("PORT")
	if thePort == "" {
		thePort = "80"
	}

	err1 := template1.ExecuteTemplate(w, "messageboard.gohtml", nil)
	HandleError(w, err1)
}

//choicepage (in-between stop between either the messageboard or the database)
func choicepage(w http.ResponseWriter, r *http.Request) {
	//if User is already logged in, bring them to the mainPage!
	aUser := getUser(w, r) //Get the User, if they exist
	aUserRole := aUser.Role
	thePort := os.Getenv("PORT")
	if thePort == "" {
		thePort = "80"
		logWriter("Defautling to this port: " + thePort)
	}
	vd := ViewData{aUser, aUser.UserName, aUserRole, thePort, 1, "Welcome to the Main page!", 0}
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//See if there is a submission for new food or updates/deletes
	//Serve the mainpage normally
	err1 := template1.ExecuteTemplate(w, "choicepage.gohtml", vd)
	HandleError(w, err1)
}

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/signup", signup)
	myRouter.HandleFunc("/mainPage", mainPage)
	myRouter.HandleFunc("/signUpUserUpdated", signUpUserUpdated)
	myRouter.HandleFunc("/documentation", documentation)
	myRouter.HandleFunc("/contact", contact)
	myRouter.HandleFunc("/choicepage", choicepage)
	myRouter.HandleFunc("/messageboard", messageboard)
	//SQL Database Stuff
	myRouter.HandleFunc("/deleteFood", deleteFood).Methods("POST")
	//myRouter.HandleFunc("/updateFood", updateFood).Methods("POST")           //Update a certain food item
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
	myRouter.HandleFunc("/checkSRC", checkSRC).Methods("POST")                   //Check if directory exists
	myRouter.HandleFunc("/deletePhotoFromS3", deletePhotoFromS3).Methods("POST") //Delete S3 Photo
	//Validation Stuff
	myRouter.HandleFunc("/checkUsername", checkUsername) //Check Username
	myRouter.HandleFunc("/loadUsernames", loadUsernames) //Loads in Usernames
	//API Checking Stuff
	myRouter.HandleFunc("/userInfoAPI", userInfoAPI).Methods("POST") //Get food information for User
	//Middleware logging
	myRouter.Handle("/", loggingMiddleware(http.HandlerFunc(logHandler)))
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//Serve the Amazon files
	myRouter.Handle("/", http.FileServer(http.Dir("./amazonimages")))
	myRouter.PathPrefix("/amazonimages/").Handler(http.StripPrefix("/amazonimages/", http.FileServer(http.Dir("./amazonimages"))))
	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed
	//Write initial entry to log
	//logWriter("Deployed superUltimateDBWebApp app.")
	//open SQL connection
	db, err = sql.Open("mysql",
		dbConnectString)
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	//Mongo Connect
	mongoClient = connectDB()
	/* Do below so our map dosen't go crazy... */
	loadedMessagesMap = make(map[int]Message)
	createTestMessages() //Create test board for messages

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

//USER/Developer API Ping
func userInfoAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: Someone posted to the userInfoAPI.\n")
	type UserSearch struct {
		UserID int `json:"UserID"`
	} //Declare incoming JSON

	type InfoResponse struct {
		ErrMsg []string `json:"ErrMsg"`
		User   AUser    `json:"User"`
	} //Define outgoing JSON

	errMsg := "" //Define the err message to return as part of the ErrMsg array in infoResponse

	responseMsg := InfoResponse{} //Declare response JSON to send back

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errMsg = "Error reading the request coming in: " + err.Error() //Assemble Error
		responseMsg.ErrMsg = append(responseMsg.ErrMsg, errMsg)        //Append Error data
		fmt.Println(errMsg)
		dataJSON, err := json.Marshal(responseMsg)
		if err != nil {
			fmt.Println("Sending response JSON in userInfoAPI: " + err.Error())
			fmt.Fprintf(w, string("Sending response JSON in userInfoAPI: "+err.Error())) //Send Error back
			return
		} else {
			fmt.Fprintf(w, string(dataJSON)) //Send Error back
			return
		}
	}
	//Marshal it into our type
	var findUser UserSearch
	json.Unmarshal(bs, &findUser)
	fmt.Printf("DEBUG: Here is our findUser: %v\n", findUser)
	//Query Mongo to get User to return
	returnedUser, successFail, message := getUserMongo(findUser.UserID)
	if successFail == true {
		errMsg = message                                         //Initialze Message to add
		responseMsg.ErrMsg = append(responseMsg.ErrMsg, message) //Append Message data
		responseMsg.User = returnedUser                          //Add User to the returned JSON
		//Marshal JSON to return
		dataJSON, err := json.Marshal(responseMsg)
		if err != nil {
			fmt.Fprintf(w, string("Sending response JSON in userInfoAPI: "+err.Error())) //Send Error back
			return
		} else {
			fmt.Fprintf(w, string(dataJSON)) //Send Message back
			return
		}
	} else {
		//Mongo query failure, sending JSON back to User
		errMsg = "Error finding User in Mongo: " + message       //Assemble Error
		responseMsg.ErrMsg = append(responseMsg.ErrMsg, message) //Append Error data
		fmt.Println(errMsg)
		dataJSON, err := json.Marshal(responseMsg)
		if err != nil {
			fmt.Fprintf(w, string("Sending response JSON in userInfoAPI: "+err.Error())) //Send Error back
			return
		} else {
			fmt.Fprintf(w, string(dataJSON)) //Send Error back
			return
		}
	}
}

/* DEBUG ZONE */
