var userID;
var thePort;

function getUserID(passedID){
    userID = passedID;
    console.log("We've set userID to " + userID);
}

function getOtherHeaderValues(message, willDisplay){
    console.log("DEBUG: Delivered the message value, it should be this: " + message);
    //Define variables for execution
    var errordisplayerDiv = document.getElementById("errordisplayer_div"); //Div Display
    var errordisplayerP = document.getElementById("errordisplayer_p"); //Message display
    //Determine if this div should display
    console.log("DEBUG: willDisplay: " + willDisplay);
    if (willDisplay === 0){
        //Set the message of the P
        errordisplayerP.innerHTML = message;
        errordisplayerDiv.style.display = "flex";
        console.log("DEBUG: We should be seeing the div now: " + willDisplay);
        //Fade pic back up after a few seconds

        setTimeout(function() { rollup(); }, 3500); //Start fading away
        setTimeout(function() { dissapearDiv(); }, 6000);

        function rollup(){
            console.log("DEBUG: We rolling up...");
            errordisplayerDiv.style.transition = "height 2s linear 0s";
            errordisplayerDiv.style.height = "0px";
            errordisplayerP.style.height = "0px";
        }

        function dissapearDiv(){
            console.log("DEBUG: Bye Bye Div");
            errordisplayerDiv.style.display = "none";
        }

    } else if (willDisplay === 1){
        //Do nothing, we don't mean to display anything
    } else {
        console.log("Error, willDisplay is incorrect: " + willDisplay);
    }
}

function getPort(passedPort){
    thePort = "http://localhost:" + passedPort;
}

