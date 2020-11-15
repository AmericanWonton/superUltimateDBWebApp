package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var GmailService *gmail.Service //This gets initialized in init

var theClientID string
var theClientSecret string
var theAccessToken string
var theRefreshToken string

type UserJSON struct {
	TheName    string `json:"TheName"`
	TheEmail   string `json:"TheEmail"`
	TheMessage string `json:"TheMessage"`
}

//Initialized at begininning of program
func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     theClientID,
		ClientSecret: theClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://" + serverAddress,
	}

	token := oauth2.Token{
		AccessToken:  theAccessToken,
		RefreshToken: theRefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		errMsg := "Unable to retrieve Gmail client: " + err.Error()
		fmt.Println(errMsg)
		logWriter(errMsg)
	}

	GmailService = srv
	if GmailService != nil {
		succMsg := "Email service is initialized"
		logWriter(succMsg)
	}
}

//Attempts to send an email to User
func signUpUserEmail(theEmail string, theRole string, fName string, lName string) bool {
	goodEmailSend := true
	theMessage := "Hello " + fName + " " + lName + ", thank you for visiting my site!\n" +
		"Please enjoy the new " + theRole + " role you created. Depending on that role, you might " +
		"be limited to certain features. Be respectable and have fun!"
	theSubject := "Welcome to SuperDBTester3000"

	var message gmail.Message

	emailTo := "To: " + theEmail + "\r\n"
	subject := "Subject: " + theSubject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + theMessage)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		errMsg := "Error sending this message to the User: " + err.Error()
		fmt.Println(errMsg)
		logWriter(errMsg)
		goodEmailSend = false
	}

	return goodEmailSend
}

func emailToMe(theEmail UserJSON) bool {
	goodEmailSend := true
	theMessage := theEmail.TheName + " has a message for you for SuperDBWebApp.\n" +
		theEmail.TheMessage + "\nYou can contact him here: " + theEmail.TheEmail
	theSubject := "SuperDBTester3000 questions"

	var message gmail.Message

	emailTo := "To: " + senderAddress + "\r\n"
	subject := "Subject: " + theSubject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + theMessage)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		errMsg := "Error sending this message to the User: " + err.Error()
		fmt.Println(errMsg)
		logWriter(errMsg)
		goodEmailSend = false
	}

	return goodEmailSend
}
