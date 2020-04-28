package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//Here's our User struct
type User struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

type Hotdog struct {
	User       User //The User whomst this hotdog belongs to
	HotDogType string
	Condiment  string
	Calories   int
	Name       string
}

type Hamburger struct {
	User       User //The User whomst this hotdog belongs to
	BurgerType string
	Condiment  string
	Calories   int
	Name       string
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
			err := bcrypt.CompareHashAndPassword(loginUser.Password, []byte(password))
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
		/* We should probobly due some field validation with ajax and mongo... */
		if _, ok := dbUsers[username]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
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
		theUser = User{username, bs, firstname, lastname, role}
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

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/signup", signUp)
	myRouter.HandleFunc("/mainPage", mainPage)
	/* Handle all files in the static path */
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
