package main

import (
	"fmt"
	"log"
	"net"
)

func writeMsg(endpoint string) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(conn, "Assfuck")
}

func main() {
	writeMsg("192.168.43.5:1366")
}
