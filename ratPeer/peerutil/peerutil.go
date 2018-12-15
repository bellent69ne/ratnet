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

	err = newSession.GenerateRSAkey()
	if err != nil {
		log.Fatal(err)
	}
	err = newSession.GenerateAESkey()
	if err != nil {
		log.Fatal(err)
	}
	var newParcel ratp.Parcel
	err = newParcel.Make(ratp.MsgData, []byte("Assfucked"))
	if err != nil {
		log.Fatal(err)
	}
	newSession.SendParcel(&newParcel)

}
