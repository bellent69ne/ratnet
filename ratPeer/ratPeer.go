package main

import (
	//"fmt"
	"github.com/bellent69ne/ratnet/ratPeer/peerutil"
	//"net"
)

func writeMsg(endpoint string) {
	/*	conn, err := net.Dial("tcp", endpoint)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(conn, "Assfuck")
	*/
}

func main() {
	//wrieMsg("192.168.43.5:1366")
	peerutil.DoSomething()
}
