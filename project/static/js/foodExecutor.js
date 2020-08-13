var userID;
function getUserID(passedID){
    userID = passedID;
    console.log("We've set userID to " + userID);
}

function revealFoodForm(foodChoice) {
    console.log(foodChoice);
    var theDiv = document.getElementById("foodFormDiv"); //Get the div
    theDiv.innerHTML = ""; //Remove any child elements if any remain
    /* For Hotdog selection */
    if (foodChoice == 0) {
        console.log("Mkay, you clicked hotdog.")
        //Create the form elements and append them to the form
        /* Let's start with form instruction first */
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hotdog multiple condiments, give a space between each condiment.";
        var hDogType = document.createElement("input");
        hDogType.setAttribute("type", "text");
        hDogType.setAttribute("id", "hotDogType");
        hDogType.setAttribute("maxlength", 48);
        hDogType.setAttribute("name", "hDogType");
        hDogType.setAttribute("placeholder", "HotDog Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("placeholder", "Condiment");
        var caloriesType = document.createElement("input");
        caloriesType.setAttribute("type", "number");
        caloriesType.setAttribute("maxlength", 8);
        caloriesType.setAttribute("id", "caloriesType");
        caloriesType.setAttribute("name", "caloriesType");
        caloriesType.setAttribute("placeholder", "Calories");
        var nameType = document.createElement("input");
        nameType.setAttribute("id", "nameType");
        nameType.setAttribute("type", "text");
        nameType.setAttribute("maxlength", 48);
        nameType.setAttribute("name", "nameType");
        nameType.setAttribute("placeholder", "HotDog Name");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        
        submitButton.addEventListener("click", function(){
            console.log("Submit button clicked, submitting hotdog data.");
            //Ajax functionality for submitting forms
            console.log("DEBUG: We're submitting the form.");
            //Field correction if fields aren't filled out
            if (hDogType.value == ""){
                hDogType.value ="NONE";
            }
            if (condimentType.value == ""){
                condimentType.value = "NONE";
            }
            if (caloriesType.value <= 0) {
                caloriesType.value = 0;
            }
            if (nameType.value == ""){
                nameType.value = "NONE";
            }
            
            //JSON String creation
            var toSend = {
                HotDogType: "",
                Condiment: "",
                Calories: Number(caloriesType.value),
                Name: "",
                UserID: userID,
                FoodID: 0,
                DateCreated: "",
                DateUpdated: ""
            };
            
            toSend.HotDogType = hDogType.value;
            toSend.Condiment = condimentType.value;
            toSend.Name = nameType.value;
            toSend.UserID = userID;
            //Stringify our JSON
            var jsonString = JSON.stringify(toSend);
            console.log("DEBUG: Here's our JSON we're sending: " + jsonString);
            //Send our hotdog to create
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/hotDogInsertWebPage', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    var dataReturned = JSON.parse(item);
                    if (dataReturned.SuccessBool === true) {
                        //Data inserted properly; clear the form fields
                        hDogType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        console.log("Hotdog submitted successfully!");
                        //alert("Hotdog submitted successfully!");
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        //location.reload(true); //Reload Page
                    } else if (dataReturned.SuccessBool === false){
                        var foulLanguage = dataReturned.SuccessMsg.includes("foul language");
                        if (foulLanguage == true){
                            hDogType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            console.log("Hotdog contained foul language. Please re-enter.");
                            //alert("There was an issue submitting your hotdog :(");
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            //location.reload(true); //Reload Page
                        } else {
                            //Data NOT inserted properly
                            hDogType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            console.log("There was an issue submitting your hotdog :(");
                            //alert("There was an issue submitting your hotdog :(");
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            //location.reload(true); //Reload Page
                        }
                    } else {
                        //No appropriate Response recieved
                        console.log("Error submitting data for hotDogInsertWebPage, please send again.");
                        //alert("Error submitting data for insertHotDogMongo, please send again.");
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        //location.reload(true); //Reload Page
                    }
                }
            });
            xhr.send(jsonString);
        });

        //Append "Form Data"
        theDiv.appendChild(condimentInstruction);
        theDiv.appendChild(hDogType);
        theDiv.appendChild(condimentType);
        theDiv.appendChild(caloriesType);
        theDiv.appendChild(nameType);
        theDiv.appendChild(submitButton);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else if (foodChoice == 1) { //For Hamburger Selection
        console.log("Mkay, you clicked hamburger.")
        //Create the form elements and append them to the form
        /* Let's start with form instruction first */
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hamburger multiple condiments, give a space between each condiment.";
        var hamburgType = document.createElement("input");
        hamburgType.setAttribute("type", "text");
        hamburgType.setAttribute("id", "hamburgType");
        hamburgType.setAttribute("maxlength", 48);
        hamburgType.setAttribute("name", "hamburgType");
        hamburgType.setAttribute("placeholder", "Hamburger Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("placeholder", "Condiment");
        var caloriesType = document.createElement("input");
        caloriesType.setAttribute("type", "number");
        caloriesType.setAttribute("maxlength", 8);
        caloriesType.setAttribute("id", "caloriesType");
        caloriesType.setAttribute("name", "caloriesType");
        caloriesType.setAttribute("placeholder", "Calories");
        var nameType = document.createElement("input");
        nameType.setAttribute("id", "nameType");
        nameType.setAttribute("type", "text");
        nameType.setAttribute("maxlength", 48);
        nameType.setAttribute("name", "nameType");
        nameType.setAttribute("placeholder", "Hamburger Name");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        
        submitButton.addEventListener("click", function(){
            console.log("Submit button clicked, submitting hamburger data.");
            //Ajax functionality for submitting forms
            //Field correction if fields aren't filled out
            if (hamburgType.value == ""){
                hamburgType.value ="NONE";
            }
            if (condimentType.value == ""){
                condimentType.value = "NONE";
            }
            if (caloriesType.value <= 0) {
                caloriesType.value = 0;
            }
            if (nameType.value == ""){
                nameType.value = "NONE";
            }

            //JSON String creation
            var toSend = {
                BurgerType: "",
                Condiment: "",
                Calories: Number(caloriesType.value),
                Name: "",
                UserID: 0,
                FoodID: 0,
                DateCreated: "",
                DateUpdated: ""
            };
            
            toSend.BurgerType = hamburgType.value;
            toSend.Condiment = condimentType.value;
            toSend.Name = nameType.value;
            toSend.UserID = userID;

            var jsonString = JSON.stringify(toSend);
            console.log(jsonString);
            //For SQL Database
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/hamburgerInsertWebPage', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    var dataReturned = JSON.parse(item);
                    if (dataReturned.SuccessBool === true) {
                        //Data inserted properly; clear the form fields
                        hamburgType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        alert("Hamburger submitted successfully!")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        location.reload(true); //Reload Page
                    } else if (dataReturned.SuccessBool === false){
                        var foulLanguage = dataReturned.SuccessMsg.includes("foul language");
                        if (foulLanguage == true){
                            //Data NOT inserted properly
                            hamburgType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            alert("Hamburger contained foul language. Please re-enter.");
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            location.reload(true); //Reload Page
                        }else{
                            //Data NOT inserted properly
                            hamburgType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            alert("There was an issue submitting your hamburger :(")
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            location.reload(true); //Reload Page
                        }
                    } else {
                        //No appropriate Response recieved
                        alert("Error submitting data, please send again.")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        location.reload(true); //Reload Page
                    }
                }
            });
            xhr.send(jsonString);
        });
        //Append "Form Data"
        theDiv.appendChild(condimentInstruction);
        theDiv.appendChild(hamburgType);
        theDiv.appendChild(condimentType);
        theDiv.appendChild(caloriesType);
        theDiv.appendChild(nameType);
        theDiv.appendChild(submitButton);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else if (foodChoice == 2){//For IT/Admin Hotdog Selection
        console.log("Mkay, you clicked hotdog as an Admin or IT.")
        //Create the form elements and append them to the form
        /* Let's start with form instruction first */
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hotdog multiple condiments, give a space between each condiment.";
        var hDogType = document.createElement("input");
        hDogType.setAttribute("type", "text");
        hDogType.setAttribute("id", "hotDogType");
        hDogType.setAttribute("maxlength", 48);
        hDogType.setAttribute("name", "hDogType");
        hDogType.setAttribute("placeholder", "HotDog Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("placeholder", "Condiment");
        var caloriesType = document.createElement("input");
        caloriesType.setAttribute("type", "number");
        caloriesType.setAttribute("maxlength", 8);
        caloriesType.setAttribute("id", "caloriesType");
        caloriesType.setAttribute("name", "caloriesType");
        caloriesType.setAttribute("placeholder", "Calories");
        var nameType = document.createElement("input");
        nameType.setAttribute("id", "nameType");
        nameType.setAttribute("type", "text");
        nameType.setAttribute("maxlength", 48);
        nameType.setAttribute("name", "nameType");
        nameType.setAttribute("placeholder", "HotDog Name");
        var userIDInput = document.createElement("input");
        userIDInput.setAttribute("id", "userIDInput");
        userIDInput.setAttribute("type", "number");
        userIDInput.setAttribute("maxlength", 8);
        userIDInput.setAttribute("name", "userIDInput");
        userIDInput.setAttribute("placeholder", "userID");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        
        submitButton.addEventListener("click", function(){
            console.log("Submit button clicked, submitting hotdog data.");
            //Ajax functionality for submitting forms
            console.log("DEBUG: We're submitting the form.");
            //Field correction if fields aren't filled out
            if (hDogType.value == ""){
                hDogType.value ="NONE";
            }
            if (condimentType.value == ""){
                condimentType.value = "NONE";
            }
            if (caloriesType.value <= 0) {
                caloriesType.value = 0;
            }
            if (nameType.value == ""){
                nameType.value = "NONE";
            }
            if (userIDInput.value.length == 0) {
                userIDInput.innerHTML = userID;
                userIDInput.value = userID;
            }

            //JSON String creation
            var toSend = {
                HotDogType: "",
                Condiment: "",
                Calories: Number(caloriesType.value),
                Name: "",
                UserID: 0,
                FoodID: 0,
                DateCreated: "",
                DateUpdated: ""
            };
            var theIDNumber = Number(userIDInput.value);
            toSend.HotDogType = hDogType.value;
            toSend.Condiment = condimentType.value;
            toSend.Name = nameType.value;
            toSend.UserID = theIDNumber;

            var jsonString = JSON.stringify(toSend);
            console.log(jsonString);
            //For SQL Insertion
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/hotDogInsertWebPage', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    var dataReturned = JSON.parse(item);
                    if (dataReturned.SuccessBool === true) {
                        //Data inserted properly; clear the form fields
                        hDogType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        userIDInput.value = "";
                        alert("Hotdog submitted successfully!")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        location.reload(true); //Reload Page
                    } else if (dataReturned.SuccessBool === false){
                        var foulLanguage = dataReturned.SuccessMsg.includes("foul language");
                        if (foulLanguage == true){
                            //Data NOT inserted properly
                            hDogType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            userIDInput.value = "";
                            alert("Hotdog contained foul language. Please re-enter.");
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            location.reload(true); //Reload Page
                        } else {
                            //Data NOT inserted properly
                            hDogType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            userIDInput.value = "";
                            alert("There was an issue submitting your hotdog :(");
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            location.reload(true); //Reload Page
                        }
                    } else {
                        //No appropriate Response recieved
                        alert("Error submitting data, please send again.")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        location.reload(true); //Reload Page
                    }
                }
            });
            xhr.send(jsonString);
        });
        //Append "Form Data"
        theDiv.appendChild(condimentInstruction);
        theDiv.appendChild(hDogType);
        theDiv.appendChild(condimentType);
        theDiv.appendChild(caloriesType);
        theDiv.appendChild(nameType);
        theDiv.appendChild(userIDInput);
        theDiv.appendChild(submitButton);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else if (foodChoice == 3){//For IT/Admin Hamburger Selection
        console.log("Mkay, you clicked hamburger.")
        //Create the form elements and append them to the form
        /* Let's start with form instruction first */
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hamburger multiple condiments, give a space between each condiment.";
        var hamburgType = document.createElement("input");
        hamburgType.setAttribute("type", "text");
        hamburgType.setAttribute("id", "hamburgType");
        hamburgType.setAttribute("maxlength", 48);
        hamburgType.setAttribute("name", "hamburgType");
        hamburgType.setAttribute("placeholder", "Hamburger Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("placeholder", "Condiment");
        var caloriesType = document.createElement("input");
        caloriesType.setAttribute("type", "number");
        caloriesType.setAttribute("maxlength", 8);
        caloriesType.setAttribute("id", "caloriesType");
        caloriesType.setAttribute("name", "caloriesType");
        caloriesType.setAttribute("placeholder", "Calories");
        var nameType = document.createElement("input");
        nameType.setAttribute("id", "nameType");
        nameType.setAttribute("type", "text");
        nameType.setAttribute("maxlength", 48);
        nameType.setAttribute("name", "nameType");
        nameType.setAttribute("placeholder", "Hamburger Name");
        var userIDInput = document.createElement("input");
        userIDInput.setAttribute("id", "userIDInput");
        userIDInput.setAttribute("type", "number");
        userIDInput.setAttribute("maxlength", 8);
        userIDInput.setAttribute("name", "userIDInput");
        userIDInput.setAttribute("placeholder", "userID");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        
        submitButton.addEventListener("click", function(){
            console.log("Submit button clicked, submitting hamburger data.");
            //Ajax functionality for submitting forms
            console.log("DEBUG: We're submitting the form.");
            //Field correction if fields aren't filled out
            if (hamburgType.value == ""){
                hamburgType.value ="NONE";
            }
            if (condimentType.value == ""){
                condimentType.value = "NONE";
            }
            if (caloriesType.value <= 0) {
                caloriesType.value = 0;
            }
            if (nameType.value == ""){
                nameType.value = "NONE";
            }
            if (userIDInput.value.length == 0) {
                userIDInput.innerHTML = userID;
                userIDInput.value = userID;
            }

            //JSON String creation
            var toSend = {
                BurgerType: "",
                Condiment: "",
                Calories: Number(caloriesType.value),
                Name: "",
                UserID: 0,
                FoodID: 0,
                DateCreated: "",
                DateUpdated: ""
            };
            var theIDNumber = Number(userIDInput.value);
            toSend.BurgerType = hamburgType.value;
            toSend.Condiment = condimentType.value;
            toSend.Name = nameType.value;
            toSend.UserID = theIDNumber;

            var jsonString = JSON.stringify(toSend);
            console.log(jsonString);
            //For SQL Database
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/hamburgerInsertWebPage', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    var dataReturned = JSON.parse(item);
                    if (dataReturned.SuccessBool === true) {
                        //Data inserted properly; clear the form fields
                        hamburgType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        userIDInput.value = "";
                        alert("Hamburger submitted successfully!")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        location.reload(true); //Reload Page
                    } else if (dataReturned.SuccessBool === false){
                        var foulLanguage = dataReturned.SuccessMsg.includes("foul language");
                        if (foulLanguage == true){
                            //Data NOT inserted properly
                            hamburgType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            userIDInput.value = "";
                            alert("Hamburger contained foul language. Please re-enter.");
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            location.reload(true); //Reload Page
                        } else{
                            //Data NOT inserted properly
                            hamburgType.value = "";
                            condimentType.value = "";
                            caloriesType.value = "";
                            nameType.value = "";
                            userIDInput.value = "";
                            alert("There was an issue submitting your hamburger :(");
                            theDiv.innerHTML = ""; //Remove any child elements if any remain
                            location.reload(true); //Reload Page
                        }
                    } else {
                        //No appropriate Response recieved
                        alert("Error submitting data, please send again.")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                        location.reload(true); //Reload Page
                    }
                }
            });
            xhr.send(jsonString);
        });
        //Append "Form Data"
        theDiv.appendChild(condimentInstruction);
        theDiv.appendChild(hamburgType);
        theDiv.appendChild(condimentType);
        theDiv.appendChild(caloriesType);
        theDiv.appendChild(nameType);
        theDiv.appendChild(userIDInput);
        theDiv.appendChild(submitButton);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else {
        console.log("Whoops, we got a problem. Wrong food choice came in: " + foodChoice);
        location.reload(true); //Reload Page
    }
}