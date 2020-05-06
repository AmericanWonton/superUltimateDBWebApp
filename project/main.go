package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//Here's our User struct
type User struct {
	UserName string
	Password string //This was formally a []byte but we are changing our code to fit the database better
	First    string
	Last     string
	Role     string
	UserID   int
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
	User       User   `json:"User"` //The User whomst this hotdog belongs to
	BurgerType string `json:"BurgerType"`
	Condiment  string `json:"Condiment"`
	Calories   int    `json:"Calories"`
	Name       string `json:"Name"`
	UserID     int    `json:"UserID"`
}

//Here is our ViewData struct
type ViewData struct {
	User     User
	UserName string
}

//Here's our session struct
type session struct {
	username     string
	lastActivity time.Time
}

//Session Database info
var dbUsers = map[string]User{}       // user ID, user
var dbSessions = map[string]session{} // session ID, session
var dbSessionsCleaned time.Time

//mySQL database declarations
var db *sql.DB
var err error

const sessionLength int = 30 //Length of sessions

/* TEMPLATE DEFINITION BEGINNING */
var template1 *template.Template

/* FUNCMAP DEFINITION */
func (u User) ReturnRoleUser(theUser string) bool {
	if strings.Compare(theUser, "user") == 0 {
		return true
	} else {
		return false
	}
}

func (u User) ReturnRoleAdmin(theAdmin string) bool {
	if strings.Compare(theAdmin, "admin") == 0 {
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
	//template1 = template.Must(template.ParseGlob("templates/*"))
	template1 = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*gohtml"))
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
	//if User is already logged in, bring them to the mainPage!
	/*
		aUser := getUser(w, r) //Get the User, if they exist
		if alreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/mainPage", http.StatusSeeOther)
			return
		}
	*/
	//If a User posts a form to log in!
	if r.Method == http.MethodPost {
		//Get Form Values
		username := r.FormValue("username")
		password := r.FormValue("password")
		//Search for Users in our Database. It fails out if Username and Password aren't there.
		if loginUser, ok := dbUsers[username]; ok {
			fmt.Printf("We found the Username %v\n", username)
			//Check on Password
			err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(password))
			if err != nil {
				http.Error(w, "Username and/or password do not match", http.StatusForbidden)
				return
			}
			fmt.Printf("We found the password, %v, updating session. \n", password)
			//User logged in, directing them to the mainpage
			// create session
			sID, _ := uuid.NewV4()
			cookie := &http.Cookie{
				Name:  "session",
				Value: sID.String(),
			}
			cookie.MaxAge = sessionLength
			http.SetCookie(w, cookie)
			dbSessions[cookie.Value] = session{username, time.Now()}
			//Send to the MainPage!
			fmt.Printf("Executing the main page now with our logged in User!\n")
			http.Redirect(w, r, "/mainPage", http.StatusSeeOther)
			return
		}
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "index.gohtml", aUser)
	HandleError(w, err1)
	fmt.Printf("Homepage Endpoint Hit\n")
	getHotDogsAll(w, r)
}

