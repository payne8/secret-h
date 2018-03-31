package main

import (
	"crypto/rand"
	"fmt"
	sh "github.com/murphysean/secrethitler"
	"net/http"
	"os"
)

var theGame *sh.SecretHitler

func main() {
	//Specify a file to write all the events to
	theGame = sh.NewSecretHitler()
	var err error
	theGame.Log, err = os.OpenFile("log.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer theGame.Log.Close()

	http.HandleFunc("/api/state", APIStateHandler)
	http.HandleFunc("/api/event", APIEventHandler)
	http.HandleFunc("/sse", ServerSentEventsHandler)
	//A file handler for the static assets
	http.Handle("/", http.FileServer(http.Dir("www")))

	http.ListenAndServe(":8080", nil)
}

func GenUUIDv4() string {
	u := make([]byte, 16)
	rand.Read(u)
	//Set the version to 4
	u[6] = (u[6] | 0x40) & 0x4F
	u[8] = (u[8] | 0x80) & 0xBF
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
