FROM golang:1.10

#Compile the node application, put it in the www dir

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v

EXPOSE 8080

CMD ["app"]
