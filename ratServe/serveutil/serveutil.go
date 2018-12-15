package serveutil

import (
	//"crypto/rsa"
	"fmt"
	"github.com/bellent69ne/ratnet/ratp"
	"log"
	"net"
)

func Serve() {
	fmt.Println("Greetings from ratnet server :)...")
	ln, err := net.Listen("tcp", ":1366") //+string(ratp.PORT))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	parcel := ratp.ReceiveParcel(conn)

	fmt.Println(parcel)

	//key, ok := parcel.Attachment.(rsa.PublicKey)
	//if ok {
	//	fmt.Println("Pub")
	//	fmt.Println(key)
	//}
}
