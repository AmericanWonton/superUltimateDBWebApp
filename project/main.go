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
type user struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

//Here's our session struct
type session struct {
	username     string
	lastActivity time.Time
}

//Session Database info
var dbUsers = map[string]user{}       // user ID, user
var dbSessions = map[string]session{} // session ID, session
var dbSessionsCleaned time.Time

const sessionLength int = 30 //Length of sessions

/* TEMPLATE DEFINITION BEGINNING */
var template1 *template.Template

/* FUNCMAP DEFINITION */
var funcMap = template.FuncMap{
	"upperCase": strings.ToUpper, //upperCase is a key we can call inside of the template html file
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
	//if User is already logged in, bring them to the mainPage!
	aUser := getUser(w, r) //Get the User, if they exist
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/mainPage", http.StatusSeeOther)
		return
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
	var theUser user
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
		theUser = user{username, bs, firstname, lastname, role}
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
	aUser := getUser(w, req) //Get the User, if they exist
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/homePage", http.StatusSeeOther)
		return
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "index.gohtml", aUser)
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
