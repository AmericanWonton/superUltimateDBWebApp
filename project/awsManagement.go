package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/mgo.v2/bson"
)

//GLOBAL VARIABLES FOR FILE INSERTION
var awsuserID int = 0
var awsfoodID int = 0
var awsphotoName string = ""
var awsfoodType string = ""
var awsdateCreated string = ""
var awsdateUpdated string = ""

type PhotoInsert struct {
	UserID      int    `json:"UserID"`
	FoodID      int    `json:"FoodID"`
	PhotoID     int    `json:"PhotoID"`
	PhotoName   string `json:"PhotoName"`
	FileType    string `json:"FileType"`
	Size        int64  `json:"Size"`
	PhotoHash   string `json"PhotoHash"`
	Link        string `json:"Link"`
	FoodType    string `json:"FoodType"`
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

//Take a file form for submission
func fileInsert(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: Sweetass, you submitted a cooler file.\n")
	//Specify Form side and parse form
	maxSize := int64(1024000) // allow only 1MB of file size
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		fmt.Printf("Image too large. Max Size: %v\n", maxSize)
		log.Println(err)
		return
	}

	hiddenUserNum := r.FormValue("hiddenUserNum")
	hiddenFoodType := r.FormValue("hiddenFoodType")
	hiddenFoodNum := r.FormValue("hiddenFoodType")

	file, fileHeader, err := r.FormFile("newFile2") //Insert name of file element here
	if err != nil {
		fmt.Printf("Could not get uploaded file. Error getting file submission: %v\n", err.Error())
		log.Println(err)
		return
	}
	//Defer closing the file
	fmt.Printf("DEBUG: Writing the temporary file.\n")
	//Make Directory for writing file
	fmt.Printf("Here is just the fileheader: %v\n and here is the extension: %v\n",
		filepath.Ext(fileHeader.Filename), fileHeader.Filename)
	hexName := bson.NewObjectId().Hex()
	fileExtension := filepath.Ext(fileHeader.Filename)
	theFileName := hexName + fileExtension
	theDir, _ := os.Getwd()
	thePath := filepath.Join(theDir, "static", "images", "pictures", "31553001", "HOTDOGS")
	os.MkdirAll(thePath, os.ModePerm)
	//Write file on server
	f, err := os.OpenFile(theFileName, os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		fmt.Printf("Error opening this file to store on server.\n")
	}
	io.Copy(f, file)
	f.Close()
	file.Close()
	//Move file to folder
	thePath2 := filepath.Join(theDir, "static", "images", "pictures", "31553001", "HOTDOGS", theFileName)
	readFile, err := os.Open(theFileName)
	if err != nil {
		fmt.Printf("Error opening this file: %v\n", err.Error())
	}
	writeToFile, err := os.Create(thePath2)
	if err != nil {
		fmt.Printf("DEBUG: Error creating writeToFile: \n%v\n", err.Error())
	}
	//Move file Contents to folder
	n, err := io.Copy(writeToFile, readFile)
	if err != nil {
		fmt.Printf("Error copying the contents of the one image to the other.\n%v\n", err.Error())
	}
	fmt.Printf("DEBUG: move the contents of n: %v\n", n)
	readFile.Close()    //Close File
	writeToFile.Close() //Close File
	//Delete created file
	removeErr := os.Remove(theFileName)
	if removeErr != nil {
		fmt.Printf("Error removing the file: %v\n", removeErr.Error())
	}
	/********************* SEND FILES TO AMAZON BUCKET **********************/
	fmt.Printf("DEBUG: Sending files to Amazon.\n")
	//Give a hex for the file value
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(
			"AKIAV7U75A66QIWGBCVK",                     // id
			"BPwAJceAAI3xttbpaB6Dshdx0uB6ugvrknSXD7zi", // secret
			""), // token can be left blank for now
	})
	if err != nil {
		fmt.Printf("Could not upload file. Error creating session.\n")
	}
	fileName, err := UploadFileToS3(hexName, s, file, fileHeader)
	if err != nil {
		fmt.Printf("Could not upload file...issue uploading to Amazon.\n")
	} else {
		fmt.Printf("Image uploaded successfully to Amazon Bucket: %v\n", fileName) //print success
	}
	/********************* UPLOAD PHOTO DETAILS TO THE DATABASE **********************/
	fmt.Printf("DEBUG: Uploading file to database.\n")
	//Upload photo details to DB
	extension := filepath.Ext(fileHeader.Filename)
	fileURL := "pictures/31553001/HOTDOGS/" + hexName + extension
	insertedPhoto := insertUserPhotos("HOTDOGS", fileURL, 31553001, fileHeader.Filename, extension, fileHeader.Size,
		hexName)
	if insertedPhoto == true {
		fmt.Println("DEBUG: Inserted photo information into SQL DB.")
	} else {
		fmt.Println("DEBUG: Issue inserting photo information into SQL DB.")
	}
	//err1 := template1.ExecuteTemplate(w, "index.gohtml", someData)
	//HandleError(w, err1)
	http.Redirect(w, r, "/", 302)
}

