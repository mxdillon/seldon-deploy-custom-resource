FROM golang:1.15

WORKDIR /go/src/app
COPY . .

# Download and install required packages
RUN go get -d -v ./..
RUN go install -v ./..

CMD ["app"]
