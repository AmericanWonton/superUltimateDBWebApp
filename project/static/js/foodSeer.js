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
var aHotDog = {
    HotDogType: "",
    Condiment: "",
    Calories: 0,
    Name: "",
    UserID: 0
}
var aHamburger = {
    BurgerType: "",
    Condiment: "",
    Calories: 0,
    Name: "",
    UserID: 0
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
    //Call Ajax to delete the foodRecord
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
                alert("Hotdog updated");
                location.reload(true);
            } else if (item == 2) {
                //Hamburger Update
                console.log(item);
                alert("Hamburger Updated");
                location.reload(true);
            } else if (item == 3) {
                console.log(item);
                alert("No food item updated; something went cooky :( ");
                location.reload(true);
            } else {
                alert("No food item updated; something went cooky :( ");
                console.log("Unexpected output at foodUpdater function.");
                location.reload(true);
            }
        }
    });
    xhr.send(jsonString);
}

function foodDeleter(whichFood, whichChoice){
    console.log("DEBUG: Reached the foodDeleter.");
    var foodDeletion = {
        FoodType: whichFood,
        FoodID: whichChoice
    }; //Make JSON to send to food API

    console.log(foodDeletion);

    var jsonString = JSON.stringify(foodDeletion);

    //Call Ajax to delete the foodRecord
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/deleteFood', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            if (item == 1){
                console.log(item);
                alert("Hotdog successfully deleted");
                location.reload(true);
            }else if (item == 2){
                console.log(item);
                alert("Hamburger Deleted");
                location.reload(true);
            } else if (item == 3) {
                console.log(item);
                alert("Trouble deleting food item");
                location.reload(true);
            } else {
                console.log("Unexpected output at foodDeletion function: " + item);
                location.reload(true);
            }
        }
    });
    xhr.send(jsonString);
}