function revealFoodForm(foodChoice) {
    var theDiv = document.getElementById("foodFormDiv"); //Get the div
    theDiv.innerHTML = ""; //Remove any child elements if any remain
    /* For Hotdog selection */
    if (foodChoice == 0) {
        //Inform User
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hotdog multiple condiments, give a space between each condiment.";
        //Create the form elements and append them to the form
        /* Create Food Form */
        var documentForm = document.createElement("form");
        documentForm.setAttribute("id", "submit-picture-form");
        documentForm.setAttribute("name", "submit-picture-form");
        documentForm.setAttribute("enctype", "multipart/form-data");
        documentForm.setAttribute("action", "/mainPage");
        documentForm.setAttribute("method", "POST");
        documentForm.setAttribute("onload", "");
        var hDogType = document.createElement("input");
        hDogType.setAttribute("type", "text");
        hDogType.setAttribute("id", "hDogType");
        hDogType.setAttribute("name", "hDogType");
        hDogType.setAttribute("maxlength", 48);
        hDogType.setAttribute("placeholder", "Hotdog Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
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
        nameType.setAttribute("placeholder", "Hotdog Name");
        var userIDInput = document.createElement("input");
        userIDInput.setAttribute("id", "userIDInput");
        userIDInput.setAttribute("type", "hidden");
        userIDInput.setAttribute("name", "userIDInput");
        userIDInput.setAttribute("value", userID);
        var docButtonInput = document.createElement("input");
        docButtonInput.setAttribute("id", "docButtonInput");
        docButtonInput.setAttribute("name", "newFile");
        docButtonInput.setAttribute("type", "file");
        var hiddenUserNum = document.createElement("input");
        hiddenUserNum.setAttribute("type", "hidden");
        hiddenUserNum.setAttribute("id", "hiddenUserNum");
        hiddenUserNum.setAttribute("name", "hiddenUserNum");
        var hiddenFoodType = document.createElement("input");
        hiddenFoodType.setAttribute("type", "hidden");
        hiddenFoodType.setAttribute("id", "hiddenFoodType");
        hiddenFoodType.setAttribute("name", "hiddenFoodType");
        var hiddenFoodNum = document.createElement("input");
        hiddenFoodNum.setAttribute("type", "hidden");
        hiddenFoodNum.setAttribute("id", "hiddenFoodNum");
        hiddenFoodNum.setAttribute("name", "hiddenFoodNum");
        var hiddenFoodAction = document.createElement("input");
        hiddenFoodAction.setAttribute("type", "hidden");
        hiddenFoodAction.setAttribute("id", "hiddenFoodAction");
        hiddenFoodAction.setAttribute("name", "hiddenFoodAction");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.setAttribute("name", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        /* Add input to the above form for document selection */
        documentForm.appendChild(hDogType);
        documentForm.appendChild(condimentType);
        documentForm.appendChild(caloriesType);
        documentForm.appendChild(nameType);
        documentForm.appendChild(userIDInput);
        documentForm.appendChild(docButtonInput);
        documentForm.appendChild(hiddenUserNum);
        documentForm.appendChild(hiddenFoodType);
        documentForm.appendChild(hiddenFoodNum);
        documentForm.appendChild(hiddenFoodAction);
        documentForm.appendChild(submitButton);
        
        submitButton.addEventListener("click", function(){
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
                userIDInput.setAttribute("value", userID); 
            }
            //Give appropriate values for form to send
            hiddenFoodAction.setAttribute("value", "food_submit");
            hiddenFoodType.setAttribute("value", "HOTDOG");
            hiddenFoodNum.setAttribute("value", ""); //Not sure if needed now...
            hiddenUserNum.setAttribute("value", "");
            //Submit the form!
            documentForm.submit();
            
        });
        theDiv.appendChild(documentForm);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else if (foodChoice == 1) { //For Hamburger Selection
        //Inform User
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hamburger multiple condiments, give a space between each condiment.";
        //Create the form elements and append them to the form
        /* Create Food Form */
        var documentForm = document.createElement("form");
        documentForm.setAttribute("id", "submit-picture-form");
        documentForm.setAttribute("name", "submit-picture-form");
        documentForm.setAttribute("enctype", "multipart/form-data");
        documentForm.setAttribute("action", "/mainPage");
        documentForm.setAttribute("method", "POST");
        documentForm.setAttribute("onload", "");
        var hamburgType = document.createElement("input");
        hamburgType.setAttribute("type", "text");
        hamburgType.setAttribute("id", "hamburgType");
        hamburgType.setAttribute("name", "hamburgType");
        hamburgType.setAttribute("maxlength", 48);
        hamburgType.setAttribute("name", "hamburgType");
        hamburgType.setAttribute("placeholder", "Hamburger Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
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
        userIDInput.setAttribute("type", "hidden");
        userIDInput.setAttribute("name", "userIDInput");
        userIDInput.setAttribute("value", userID);
        var docButtonInput = document.createElement("input");
        docButtonInput.setAttribute("id", "docButtonInput");
        docButtonInput.setAttribute("name", "newFile");
        docButtonInput.setAttribute("type", "file");
        var hiddenUserNum = document.createElement("input");
        hiddenUserNum.setAttribute("type", "hidden");
        hiddenUserNum.setAttribute("id", "hiddenUserNum");
        hiddenUserNum.setAttribute("name", "hiddenUserNum");
        var hiddenFoodType = document.createElement("input");
        hiddenFoodType.setAttribute("type", "hidden");
        hiddenFoodType.setAttribute("id", "hiddenFoodType");
        hiddenFoodType.setAttribute("name", "hiddenFoodType");
        var hiddenFoodNum = document.createElement("input");
        hiddenFoodNum.setAttribute("type", "hidden");
        hiddenFoodNum.setAttribute("id", "hiddenFoodNum");
        hiddenFoodNum.setAttribute("name", "hiddenFoodNum");
        var hiddenFoodAction = document.createElement("input");
        hiddenFoodAction.setAttribute("type", "hidden");
        hiddenFoodAction.setAttribute("id", "hiddenFoodAction");
        hiddenFoodAction.setAttribute("name", "hiddenFoodAction");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.setAttribute("name", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        /* Add input to the above form for document selection */
        documentForm.appendChild(hamburgType);
        documentForm.appendChild(condimentType);
        documentForm.appendChild(caloriesType);
        documentForm.appendChild(nameType);
        documentForm.appendChild(userIDInput);
        documentForm.appendChild(docButtonInput);
        documentForm.appendChild(hiddenUserNum);
        documentForm.appendChild(hiddenFoodType);
        documentForm.appendChild(hiddenFoodNum);
        documentForm.appendChild(hiddenFoodAction);
        documentForm.appendChild(submitButton);
        
        submitButton.addEventListener("click", function(){
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
                userIDInput.setAttribute("value", userID); 
            }
            //Give appropriate values for form to send
            hiddenFoodAction.setAttribute("value", "food_submit")
            hiddenFoodType.setAttribute("value", "HAMBURGER");
            hiddenFoodNum.setAttribute("value", ""); //Not sure if needed now...
            hiddenUserNum.setAttribute("value", "");
            //Submit the form!
            documentForm.submit();
            
        });
        theDiv.appendChild(documentForm);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else if (foodChoice == 2){//For IT/Admin Hotdog Selection
        //Inform User
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hotdog multiple condiments, give a space between each condiment.";
        //Create the form elements and append them to the form
        /* Create Food Form */
        var documentForm = document.createElement("form");
        documentForm.setAttribute("id", "submit-picture-form");
        documentForm.setAttribute("name", "submit-picture-form");
        documentForm.setAttribute("enctype", "multipart/form-data");
        documentForm.setAttribute("action", "/mainPage");
        documentForm.setAttribute("method", "POST");
        documentForm.setAttribute("onload", "");
        var hDogType = document.createElement("input");
        hDogType.setAttribute("type", "text");
        hDogType.setAttribute("id", "hDogType");
        hDogType.setAttribute("name", "hDogType");
        hDogType.setAttribute("maxlength", 48);
        hDogType.setAttribute("placeholder", "Hotdog Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
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
        nameType.setAttribute("placeholder", "Hotdog Name");
        var userIDInput = document.createElement("input");
        userIDInput.setAttribute("id", "userIDInput");
        userIDInput.setAttribute("type", "number");
        userIDInput.setAttribute("maxlength", 8);
        userIDInput.setAttribute("name", "userIDInput");
        userIDInput.setAttribute("placeholder", "userID");
        var docButtonInput = document.createElement("input");
        docButtonInput.setAttribute("id", "docButtonInput");
        docButtonInput.setAttribute("name", "newFile");
        docButtonInput.setAttribute("type", "file");
        var hiddenUserNum = document.createElement("input");
        hiddenUserNum.setAttribute("type", "hidden");
        hiddenUserNum.setAttribute("id", "hiddenUserNum");
        hiddenUserNum.setAttribute("name", "hiddenUserNum");
        var hiddenFoodType = document.createElement("input");
        hiddenFoodType.setAttribute("type", "hidden");
        hiddenFoodType.setAttribute("id", "hiddenFoodType");
        hiddenFoodType.setAttribute("name", "hiddenFoodType");
        var hiddenFoodNum = document.createElement("input");
        hiddenFoodNum.setAttribute("type", "hidden");
        hiddenFoodNum.setAttribute("id", "hiddenFoodNum");
        hiddenFoodNum.setAttribute("name", "hiddenFoodNum");
        var hiddenFoodAction = document.createElement("input");
        hiddenFoodAction.setAttribute("type", "hidden");
        hiddenFoodAction.setAttribute("id", "hiddenFoodAction");
        hiddenFoodAction.setAttribute("name", "hiddenFoodAction");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.setAttribute("name", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        /* Add input to the above form for document selection */
        documentForm.appendChild(hDogType);
        documentForm.appendChild(condimentType);
        documentForm.appendChild(caloriesType);
        documentForm.appendChild(nameType);
        documentForm.appendChild(userIDInput);
        documentForm.appendChild(docButtonInput);
        documentForm.appendChild(hiddenUserNum);
        documentForm.appendChild(hiddenFoodType);
        documentForm.appendChild(hiddenFoodNum);
        documentForm.appendChild(hiddenFoodAction);
        documentForm.appendChild(submitButton);
        
        submitButton.addEventListener("click", function(){
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
                userIDInput.setAttribute("value", userID); 
            }
            //Give appropriate values for form to send
            hiddenFoodAction.setAttribute("value", "food_submit")
            hiddenFoodType.setAttribute("value", "HOTDOG");
            hiddenFoodNum.setAttribute("value", ""); //Not sure if needed now...
            hiddenUserNum.setAttribute("value", "");
            //Submit the form!
            documentForm.submit();
            
        });
        theDiv.appendChild(documentForm);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else if (foodChoice == 3){//For IT/Admin Hamburger Selection
        //Inform User
        var condimentInstruction = document.createElement("p");
        condimentInstruction.setAttribute("id", "condimentInstruction");
        condimentInstruction.innerHTML = "To give this hamburger multiple condiments, give a space between each condiment.";
        //Create the form elements and append them to the form
        /* Create Food Form */
        var documentForm = document.createElement("form");
        documentForm.setAttribute("id", "submit-picture-form");
        documentForm.setAttribute("name", "submit-picture-form");
        documentForm.setAttribute("enctype", "multipart/form-data");
        documentForm.setAttribute("action", "/mainPage");
        documentForm.setAttribute("method", "POST");
        documentForm.setAttribute("onload", "");
        var hamburgType = document.createElement("input");
        hamburgType.setAttribute("type", "text");
        hamburgType.setAttribute("id", "hamburgType");
        hamburgType.setAttribute("name", "hamburgType");
        hamburgType.setAttribute("maxlength", 48);
        hamburgType.setAttribute("name", "hamburgType");
        hamburgType.setAttribute("placeholder", "Hamburger Type");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("type", "text");
        condimentType.setAttribute("maxlength", 48);
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
        var docButtonInput = document.createElement("input");
        docButtonInput.setAttribute("id", "docButtonInput");
        docButtonInput.setAttribute("name", "newFile");
        docButtonInput.setAttribute("type", "file");
        var hiddenUserNum = document.createElement("input");
        hiddenUserNum.setAttribute("type", "hidden");
        hiddenUserNum.setAttribute("id", "hiddenUserNum");
        hiddenUserNum.setAttribute("name", "hiddenUserNum");
        var hiddenFoodType = document.createElement("input");
        hiddenFoodType.setAttribute("type", "hidden");
        hiddenFoodType.setAttribute("id", "hiddenFoodType");
        hiddenFoodType.setAttribute("name", "hiddenFoodType");
        var hiddenFoodNum = document.createElement("input");
        hiddenFoodNum.setAttribute("type", "hidden");
        hiddenFoodNum.setAttribute("id", "hiddenFoodNum");
        hiddenFoodNum.setAttribute("name", "hiddenFoodNum");
        var hiddenFoodAction = document.createElement("input");
        hiddenFoodAction.setAttribute("type", "hidden");
        hiddenFoodAction.setAttribute("id", "hiddenFoodAction");
        hiddenFoodAction.setAttribute("name", "hiddenFoodAction");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");
        submitButton.setAttribute("name", "submitButton");
        submitButton.innerHTML = "SUBMIT";
        /* Add input to the above form for document selection */
        documentForm.appendChild(hamburgType);
        documentForm.appendChild(condimentType);
        documentForm.appendChild(caloriesType);
        documentForm.appendChild(nameType);
        documentForm.appendChild(userIDInput);
        documentForm.appendChild(docButtonInput);
        documentForm.appendChild(hiddenUserNum);
        documentForm.appendChild(hiddenFoodType);
        documentForm.appendChild(hiddenFoodNum);
        documentForm.appendChild(hiddenFoodAction);
        documentForm.appendChild(submitButton);
        
        submitButton.addEventListener("click", function(){
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
                userIDInput.setAttribute("value", userID); 
            }
            //Give appropriate values for form to send
            hiddenFoodAction.setAttribute("value", "food_submit")
            hiddenFoodType.setAttribute("value", "HAMBURGER");
            hiddenFoodNum.setAttribute("value", ""); //Not sure if needed now...
            hiddenUserNum.setAttribute("value", "");
            //Submit the form!
            documentForm.submit();
            
        });
        theDiv.appendChild(documentForm);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else {
        console.log("Whoops, we got a problem. Wrong food choice came in: " + foodChoice);
        location.reload(true); //Reload Page
    }
}

//Hides the error div unless there is an error
window.addEventListener('DOMContentLoaded', function(){
    var errordisplayerDiv = document.getElementById("errordisplayer_div");

    if (errordisplayerDiv === null){
        //Do nothing, the div that needs hiding isn't on this page
    } else{
        errordisplayerDiv.style.display = "none"; //Hide this div from displaying
    }
});