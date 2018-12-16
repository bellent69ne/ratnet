package peerutil

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/bellent69ne/ratnet/ratp"
	"log"
)

func DoSomething() {
	endpoint := "127.0.0.1:1366" //+ //string(ratp.PORT)
	var curSession ratp.Session
	err := curSession.Connect(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	err = curSession.GenerateRSAkey()
	if err != nil {
		log.Fatal(err)
	}
	err = curSession.GenerateAESkey()
	curParcel, err := ratp.NewParcel(ratp.MsgHelloFriend, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = curSession.SendParcel(&curParcel)
	if err != nil {
		log.Fatal(err)
	}
	gotParcel, err := curSession.ReceiveParcel()
	if err != nil {
		log.Fatal(err)
	}

	var alien rsa.PublicKey
	err = json.Unmarshal(&alien)
	if err != nil {
		log.Fatal(err)
	}
	curSession.SetAlienKey(&alien)

	err = curSession.GenerateAESkey()
	if err != nil {
		log.Fatal(err)
	}

	curParcel, err = ratp.NewParcel(ratp.HaveAGift, curSession.AesKey())
	if err != nil {
		log.Fatal(err)
	}

	err = curSession.SendParcel(curParcel)

}