//signUp
func signUp(w http.ResponseWriter, req *http.Request) {
	//See if user is already logged in
	if alreadyLoggedIn(w, req) {
		//If already logged in, put them back at the main menu
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	//If User is NOT already logged in, wait till they post a form!
	var theUser User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		role := req.FormValue("role")
		// username taken?
		/* We have field validation with Ajax...do we need this?
		if _, ok := dbUsers[username]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		*/
		// create session
		sID, _ := uuid.NewV4()
		newCookie := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		newCookie.MaxAge = sessionLength
		http.SetCookie(w, newCookie)
		dbSessions[newCookie.Value] = session{username, time.Now()}
		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		//Make User and USERID
		goodNum := false
		theID := 0
		row, err := db.Query(`SELECT user_id FROM users;`)
		check(err)
		defer row.Close()

		for goodNum == false {
			//Build the random, unique integer to be assigned to this User
			goodNumFound := true //A second checker to break this loop
			randInt := 0         //The random integer added onto ID
			var databaseID int   //The ID returned from the database while searching
			randIntString := ""  //The integer built through a string...
			min, max := 0, 9     //The min and Max value for our randInt
			for i := 0; i < 8; i++ {
				randInt = rand.Intn(max-min) + min
				randIntString = randIntString + strconv.Itoa(randInt)
			}
			theID, err = strconv.Atoi(randIntString)
			if err != nil {
				fmt.Println(err)
				return
			}
			//Check to see if the built number is taken.
			for row.Next() {
				err = row.Scan(&databaseID)
				check(err)
				if databaseID == theID {
					//Found the number, need to create another one!
					fmt.Printf("Found the ID, %v, in the database: %v. Creating another one...\n",
						theID, databaseID)
					goodNumFound = false
					break
				} else {

				}
			}
			//Final check to see if we need to go through this loop again
			if goodNumFound == false {
				goodNum = false
			} else {
				goodNum = true
			}
		}

		//Add User to the SQL Database
		ourInsertStatement := `INSERT INTO users VALUES ("` + username + `", "` + password + `", "` + firstname + `", "` +
			lastname + `", "` + role + `", "` + strconv.Itoa(theID) + `");`
		stmt, err := db.Prepare(ourInsertStatement)
		check(err)
		defer stmt.Close()

		r, err := stmt.Exec()
		check(err)

		n, err := r.RowsAffected()
		check(err)

		fmt.Fprintln(w, "INSERTED RECORD", n)

		theUser = User{username, string(bs), firstname, lastname, role, theID}
		dbUsers[username] = theUser
		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	} else {
		err1 := template1.ExecuteTemplate(w, "signup.gohtml", nil)
		HandleError(w, err1)
	}

	fmt.Printf("Signup Endpoint Hit\n")
}

//mainPage
func mainPage(w http.ResponseWriter, req *http.Request) {
	//if User is already logged in, bring them to the mainPage!
	aUser := getUser(w, req)              //Get the User, if they exist
	vd := ViewData{aUser, aUser.UserName} //POSSIBLY DEBUG
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	/* Execute template, handle error */
	/* DEGUG STUFF */
	fmt.Println("Is this a problem area?....")
	/*
		err1 :=  template.Must(template1.Clone()).Funcs(template1.FuncMap{
			"is"
		})
	*/
	err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
	HandleError(w, err1)
	fmt.Printf("Homepage Endpoint Hit\n")
}

//POST mainpage
func postHotDog(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inserting hotdog record.")
	//Collect JSON from Postman or wherever
	reqBody, _ := ioutil.ReadAll(req.Body)
	//Marshal it into our type
	var postedHotDog Hotdog
	json.Unmarshal(reqBody, &postedHotDog)
	//Debug
	fmt.Printf("Here is our hotdog: \n%v\n", postedHotDog)
	/*
		sqlStatement := "INSERT INTO hot_dogs (TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES" +
			"(" + "\"" + postedHotDog.HotDogType + "\", \"" + postedHotDog.Condiment + "\", \"" +
			string(postedHotDog.Calories) + "\", \"" + postedHotDog.Name + "\", \"" + string(postedHotDog.UserID) +
			"\");"
		fmt.Printf("Here is our sql statement:\n\n%v\n", sqlStatement)
	*/
	stmt, err := db.Prepare("INSERT INTO hot_dogs (TYPE, CONDIMENT, CALORIES, NAME, USER_ID) VALUES (:TYPE, :CONDIMENT, :CALORIES, :NAME, :USER_ID)")
	defer stmt.Close()

	r, err := stmt.Exec(sql.Named("TYPE", postedHotDog.HotDogType), sql.Named("CONDIMENT", postedHotDog.Condiment),
		sql.Named("CALORIES", postedHotDog.Calories), sql.Named("NAME", postedHotDog.Name), sql.Named("USER_ID", postedHotDog.UserID))
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Fprintln(w, "INSERTED RECORD", n)
}

