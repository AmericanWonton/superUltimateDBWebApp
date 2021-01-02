var openSignInWindow = false;

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
            //Go to signup Page
            window.location.replace("/signup");
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
    var signUp = document.getElementById("signin-ask-text");
    if (signUp === null){
        //Do nothing, this isn't on this page
    } else {
        //Declare the variables on the window
        var divForm = document.getElementById("divform");
        divForm.style = "display: none";
        //Listen for the button click
        signUp.addEventListener("click", function(){
            
            if (openSignInWindow === false){
                divForm.style = "display: flex";
                divForm.style = "flex-flow: wrap";
                divForm.style = "align-content: center";
                divForm.style = "justify-content: center";
                divForm.style = "width: 100%";
                divForm.style = "padding: 1rem";
                openSignInWindow = true;
            } else {
                divForm.style = "display: none";
                openSignInWindow = false;
            }
            console.log("openSignIn is: " + openSignInWindow);
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