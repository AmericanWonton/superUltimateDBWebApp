function pictureSubmit(theFileForm){
    console.log("Submitting food photo.");
    theFileForm.submit();
}

function srcCleaner(theSrcString){
    var newString = "";
    if (theSrcString.includes("\\")){
        console.log("Replacing this string: " + theSrcString);
        
    }
}

function fileUploadCheck(){
    var hasFile = false;
    if (document.getElementById("photo-update-input").value != ""){
        //We have a file, return true
        hasFile = true;
    }else {
        //User does not wish to update their file, return false
        hasFile = false;
    }
    return hasFile;
}

function photoUpdate(theFoodUpdate){
    //Make JSON to send to food photo API
    var picFoodUpdate = {
        PhotoFoodType: theFoodUpdate.FoodType,
        PhotoFoodID: theFoodUpdate.FoodID,
        PhotoTheHamburger: theFoodUpdate.TheHamburger,
        PhotoTheHotDog: theFoodUpdate.TheHotDog
    };

    var theFoodPicForm = document.getElementById("photo-update-form");
    //Set the foodID variables
    if (picFoodUpdate.PhotoFoodType == "hotdog"){
        //Get hotdog information and enter it into form variables
        document.getElementById("hiddenFoodID").value = Number(picFoodUpdate.PhotoFoodID);
        document.getElementById("hiddenUserID").value = Number(picFoodUpdate.PhotoTheHotDog.UserID);
        document.getElementById("hiddenCurrentSRC").value = picFoodUpdate.PhotoTheHotDog.PhotoSrc;
        document.getElementById("hiddenFoodType").value = "hotdog";
        document.getElementById("hiddenPhotoID").value = Number(picFoodUpdate.PhotoTheHotDog.PhotoID);
    } else if (picFoodUpdate.PhotoFoodType == "hamburger"){
        //Get Hamburger information and enter it into form variables
        document.getElementById("hiddenFoodID").value = Number(picFoodUpdate.PhotoFoodID);
        document.getElementById("hiddenUserID").value = Number(picFoodUpdate.PhotoTheHamburger.UserID);
        document.getElementById("hiddenCurrentSRC").value = picFoodUpdate.PhotoTheHamburger.PhotoSrc;
        document.getElementById("hiddenFoodType").value = "hamburger";
        document.getElementById("hiddenPhotoID").value = Number(picFoodUpdate.PhotoTheHamburger.PhotoID);
    } else {
        console.log("Error updating foodID and User variables: " + picFoodUpdate.PhotoFoodType);
    }
    theFoodPicForm.submit();
}