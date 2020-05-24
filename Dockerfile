  
# Use the official go docker image built on debian.
FROM golang:1.14.2

# Grab the source code and add it to the workspace.
ADD . /go/src/chatroom

# Install revel and the revel CLI.
RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y sqlite3 libsqlite3-dev

RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel
RUN go get github.com/skip2/go-qrcode
RUN go get github.com/liyue201/goqr
RUN go get -u cloud.google.com/go/storage
# ORM
RUN go get -u github.com/jinzhu/gorm

RUN mkdir /db
RUN /usr/bin/sqlite3 /db/prod.db
CMD /bin/bash

# Use the revel CLI to start up our application.
ENTRYPOINT revel run chatroom prod
# Open up the port where the app is running.
EXPOSE 8080
EXPOSE 65080
#VOLUME D:/sqlite/data/:/db
VOLUME home/sqlite/data/:/db