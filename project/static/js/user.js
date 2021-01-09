var openSignInWindow = false;
var openSignUpWindow = false;

//Used to control which link to send our user to
function navigateHeader(whichLink) {
    switch (whichLink) {
        case 1:
            //Go to ContactDev
            window.location.replace("/contact");
            break;
        case 2:
            //Go to Documentation
            window.location.replace("/documentation");
            break;
        case 3:
            //Go to Documentation
            window.location.replace("/mainPage");
            break;
        case 4:
            //Go to Documentation
            window.location.replace("/messageboard");
            break;
        case 5:
            //Go to Index
            window.location.replace("/");
            break;
        case 6:
            //Go to choicepage Page
            window.location.replace("/choicepage");
            break;
        default:
            console.log("Error, wrong whichLink entered: " + whichLink);
            break;
    }
}

//Set the divs to all closed so they can open correctly in the function below
window.addEventListener('DOMContentLoaded', function(){
    var bod1 = document.getElementById("bodOpen1");
    if (bod1 === null){
        //Nothing in here, not Documentation page
    } else {
        //Declare all Divs that need to start off with 'none' display value
        var bod2 = document.getElementById("bodOpen2");
        var bod3 = document.getElementById("bodOpen3");
        var bod4 = document.getElementById("bodOpen4");
        var bod5 = document.getElementById("bodOpen5");
        var bod6 = document.getElementById("bodOpen6");
        var bod7 = document.getElementById("bodOpen7");
        var bod7 = document.getElementById("bodOpen8");
        //Set all the bodys to 'none' display
        bod1.style.display = "none";
        bod2.style.display = "none";
        bod3.style.display = "none";
        bod4.style.display = "none";
        bod5.style.display = "none";
        bod6.style.display = "none";
        bod7.style.display = "none";
        bod8.style.display = "none";
        console.log("DEBUG: We should be hiding the divs.");
    }
});

//Open the correctDivs when clicked
function documentDivDisplay(whichDiv){
    console.log("DEBUG: should be displaying this link: " + whichDiv);
    switch(whichDiv){
        case 1:
            //Display or not display 1st Div
            var theDiv = document.getElementById("bodOpen1");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 2:
            //Display or not display 2nd Div
            var theDiv = document.getElementById("bodOpen2");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 3:
            //Display or not display 3rd Div
            var theDiv = document.getElementById("bodOpen3");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 4:
            //Display or not display 4th Div
            var theDiv = document.getElementById("bodOpen4");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 5:
            //Display or not display 5th Div
            var theDiv = document.getElementById("bodOpen5");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 6:
            //Display or not display 6th Div
            var theDiv = document.getElementById("bodOpen6");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 7:
            //Display or not display 6th Div
            var theDiv = document.getElementById("bodOpen7");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        case 8:
            //Display or not display 6th Div
            var theDiv = document.getElementById("bodOpen8");
            if (theDiv.style.display === "none"){
                theDiv.style.display = "flex";
            } else {
                theDiv.style.display = "none";
            }
            break;
        default:
            console.log("Error, incorrect div was opened.");
            break;
    }
}

//Close Ajax windows after submission for page reloads
window.addEventListener('DOMContentLoaded', function(){
    var signUp = document.getElementById("signup-ask-text");

        if (signUp === null){
            //Do nothing, it isn't on this page
        } else {
            //Declare form variables and make sure they are closed
            openSignInWindow = false;
            openSignUpWindow = false;
            var divformDivSignUp = document.getElementById("divformDivSignUp");
            divformDivSignUp.style = "display: none";
            var divformDivLogin = document.getElementById("divformDivLogin");
            divformDivLogin.style = "display: none";
        }
});

//Testing stuff
function testFormSubmit(){
    var theForm = document.getElementById('DEBUGpostForm');
    theForm.submit();
}