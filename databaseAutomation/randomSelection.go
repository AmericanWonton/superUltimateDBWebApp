package main

import (
	"encoding/hex"
	"math/rand"
	"time"
)

type RandomPerson struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Street struct {
				Number int    `json:"number"`
				Name   string `json:"name"`
			} `json:"street"`
			City        string `json:"city"`
			State       string `json:"state"`
			Country     string `json:"country"`
			Postcode    int    `json:"postcode"`
			Coordinates struct {
				Latitude  string `json:"latitude"`
				Longitude string `json:"longitude"`
			} `json:"coordinates"`
			Timezone struct {
				Offset      string `json:"offset"`
				Description string `json:"description"`
			} `json:"timezone"`
		} `json:"location"`
		Email string `json:"email"`
		Login struct {
			UUID     string `json:"uuid"`
			Username string `json:"username"`
			Password string `json:"password"`
			Salt     string `json:"salt"`
			Md5      string `json:"md5"`
			Sha1     string `json:"sha1"`
			Sha256   string `json:"sha256"`
		} `json:"login"`
		Dob struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"dob"`
		Registered struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"registered"`
		Phone string `json:"phone"`
		Cell  string `json:"cell"`
		ID    struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"id"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
		Nat string `json:"nat"`
	} `json:"results"`
	Info struct {
		Seed    string `json:"seed"`
		Results int    `json:"results"`
		Page    int    `json:"page"`
		Version string `json:"version"`
	} `json:"info"`
}

//Returns a hashed string of bytes for our User struct
func randomPassword() string {
	//slice of string for password generation
	somePasswords := []string{"bigBrownFox69", "sirPants22", "elvishmanhood1", "sirDankulus",
		"bigTudeDude", "seniorSeniorson", "thebottleO34"}
	useridmin, useridmax := 0, 6                                                //The min and Max value for our randInt
	bsString := []byte(somePasswords[rand.Intn(useridmax-useridmin)+useridmin]) //Encode Password
	encodedString := hex.EncodeToString(bsString)                               //Encode Password Pt2

	return encodedString
}

func randomUsername(fName string, lName string) string {
	ourUsername := ""
	ourUsername = ourUsername + string(fName[0:3]) + string(lName[0:3]) + "2020"
	return ourUsername
}
