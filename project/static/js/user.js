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
            //Go to Index
            window.location.replace("/");
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

//Listen for divs to hide on the documentation page
window.addEventListener('DOMContentLoaded', function(){
    var testDiv = document.getElementById("hedOpen1");

    if (testDiv === null){
        //Do nothing, this page dosen't have divs to hide
        console.log("DEBUG: No documentation page...");
    } else{
        console.log("DEBUG: Documentation Page...");
        var divs = document.getElementsByClassName("body_opening_div");
        var j;
        for (j = 0; j < divs.length; j++){
            divs[j].display = "none";
        }

        
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