var userID;
function getUserID(passedID){
    userID = passedID;
    console.log("We've set userID to " + userID);
}

function revealFoodForm(foodChoice) {
    var theDiv = document.getElementById("foodFormDiv"); //Get the div
    theDiv.innerHTML = ""; //Remove any child elements if any remain
    /* For Hotdog selection */
    if (foodChoice == 0) {
        console.log("Mkay, you clicked hotdog.")
        //Create the form elements and append them to the form
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
            var condimentDEBUG = String(condimentType.value);
            var nameDEBUG = nameType.value;
            console.log("CondimentDEBUG is: " + condimentDEBUG + " and nameDEBUG is: " + nameDEBUG);
            //JSON String creation
            var toSend = {
                HotDogType: "",
                Condiment: "",
                Calories: Number(caloriesType.value),
                Name: "",
                UserID: 0
            };
            
            toSend.HotDogType = document.getElementById("hotDogType").value;
            toSend.Condiment = condimentType.value;
            toSend.Name = nameType.value;
            toSend.UserID = userID;

            var jsonString = JSON.stringify(toSend);
            console.log(jsonString);

            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/insertHotDog', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    if (item == 'Successful Insert') {
                        //Data inserted properly; clear the form fields
                        hDogType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        alert("Hotdog submitted successfully!")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                    } else if (item == 'Unsuccessful Insert'){
                        //Data NOT inserted properly
                        hDogType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        alert("There was an issue submitting your hotdog :(")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                    } else {
                        //No appropriate Response recieved
                        alert("Error submitting data, please send again.")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                    }
                }
            });
            xhr.send(jsonString);
            
        });
        //Append "Form Data"
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
        nameType.setAttribute("placeholder", "HotDog Name");
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
            //JSON String creation
            var toSend = {
                HamburgerType: "",
                Condiment: "",
                Calories: Number(caloriesType.value),
                Name: "",
                UserID: 0
            };
            
            toSend.HamburgerType = document.getElementById("hamburgType").value;
            toSend.Condiment = condimentType.value;
            toSend.Name = nameType.value;
            toSend.UserID = userID;

            var jsonString = JSON.stringify(toSend);
            console.log(jsonString);

            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/insertHamburger', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    if (item == 'Successful Insert') {
                        //Data inserted properly; clear the form fields
                        hamburgType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        alert("Hamburger submitted successfully!")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                    } else if (item == 'Unsuccessful Insert'){
                        //Data NOT inserted properly
                        hamburgType.value = "";
                        condimentType.value = "";
                        caloriesType.value = "";
                        nameType.value = "";
                        alert("There was an issue submitting your hamburger :(")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                    } else {
                        //No appropriate Response recieved
                        alert("Error submitting data, please send again.")
                        theDiv.innerHTML = ""; //Remove any child elements if any remain
                    }
                }
            });
            xhr.send(jsonString);
            
        });
        //Append "Form Data"
        theDiv.appendChild(hamburgType);
        theDiv.appendChild(condimentType);
        theDiv.appendChild(caloriesType);
        theDiv.appendChild(nameType);
        theDiv.appendChild(submitButton);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else {
        console.log("Whoops, we got a problem. Wrong food choice came in.")
    }
}