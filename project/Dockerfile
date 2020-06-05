# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.14.1
#Maintainer info
LABEL maintainer="jbkeller0303@gmail.com"
LABEL build_date="2020-30-05"
#Set Golang Environment
ENV GO111MODULE=on
#Make the app directory where our code will live
RUN mkdir /app
# Copy the local package files to the container's workspace.
ADD . /app
#Speicify the work directory
WORKDIR /app
# Add this go mod download command to pull in any dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
#Run the go directory
# MIGHT NOT WORK: RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
# Expose port 8080 to the outside world
EXPOSE 8080
ENTRYPOINT [ "/app/databaseapp" ]
#Check off our binary executeable
# MIGHT NOT WORK: CMD ["/app/main"]