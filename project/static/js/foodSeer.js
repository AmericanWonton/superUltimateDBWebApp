var userID;
function getUserID(passedID){
    userID = passedID;
    console.log("We've set userID to " + userID);
}
window.addEventListener("load", function(){
    console.log("DEBUG: Hey, here's foodSeer.js");
});

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

seeFoodButton.addEventListener("click", function(){
    toSend.UserID = userID;
    var jsonString = JSON.stringify(toSend);

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/getAllFoodUser', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var dataOBJ = JSON.parse(item);
            if (dataOBJ.SuccessMessage == "Success") {
                //Clear all arrays, just in case
                hotdogIDArray = []; //Declare Array for id reference use
                hamburgerIDArray = []; //Declare Array for id reference use
                hotdogsArray = []; //All the hot dog object arrays
                hamburgArray = []; //All the hamburger object arrays
                hamSeerDiv.innerHTML = ""; //Clear hamseer div
                hdogSeer.innerHTML = ""; //Clear hotdog div
                var hotdogsGoodCount = true; //Used for not displaying hotdogs
                var hamburgerGoodCount = true; //Used for displaying hotdogs
                //Update Hotdog and Hamburger Arrays
                if (dataOBJ.TheHamburgers[0].BurgerType == "DEBUGTYPE"){
                    console.log("DEBUG: We have no burgers to deliver");
                    hamburgerGoodCount = false;
                } else {
                    for (var x = 0; x < dataOBJ.ID_Hamburgers.length; x++) {
                        hamburgerIDArray.push(dataOBJ.ID_Hamburgers[x]);
                        hamburgerChoice = hamburgerChoice + 1;
                        aHamburger.BurgerType = dataOBJ.TheHamburgers[x].BurgerType;
                        aHamburger.Condiment = dataOBJ.TheHamburgers[x].Condiment;
                        aHamburger.Calories = dataOBJ.TheHamburgers[x].Calories;
                        aHamburger.Name = dataOBJ.TheHamburgers[x].Name;
                        aHamburger.UserID = dataOBJ.TheHamburgers[x].UserID;
                        hamburgArray.push(aHamburger); //Add the Hamburger to the Hamburger array
                    }
                }
                
                if (dataOBJ.TheHotDogs[0].HotDogType == "DEBUGTYPE"){
                    console.log("DEBUG: WE AIN'T FOUND SHIT FOR HOTDOGS");
                    hotdogsGoodCount = false;
                } else {
                    for (var y = 0; y < dataOBJ.ID_HotDogs.length; y++) {
                        console.log("DEBUG: We're adding hotdogs ids");
                        hotdogIDArray.push(dataOBJ.ID_HotDogs[y]);
                        hotDogChoice = hotDogChoice + 1;
                        aHotDog.HotDogType = dataOBJ.TheHotDogs[y].HotDogType;
                        aHotDog.Condiment = dataOBJ.TheHotDogs[y].Condiment;
                        aHotDog.Calories = dataOBJ.TheHotDogs[y].Calories;
                        aHotDog.Name = dataOBJ.TheHotDogs[y].Name;
                        aHotDog.UserID = dataOBJ.TheHotDogs[y].UserID;
                        hotdogsArray.push(aHotDog); //Add the hotdog to the hotdog array
                    }
                }
                
                //Display all Hamburgers
                if (hamburgerGoodCount == false){
                    console.log("We won't be displaying hamburgers.");
                } else {
                    for (var y = 0; y < dataOBJ.ID_Hamburgers.length; y++) {
                        //Create container for Hamburger info
                        var udContainerDiv = document.createElement("div");
                        udContainerDiv.setAttribute("id", "udContainerDiv");

                        //Create ID info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "HamburgerID";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = y;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Type info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Hamburger Type";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHamburgers[y].BurgerType;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Condiment info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Condiment";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHamburgers[y].Condiment;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Calories info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Calories";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHamburgers[y].Calories;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Name info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Name";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHamburgers[y].Name;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //All the divs fields for the container are added. Add this to the hamseer!
                        hamSeerDiv.appendChild(udContainerDiv);
                    }
                    hamSeerDiv.style.display = "block"; //Display the formatted Hamburger Divs
                }
                
                //Display all HotDogs
                if (hotdogsGoodCount == false){
                    console.log("We won't be displaying hot dogs.");
                } else {
                    for (var i = 0; i < dataOBJ.TheHotDogs.length; i++) {
                        console.log("We're displaying hotodogs");
                        //Create container for hotdog info
                        var udContainerDiv = document.createElement("div");
                        udContainerDiv.setAttribute("id", "udContainerDiv");

                        //Create ID info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "HotdogID";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = i;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Type info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Hotdog Type";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHotDogs[i].HotDogType;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Condiment info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Condiment";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHotDogs[i].Condiment;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Calories info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Calories";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHotDogs[i].Calories;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //Create Name info section
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = "Name";
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container
                        var udcFieldInfoTxtDiv = document.createElement("div");
                        udcFieldInfoTxtDiv.setAttribute("id", "udcFieldInfoTxtDiv");
                        var udcfText = document.createElement("p");
                        udcfText.setAttribute("id", "udcfText");
                        udcfText.innerHTML = dataOBJ.TheHotDogs[i].Name;
                        udcFieldInfoTxtDiv.appendChild(udcfText); //Append text to Div
                        udContainerDiv.appendChild(udcFieldInfoTxtDiv); //Append infoDiv to container

                        //All the divs fields for the container are added. Add this to the hamseer!
                        hdogSeerDiv.appendChild(udContainerDiv);

                    }
                    hdogSeerDiv.style.display = "block"; //Display the formattted Hotdog divs
                }
                
                foodUpdater.style.display = "block"; //Display the foodUpdater div

            } else {
                //No data retrieved
                console.log("Either User has no data or no data could be retrieved.");
                alert("Either User has no data or no data could be retrieved.");
                //DEBUG: PRINT A MESSAGE TO USER HERE.
            }
            
        }
    });
    xhr.send(jsonString);
});

