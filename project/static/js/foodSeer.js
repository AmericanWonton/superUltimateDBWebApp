var userID;
function getUserID(passedID){
    userID = passedID;
    console.log("We've set userID to " + userID);
}

var seeFoodButton;
var foodSeerDiv;
var foodUpdater;
var hdogSeerDiv;
var hamSeerDiv;
var foodUpdateBtn;
var foodDeleteBtn;
var hotdogIDArray = new Array(); //Declare Array for id reference use
var hamburgerIDArray = new Array(); //Declare Array for id reference use
var hotdogsArray = new Array(); //All the hot dog object arrays
var hamburgArray = new Array(); //All the hamburger object arrays
var hotdogPicArray = new Array(); //All the hotdog pictures
var hamburgPicArray = new Array();  //All the hamburger pictures
var hDogSrc = new Array();  //A list of all food src for hotdogs
var hamSrc = new Array();   //A list of all food src for hamburgers
var mongoHotDogFoodID = new Array(); //All of the Mongo food IDS for Hotdogs
var mongoHamburgerFoodID = new Array(); //All of the Mongo Food IDS for Hamburgers
var hamburgerChoice = 1; //When the User updates data, this displays the seletced Burger Choices
var hotDogChoice = 1; //When the User updates data, this displays the seletced Hotdog Choices
var toSend = {
    UserName: "",
    Password: "", //This was formally a []byte but we are changing our code to fit the database better
    First:    "",
    Last:     "",
    Role:     "",
    UserID:   0
};
var emptyArray = new Array(); //only used to put this value in aHotdog and aHamburger
var aHotDog = {
    HotDogType: "",
    Condiments: emptyArray,
    Calories: 0,
    Name: "",
    FoodID: 0,
    UserID: 0,
    PhotoID: 0,
    PhotoSrc: "",
    DateCreated: "",
    DateUpdated: ""
}
var aHamburger = {
    BurgerType: "",
    Condiments: emptyArray,
    Calories: 0,
    Name: "",
    FoodID: 0,
    UserID: 0,
    PhotoID: 0,
    PhotoSrc: "",
    DateCreated: "",
    DateUpdated: ""
}

function elementPasser(seeFood, foodSeer, foodUpDate, hdogSeer, hamSeer, foodUpdateButton, foodDeleteButton){
    seeFoodButton = seeFood;
    foodSeerDiv = foodSeer;
    foodUpdater = foodUpDate;
    hdogSeerDiv = hdogSeer;
    hamSeerDiv = hamSeer;
    foodUpdateBtn = foodUpdateButton;
    foodDeleteBtn = foodDeleteButton;
}

function foodChanger(whichFood, whichChoice, hamburgObj, hotdogObj){
    console.log("DEBUG: Reached the foodUpdater.");
    //Make JSON to send to food API
    var foodUpdate = {
        FoodType: whichFood,
        FoodID: whichChoice,
        TheHamburger: hamburgObj,
        TheHotDog: hotdogObj
    };

    var jsonString = JSON.stringify(foodUpdate);

    //Call Ajax to update the foodRecord (SQL Database)
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/updateFood', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            console.log("DEBUG: " + item);
            if (item == 1) {
                //Hotdog Update
                console.log(item);
                console.log("Hotdog Updated!")
                location.reload();
            } else if (item == 2) {
                //Hamburger Update
                console.log(item);
                console.log("Hamburger Updated");
                location.reload();
            } else if (item == 3) {
                console.log(item);
                alert("No food item updated; something went cooky :(");
                location.reload();
            } else if (item == 4){
                console.log(item);
                alert("The food contained derogatory terms. Please update food again.");
                location.reload();
            } else {
                alert("No food item updated; something went cooky :( ");
                console.log("Unexpected output at foodUpdater function.");
                location.reload();
            }
        }
    });
    xhr.send(jsonString);
    //Call Ajax to update the foodRecord (Mongo Database)
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/foodUpdateMongo', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            console.log("DEBUG: " + item);
            if (item == 1) {
                //Hotdog Update
                console.log(item);
                console.log("Hotdog updated");
                location.reload();
            } else if (item == 2) {
                //Hamburger Update
                console.log(item);
                console.log("Hamburger Updated");
                location.reload();
            } else if (item == 3) {
                console.log(item);
                alert("No food item updated; something went cooky :( ");
                location.reload();
            } else if (item ==4){
                console.log(item);
                alert("The food contained derogatory terms. Please update food again.");
                location.reload();
            } else {
                alert("No food item updated; something went cooky :( ");
                console.log("Unexpected output at foodUpdater function.");
                location.reload();
            }
        }
    });
    xhr.send(jsonString);

    /* CHECK TO SEE IF THERE'S A PHOTO SUBMITTED TO CHANGE */
    var photoSubmission = fileUploadCheck();
    if (photoSubmission == true){
        //Begin process of submitting submitted photo for an update
        photoUpdate(foodUpdate);
    } else {

    }
}

