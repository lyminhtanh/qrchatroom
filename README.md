# Welcome to QR Chatroom

##User flow:

1. User loads application on computer -> YES

2. User requests a unique URL QR code -> YES

3. User scans QR code with phone to connect -> YES
the two devices via a unique URL

4. User can now submit messages on either the -> YES
computer or phone and access the message on both devices

##Considerations:

1. Multiple users should be able to use the -> YES
application at the same time

2. The URL should be easy to type or contain -> YES /room/<roomName>
a shortened version for ease of access without using the QR code

## Extra credit

1. Make messages live update, user doesn’t have -> YES (Websocket)
to refresh to see new messages
2. Make messages ephemeral, set to self destruct -> NOT YET
after a certain period of time
3. Add the ability to send files of arbitrary data -> NOT YET
4. Containerize your application using Docker. -> YES (Docker)
5. Use Golang (make sure to support go modules) or typescript-> YES (Golang with Revel FW)

### Technologies
- Golang Revel Web Framework 
- Websocket
- Docker
- Google app engine to deploy dockerized application (app.yaml config)
- Google Cloud Storage to store QR code image
- SQLite3 database

## Next Steps
1. Apply Dependency Injection (google wire project) set up DB type connection and cloud provider

2. Responsive layout for better displaying on Mobile device

3. Send/ receive media objects

4. Caching to decrease number of DB connection and reduce latency messaging using Redis

5. Handle fault tolerance using Hystryx

6. Fix some existing bugs

### Start the web server on local:

 1. Set up environment variables:
     HTTP_ADDR, 
     GOOGLE_PROJECT_ID, 
     GOOGLE_BUCKET, 
     GOOGLE_APPLICATION_CREDENTIALS (no needed on production)
     
 2. revel run -a chatroom


 3. Go to http://<YOUR_IP>:8080/ and you'll see the home page