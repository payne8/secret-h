package main

import (
	"crypto/rand"
	"fmt"
	sh "github.com/murphysean/secrethitler"
	"net/http"
	"os"
)

var theGame *sh.SecretHitler
var theGameFile string

type Writer struct {
	Name string
}

func (w Writer) Write(b []byte) (int, error) {
	f, err := os.OpenFile(w.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("write:", theGameFile, err)
		return 0, err
	}
	defer f.Close()

	return f.Write(b)
}

func main() {
	//Specify a file to write all the events to
	theGame = sh.NewSecretHitler()
	var w Writer
	theGameFile = GenUUIDv4() + ".json"
	w.Name = theGameFile
	theGame.Log = w

	http.HandleFunc("/api/state", APIStateHandler)
	http.HandleFunc("/api/event", APIEventHandler)
	http.HandleFunc("/sse", ServerSentEventsHandler)
	//A file handler for the static assets
	http.Handle("/", http.FileServer(http.Dir("www/dist")))

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