function foodDeleter(whichFood, whichChoice, whichUserID){
    console.log("DEBUG: Reached the foodDeleter.");
    var foodDeletion = {
        FoodType: whichFood,
        FoodID: whichChoice,
        UserID: whichUserID
    }; //Make JSON to send to food API for SQL

    var jsonString = JSON.stringify(foodDeletion);
    console.log("Here is our json string for food deletion: " + jsonString);
    //Delete Food Photo First from Amazon and directory
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/deletePhotoFromS3', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            if (item == 1){
                console.log(item);
            }else if (item == 2){
                console.log(item);
            } else if (item == 3) {
                console.log(item);
            } else {
                console.log("Unexpected output at foodDeletion function: " + item);
            }
        }
    });
    xhr.send(jsonString);
    //Call Ajax to delete the foodRecord(SQL Database)
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/deleteFood', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            if (item == 1){
                console.log(item);
            }else if (item == 2){
                console.log(item);
            } else if (item == 3) {
                console.log(item);
            } else {
                console.log("Unexpected output at foodDeletion function: " + item);
            }
        }
    });
    xhr.send(jsonString);
    //Call Ajax to delete the foodRecord(Mongo Database)
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/foodDeleteMongo', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            if (item == 1){
                console.log(item);
            }else if (item == 2){
                console.log(item);
            } else if (item == 3) {
                console.log(item);
            } else {
                console.log("Unexpected output at foodDeletion function for Mongo: " + item);
            }
        }
    });
    xhr.send(jsonString);

    //Reload location once finished making API calls
    location.reload();
}

function srcChecker(theHamburger, theHotdog, whichFood){
    var srcThere = true;
    //Determine if we're sending the hotdog's foodSrc or the Hamburger's
    var theSrc = "";
    if (whichFood === 0){
        //Send Hamburger Data
        theSrc = String(theHamburger.PhotoSrc);
        console.log("Here is our photo src: " + theSrc);
    } else if (whichFood === 1){
        //Send Hotdog Data
        theSrc = String(theHotdog.PhotoSrc);
        console.log("Here is our photo src: " + theSrc);
    } else {
        console.log("Error, unidentified 'whichFood' returned: " + whichFood);
    }

    var photoChecker = {
        TheSrc: theSrc,
        TheFood: whichFood
    };
    

    var jsonString = JSON.stringify(photoChecker);
    console.log("Checking to see if photo exists"); //DEBUG
    //Call Ajax to check the fileRecord
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/checkSRC', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var returnOBJ = JSON.parse(item);
            if (returnOBJ.FoundSRC == true) {
                //PicutreSRC found
                console.log("Photo found"); //DEBUG
                srcThere = true;
            } else {
                //PictureSRC not found
                console.log("Photo not found"); //DEBUG
                srcThere = false;
            }
        }
    });
    xhr.send(jsonString);

    return srcThere
}