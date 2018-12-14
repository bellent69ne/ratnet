package peerutil

import (
	"github.com/bellent69ne/ratnet/ratp"
	"log"
)

func DoSomething() {
	endpoint := "127.0.0.1:1366" //+ //string(ratp.PORT)
	var newSession ratp.Session
	err := newSession.Connect(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	newSession.HelloFriend()
}
