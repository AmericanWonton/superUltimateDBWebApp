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

//Listen for button to submit email
window.addEventListener('DOMContentLoaded', function(){
    var button = document.getElementById("submitB");

    if (button === null){
        //Do nothing...this is for the 'Contact' page
    } else {
        button.addEventListener("click", function(){
            var name = document.getElementById("YourNameInput");
            var email = document.getElementById("YourEmailInput");
            var message = document.getElementById("YourMessageInput");
    
            var nameVal = String(name.value);
            var emailVal = String(email.value);
            var messageVal = String(message.value);
    
            console.log("The name is " + nameVal);
            console.log("The email is " + emailVal);
            console.log("The message is " + messageVal);
    
            //This is executed in Ajax
            function goToHomeScreen(){
                window.location.replace("/");
            }
    
    
            var UserJSON = {
                TheName: nameVal,
                TheEmail: emailVal,
                TheMessage:  messageVal
            };
            var jsonString = JSON.stringify(UserJSON);
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/contact', true);
            xhr.setRequestHeader("Content-Type", "application/json");
            xhr.addEventListener('readystatechange', function(){
                if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                    var item = xhr.responseText;
                    var successMSG = JSON.parse(item);
                    if (successMSG.SuccessNum == 0){
                        //Create Text to inform User
                        var informDiv = document.getElementById("informResultDiv");
                        var informP = document.getElementById("informTxtP");
                        informP.innerHTML = "Message submitted! Thanks, I'll take a look!";
                        informDiv.display = "block";
                        name.innerHTML = "";
                        email.innerHTML = "";
                        message.innerHTML = "";
                        //Go back to homescreen after a few seconds
                        setTimeout(goToHomeScreen(), 3000);
                    } else {
                        //Create Text to inform User
                        var informDiv = document.getElementById("informResultDiv");
                        var informP = document.getElementById("informTxtP");
                        informP.innerHTML = "There was an issue sending your message...";
                        informDiv.display = "block";
                        name.innerHTML = "";
                        email.innerHTML = "";
                        message.innerHTML = "";
                    }
                }
            });
            xhr.send(jsonString);
        });
    }
});

//Listen for User to click the Sign In button
window.addEventListener('DOMContentLoaded', function(){
    var signIn = document.getElementById("signin-ask-text");
    if (signIn === null){
        //Do nothing, this isn't on this page
    } else {
        //Declare the variables on the window
        var divformDivLogin = document.getElementById("divformDivLogin");
        divformDivLogin.style = "display: none";
        //Listen for the button click
        signIn.addEventListener("click", function(){
            //First check to see if other sheet is open; if yes, close it
            var divformDivSignUp = document.getElementById("divformDivSignUp");
            if (openSignUpWindow === true) {
                openSignUpWindow = false;
                divformDivSignUp.style = "display: none";
            }
            //Open this form if needed
            if (openSignInWindow === false){
                divformDivLogin.style = "display: flex";
                divformDivLogin.style = "flex-flow: wrap";
                divformDivLogin.style = "align-content: center";
                divformDivLogin.style = "justify-content: center";
                divformDivLogin.style = "width: 100%";
                divformDivLogin.style = "padding: 1rem";
                openSignInWindow = true;
            } else {
                divformDivLogin.style = "display: none";
                openSignInWindow = false;
            }
        });
    }
});

//Listen for User to click the Sign Up button
window.addEventListener('DOMContentLoaded', function(){
    var signUp = document.getElementById("signup-ask-text");

    if (signUp === null){
        //Do nothing, it isn't on this page
    } else {
        //Declare the variables on the window
        var divformDivSignUp = document.getElementById("divformDivSignUp");
        divformDivSignUp.style = "display: none";
        //Listen for the button click
        signUp.addEventListener("click", function(){
            //First check to see if other sheet is open; if yes, close it
            var divformDivLogin = document.getElementById("divformDivLogin");
            if (openSignInWindow === true) {
                openSignInWindow = false;
                divformDivLogin.style = "display: none";
            }
            //open this sheet if needed
            if (openSignUpWindow === false){
                divformDivSignUp.style = "display: flex";
                divformDivSignUp.style = "flex-flow: wrap";
                divformDivSignUp.style = "align-content: center";
                divformDivSignUp.style = "justify-content: center";
                divformDivSignUp.style = "width: 100%";
                divformDivSignUp.style = "padding: 1rem";
                openSignUpWindow = true;
            } else {
                divformDivLogin.style = "display: none";
                openSignUpWindow = false;
            }
        });
    }
});

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


//Testing stuff
function testFormSubmit(){
    var theForm = document.getElementById('DEBUGpostForm');
    theForm.submit();
}