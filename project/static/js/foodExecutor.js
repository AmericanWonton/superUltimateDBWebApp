function revealFoodForm(foodChoice) {
    var theDiv = document.getElementById("foodFormDiv"); //Get the div
    //var theForm = document.get
    /* For Hotdog selection */
    if (foodChoice == 0) {
        console.log("Mkay, you clikced hotdog.")
        theDiv.style.display = "block";
    } else if (foodChoice == 1) {
        /* For Hamburger Selection */
        onsole.log("Mkay, you clikced hamburger.")
        theDiv.style.display = "block";
    } else {
        console.log("Whoops, we got a problem. Wrong food choice came in.")
    }
}