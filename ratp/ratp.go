package ratp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func Serve() {
	fmt.Println("ratnet server is running...")
	fmt.Println("waiting for connection...")

	ln, err := net.Listen("tcp", ":1366")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	dataRead, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(dataRead))

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
