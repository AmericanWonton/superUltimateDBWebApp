function revealFoodForm(foodChoice) {
    var theDiv = document.getElementById("foodFormDiv"); //Get the div
    theDiv.innerHTML = ""; //Remove any child elements if any remain
    var theForm = document.createElement("form");
    theForm.setAttribute("method", "POST");
    //var theForm = document.get
    /* For Hotdog selection */
    if (foodChoice == 0) {
        console.log("Mkay, you clikced hotdog.")
        //Create the form elements and append them to the form
        var user = document.createElement("input");
        user.setAttribute("id", "userTxt");
        user.setAttribute("display", "hidden");
        user.setAttribute("name", "userTXT");
        user.setAttribute("placeholder", "userTXT");
        var hDogType = document.createElement("input");
        hDogType.setAttribute("id", "hotDogType");
        hDogType.setAttribute("name", "hDogType");
        hDogType.setAttribute("placeholder", "hDogType");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("placeholder", "condimentType");
        var caloriesType = document.createElement("input");
        caloriesType.setAttribute("id", "caloriesType");
        caloriesType.setAttribute("name", "caloriesType");
        caloriesType.setAttribute("placeholder", "caloriesType");
        var nameType = document.createElement("input");
        nameType.setAttribute("id", "nameType");
        nameType.setAttribute("name", "nameType");
        nameType.setAttribute("placeholder", "nameType");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");

        theForm.appendChild(user);
        theForm.appendChild(hDogType);
        theForm.appendChild(condimentType);
        theForm.appendChild(caloriesType);
        theForm.appendChild(nameType);
        theForm.appendChild(submitButton);

        theDiv.appendChild(theForm);
        //Dispaly the form to click on
        theDiv.style.display = "block";



    } else if (foodChoice == 1) {
        /* For Hamburger Selection */
        console.log("Mkay, you clikced hamburger.")
        
        var user = document.createElement("input");
        user.setAttribute("id", "userTxt");
        user.setAttribute("display", "hidden");
        user.setAttribute("name", "userTXT");
        user.setAttribute("placeholder", "userTXT");
        var hamBurgType = document.createElement("input");
        hamBurgType.setAttribute("id", "hamBurgType");
        hamBurgType.setAttribute("name", "hamBurgType");
        hamBurgType.setAttribute("placeholder", "hamBurgType");
        var condimentType = document.createElement("input");
        condimentType.setAttribute("id", "condimentType");
        condimentType.setAttribute("name", "condimentType");
        condimentType.setAttribute("placeholder", "condimentType");
        var caloriesType = document.createElement("input");
        caloriesType.setAttribute("id", "caloriesType");
        caloriesType.setAttribute("name", "caloriesType");
        caloriesType.setAttribute("placeholder", "caloriesType");
        var nameType = document.createElement("input");
        nameType.setAttribute("id", "nameType");
        nameType.setAttribute("name", "nameType");
        nameType.setAttribute("placeholder", "nameType");
        var submitButton = document.createElement("button");
        submitButton.setAttribute("id", "submitButton");

        theForm.appendChild(user);
        theForm.appendChild(hamBurgType);
        theForm.appendChild(condimentType);
        theForm.appendChild(caloriesType);
        theForm.appendChild(nameType);
        theForm.appendChild(submitButton);

        theDiv.appendChild(theForm);
        //Dispaly the form to click on
        theDiv.style.display = "block";
    } else {
        console.log("Whoops, we got a problem. Wrong food choice came in.")
    }
}