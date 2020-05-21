  
# Use the official go docker image built on debian.
FROM golang:1.14.2

# Grab the source code and add it to the workspace.
ADD . /go/src/chatroom

# Install revel and the revel CLI.
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel

# Use the revel CLI to start up our application.
ENTRYPOINT revel run chatroom prod
# Open up the port where the app is running.
EXPOSE 8080