//Call the food deleter function to run a query and delete food,(if the selection is appropriate)
foodDeleteBtn.addEventListener("click", function(){
    //Declare input variable names to work with
    var foodChoice = document.getElementById("foodChoice");
    var foodID = document.getElementById("foodID");
    

    //Assign fields to pass to food deleter
    if (foodChoice.value == "hotdog"){
        //See if the foodID field was assigned correctly
        if (foodID.value < 0 || foodID.value > hotdogsArray.length - 1){
            //bad selection
            console.log("User inputted a wrong selection: " + foodID.value);
            alert("Please choose an ID for hotdog listed below.");
            //location.reload(true);
        } else {    
            //Good Selection
            console.log("DEBUG: Deleting a hotdog");
            foodDeleter(foodChoice.value, hotdogIDArray[foodID.value]);
        }
    } else if (foodChoice.value == "hamburger"){
        //See if the foodID field was assigned correctly
        if (foodID.value < 0 || foodID.value > hamburgArray.length - 1){
            //bad selection
            console.log("User inputted a wrong selection: " + foodID.value);
            alert("Please choose an ID for hotdog listed below.");
        } else {    
            //Good Selection
            console.log("DEBUG: deleting a hamburger");
            foodDeleter(foodChoice.value, hamburgerIDArray[foodID.value]);
        }
    } else {
        console.log("Error, User selected nothing in the foodchoice" +
        " field or wrong value got put in place.");
        alert("Please select hamburger or hotdog food to delete.");
        location.reload(true);
    }
});
//call the foodChanger function and update the food if the selection is good
foodUpdateBtn.addEventListener("click", function(){
    //Declare input variable names to work with
    var foodChoice = document.getElementById("foodChoice");
    var foodID = document.getElementById("foodID");
    var foodIDNumber = Number(foodID.value);
    var foodType = document.getElementById("foodType");
    var condimentType = document.getElementById("condimentType");
    var calories = document.getElementById("calories");
    var caloriesNumber = Number(calories.value);
    var foodName = document.getElementById("foodName");

    //Assign fields to pass to food deleter
    if (foodChoice.value == "hotdog"){
        //See if the foodID field was assigned correctly
        if (foodIDNumber < 0 || foodIDNumber > hotdogsArray.length - 1){
            //bad selection
            console.log("User inputted a wrong selection: " + foodIDNumber);
            alert("Please choose an ID for hotdog listed below.");
            location.reload(true);
        } else {    
            //Good Selection
            aHotDog.HotDogType = foodType.value;
            aHotDog.Condiment = condimentType.value;
            aHotDog.Calories = caloriesNumber;
            aHotDog.Name = foodName.value;
            aHotDog.UserID = userID;
            foodChanger(foodChoice.value, hotdogIDArray[foodIDNumber], aHamburger, aHotDog);
        }
    } else if (foodChoice.value == "hamburger"){
        //See if the foodID field was assigned correctly
        if (foodIDNumber < 0 || foodIDNumber > hamburgArray.length - 1){
            //bad selection
            console.log("User inputted a wrong selection: " + foodIDNumber);
            alert("Please choose an ID for hotdog listed below.");
            location.reload(true);
        } else {    
            //Good Selection
            aHamburger.BurgerType = foodType.value;
            aHamburger.Condiment = condimentType.value;
            aHamburger.Calories = caloriesNumber;
            aHamburger.Name = foodName.value;
            aHamburger.UserID = userID;
            foodChanger(foodChoice.value, hamburgerIDArray[foodIDNumber], aHamburger, aHotDog);
        }
    } else {
        console.log("Error, User selected nothing in the foodchoice" +
        " field or wrong value got put in place.");
        alert("Please select hamburger or hotdog food to update.");
        location.reload(true);
    }
});