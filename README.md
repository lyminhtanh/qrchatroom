# Welcome to QR Chatroom

##User flow:


User loads application on computer



User requests a unique URL QR code



User scans QR code with phone to connect
the two devices via a unique URL



User can now submit messages on either the
computer or phone and access the message on both devices



##Considerations:


Multiple users should be able to use the
application at the same time



The URL should be easy to type or contain
a shortened version for ease of access without using the QR code



Extra credit
(we don’t expect you to do all of these, but pick a few you’re interested in to show off your skills):


Make messages live update, user doesn’t have
to refresh to see new messages



Make messages ephemeral, set to self destruct
after a certain period of time



Add the ability to send files of arbitrary
data



Containerize your application using Docker.



Use Golang (make sure to support go modules)
or typescript

### Start the web server:

   revel run -a chatroom

### Go to http://localhost:9000/ and you'll see the home page



## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


## Help

##