// UploadFileToS3 saves a file to aws bucket and returns the url to // the file and an error if there's any
func UploadFileToS3(aHex string, s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error, bool) {
	// get the file size and read
	successfulFileSend := true
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	stringUserID := strconv.Itoa(awsuserID)
	tempFileName := "pictures/" + stringUserID + "/" + awsfoodType + "/" + aHex + filepath.Ext(fileHeader.Filename)

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("webapppics"),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		fmt.Printf("Error submitting file for Amazon bucket in UploadFileToS3: \n%v\n", err.Error())
		successfulFileSend = false
		return "", err, successfulFileSend
	}

	return tempFileName, err, successfulFileSend
}

func fileInsertDBS(foodtype string, fileURL string, userID int, foodID int, photoID int,
	theFileName string, theExtension string, theSize int64, theHex string) bool {

	successfulInsert := true
	fmt.Printf("DEBUG: Inserting photos into SQL.\n")
	theTimeNow := time.Now()
	//Which Type of food?
	if strings.Contains(foodtype, "HOTDOG") {
		fmt.Printf("Inserting Hotdog Photo\n")
		theStatement := "INSERT INTO user_photos" +
			"(USER_ID, FOOD_ID, PHOTO_ID, PHOTO_NAME, FILE_TYPE, SIZE, PHOTO_HASH, LINK, FOOD_TYPE, DATE_CREATED, DATE_UPDATED) " +
			"VALUES(?,?,?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(theStatement)

		r, err := stmt.Exec(userID, foodID, photoID, theFileName, theExtension, theSize, theHex,
			fileURL, "HOTDOG", theTimeNow.Format("2006-01-02 15:04:05"),
			theTimeNow.Format("2006-01-02 15:04:05"))
		check(err)

		n, err := r.RowsAffected()
		check(err)
		fmt.Printf("%v rows effected.\n", n)
		stmt.Close() //Close the SQL

		//INSERT INTO MongoDB
		photoInsertion := PhotoInsert{
			UserID:      userID,
			FoodID:      foodID,
			PhotoID:     photoID,
			PhotoName:   theFileName,
			FileType:    theExtension,
			Size:        theSize,
			PhotoHash:   theHex,
			Link:        fileURL,
			FoodType:    "HOTDOG",
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Collect Data for Mongo
		photoCollection := mongoClient.Database("superdbtest1").Collection("photos") //Here's our collection
		collectedUsers := []interface{}{photoInsertion}
		//Insert Our Data
		insertManyResult, err2 := photoCollection.InsertMany(theContext, collectedUsers)
		if err2 != nil {
			fmt.Printf("We had an error inserting a photo into MongoSQL: %v\n", err2.Error())
		} else {
			fmt.Printf("Inserting photo for hotdogs was successful: %v\n", insertManyResult)
		}
	} else if strings.Contains(foodtype, "HAMBURGER") {
		fmt.Printf("Inserting Hamburger Photo into SQL\n")
		theStatement := "INSERT INTO user_photos" +
			"(USER_ID, FOOD_ID, PHOTO_ID, PHOTO_NAME, FILE_TYPE, SIZE, PHOTO_HASH, LINK, FOOD_TYPE, DATE_CREATED, DATE_UPDATED) " +
			"VALUES(?,?,?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(theStatement)

		r, err := stmt.Exec(userID, foodID, photoID, theFileName, theExtension, theSize, theHex,
			fileURL, "HAMBURGER", theTimeNow.Format("2006-01-02 15:04:05"),
			theTimeNow.Format("2006-01-02 15:04:05"))
		check(err)

		n, err := r.RowsAffected()
		check(err)
		fmt.Printf("%v rows effected.\n", n)
		stmt.Close() //Close the SQL

		//INSERT INTO MongoDB
		photoInsertion := PhotoInsert{
			UserID:      userID,
			FoodID:      foodID,
			PhotoID:     photoID,
			PhotoName:   theFileName,
			FileType:    theExtension,
			Size:        theSize,
			PhotoHash:   theHex,
			Link:        fileURL,
			FoodType:    "HAMBURGER",
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Collect Data for Mongo
		photoCollection := mongoClient.Database("superdbtest1").Collection("photos") //Here's our collection
		collectedUsers := []interface{}{photoInsertion}
		//Insert Our Data
		insertManyResult, err2 := photoCollection.InsertMany(theContext, collectedUsers)
		if err2 != nil {
			fmt.Printf("We had an error inserting a photo into MongoSQL: %v\n", err2.Error())
			successfulInsert = false
		} else {
			fmt.Printf("Inserting photo for hamburgers was successful: %v\n", insertManyResult)
		}
	} else {
		fmt.Printf("Did not recieve the correct foodtype: %v\n", foodtype)
		successfulInsert = false
	}

	return successfulInsert
}