//GET mainpage
func getHotDogSingular(w http.ResponseWriter, req *http.Request) {
	//Get the string map of our variables from the request
	fmt.Println("Finding hotdog singular")
	//Collect JSON from Postman or wherever
	reqBody, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("Here's our body: \n%v\n", reqBody)
	//Marshal it into our type
	var postedHotDog Hotdog
	json.Unmarshal(reqBody, &postedHotDog)
	fmt.Printf("Here is our postedHotDog: %v\n", postedHotDog)

	rows, err := db.Query(`SELECT * FROM hot_dogs WHERE NAME = 'HOT_AND_READY';`)
	check(err)
	defer rows.Close()
	var id int64
	var theUser string
	var dogType string
	var condiment string
	var calories int
	var hotdogName string
	var userID string
	count := 0
	for rows.Next() {
		err = rows.Scan(&id, &theUser, &dogType, &condiment, &calories, &hotdogName, &userID)
		check(err)
		fmt.Printf("Retrieved Record: %v\n", hotdogName)
		count++
	}
	//If nothing returned from the rows
	if count == 0 {
		fmt.Printf("Nothing returned for this query.\n")
		return
	} else {
		//Assign the returned name to our object
		fmt.Printf("Hotdog name is: %v\n", hotdogName)
		//Compare to see if the name matches the name we posted
		if strings.Compare(postedHotDog.Name, hotdogName) == 0 {
			fmt.Printf("Hey, our query %v matches our posted JSON, %v \n", hotdogName, postedHotDog.Name)
		} else {
			fmt.Printf("Whooops, our query, %v, does not match our JSON, %v\n", hotdogName, postedHotDog.Name)
		}
	}
	//DEBUG, need to see how rows are returned.
	fmt.Printf("Here is our rows returned:\nID:%v\nTheUser:%v\nDog Type:%v\nCondiment:%v\nCalories:%v\nHotdogname:%v\nuserID: %v\n",
		id, theUser, dogType, condiment, calories, hotdogName, userID)
}

func getHotDogsAll(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query(`SELECT TYPE FROM hot_dogs;`)
	check(err)
	defer rows.Close()
	// data to be used in query
	var s, name string
	s = "RETRIEVED RECORDS:\n"

	// query
	/* From the documentation, Next returns the next row in the line of rows we asked for from the 'rows' variable above.
	It returns false if there's no row up next, (so basically, it's really good for loops) */
	for rows.Next() {
		/* Scan copies the columns in the current row and copies them to a destination. So we set the destination,
		(that 'name' string variable above), and point it to that */
		err = rows.Scan(&name)
		check(err)       //Check to make sure there was no error doing that above.
		s += name + "\n" //We keep adding the name returned and a newline for printing later.
	}

	fmt.Printf("Here's the records, fucker: \n%v\n", s)
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/signup", signUp)
	myRouter.HandleFunc("/mainPage", mainPage)
	//Database Stuff
	myRouter.HandleFunc("/postHotDog", postHotDog).Methods("POST")            //Post a hotdog!
	myRouter.HandleFunc("/scadoop", getHotDogsAll).Methods("GET")             //Get ALL Hotdogs!
	myRouter.HandleFunc("/getHDogSingular", getHotDogSingular).Methods("GET") //Get a SINGULAR hotdog
	//Validation Stuff
	myRouter.HandleFunc("/checkUsername", checkUsername) //Check Username
	myRouter.HandleFunc("/loadUsernames", loadUsernames) //Loads in Usernames
	//Serve our CSS files...
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("."+"/static/"))))

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

	/* DEBUG
	elHotDog := Hotdog{
		User{"butthole", []byte("dingus"), "First", "Last", "Role", 38298457},
		"dogtype",
		"condiment",
		650,
		"Name",
		38298457,
	}

	q, err := json.Marshal(elHotDog)
	if err != nil {
		fmt.Println("There's an error marshalling.")
	}
	fmt.Printf("Here's our JSON: %v\n", string(q))
	*/
	/*
		quotes := "quatation marks"
		bigPeener := "Here's my\"" + quotes + "\""
		fmt.Println(bigPeener)
	*/
	/* DEBUG
	bs, err := bcrypt.GenerateFromPassword([]byte("pWord2"), bcrypt.MinCost)
	if err != nil {
		return
	}
	fmt.Printf("Our hashed password is: %v\n", bs)

	err2 := bcrypt.CompareHashAndPassword(bs, []byte("pWord2"))
	if err != nil {
		return
	}
	fmt.Printf("Err2 is %v\n", err2)

	theBSString := string(bs)
	fmt.Printf("Here's our byte array as a string:\n%v\n", theBSString)
	theStringBS := []byte(theBSString)
	fmt.Printf("Here's our string BS back to a BS: \n%v\n", theStringBS)
	*/

	//Handle Requests
	handleRequests()
}

//Check errors in our mySQL errors
func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
