package peerutil

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/bellent69ne/ratnet/ratp"
	"log"
)

//func GetFriends() []byte {
func DoSomething() {
	endpoint := "127.0.0.1:1366" //+ //string(ratp.PORT)
	var curSession ratp.Session
	err := curSession.Connect(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	// close session when finished
	defer curSession.Conn.Close()

	// Make initial handshaking for current session
	if !ratp.Handshake(&curSession) {
		//return nil
		return
	}

	// Make session secure
	if !ratp.SecureSession(&curSession) {
		return
	}

	addrs := ratp.GetFriendsAddrs(&curSession)
	fmt.Println(addrs)
}
