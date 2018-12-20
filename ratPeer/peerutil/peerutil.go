package peerutil

import (
	"bufio"
	"fmt"
	"github.com/bellent69ne/ratnet/ratPeer/copycat"
	"github.com/bellent69ne/ratnet/ratp"
	"log"
	"os"
	"strings"
)

const (
	UpdatePEERs = "update peers\n"
	ShowPEERS   = "show peers\n"
	LOOT        = "loot\n"
)

func JustDoIt() {
	var serverSes *ratp.Session
	var serverAddr string
	for {
		fmt.Println("Server address: ")
		serverAddr := GetUserInput()
		serverSes = ServerSession(serverAddr)
		if serverSes != nil {
			break
		}
	}

	sayHello(serverSes)
	go Serve(serverAddr)

	shell(serverAddr)
}

func shell(serverAddr string) {
	var addrs []string
	for {
		fmt.Print("RATnet$ ")
		str := GetUserInput()

		switch str {
		case UpdatePEERs:
			{
				serverSes := ServerSession(serverAddr)
				if serverSes == nil {
					fmt.Println("Couldn't update peers...")
					break
				}
				addrs = GetAddresses(serverSes)
			}
		case ShowPEERS:
			printPeers(addrs)

		default:
			DoSomeStuff(str, addrs)
		}
	}
}

func DoSomeStuff(str string, addrs []string) {
	splitted := strings.Split(str, " ")

	switch splitted[0] {
	case LOOT:
		// download the file specified
		{
			if len(splitted) > 2 {
				fmt.Printf("command \"%s\" has to have only one argument",
					splitted[0])
				return
			}
			// Loot
		}
	default:
		fmt.Printf("command \"%s\" doesn't exist\n\n", splitted[0])
	}
}

func Loot(addrs []string, url string) {
	peerSes := getSecureSession(addrs)
	if peerSes != nil {
		log.Println("Couldn't initiate secure connection with existing peers")
		return
	}
	// Do this shit
	parcel, err := ratp.NewParcel(ratp.MsgData, []byte(url))
	if err != nil {
		log.Println(err)
		return
	}

	err = peerSes.SendParcel(&parcel)
	if err != nil {
		log.Println(err)
		return
	}

	stream := make(chan []byte)

	go getData(peerSes, stream)
	err = copycat.WriteToFile(copycat.Filename(&url), stream)
	if err != nil {
		log.Println(err)
		return
	}
}

func getData(peerSes *ratp.Session, stream chan []byte) {
	var doneParcel ratp.Parcel
	for string(doneParcel.Message) != ratp.MsgDone {
		parcel, err := peerSes.ReceiveParcel()
		if err != nil {
			log.Println(err)
			stream <- nil
		}

		stream <- parcel.Attachment
	}
}

func getSecureSession(addrs []string) *ratp.Session {
	var peerSession ratp.Session
	for _, val := range addrs {
		err := peerSession.Connect(val + ratp.PeerPort)
		if err != nil {
			log.Println(err)
			return nil
		}

		if !ratp.Handshake(&peerSession) {
			log.Println("Couldn't handshake with ", val)
			return nil
		}

		if !ratp.SecureSession(&peerSession) {
			log.Println("Couldn't secure session with ", val)
			return nil
		}
	}

	return &peerSession
}

func printPeers(addrs []string) {
	fmt.Println("Peers ip addresses\n")
	for _, val := range addrs {
		fmt.Println(val)
	}
}

func sayHello(serverSes *ratp.Session) {
	done := false
	for !done {
		done = ratp.SayHelloFriend(serverSes)
	}
	serverSes.Conn.Close()
}

func ServerSession(serverAddr string) *ratp.Session {
	endpoint := serverAddr + ratp.ServerPort
	var serverSession ratp.Session
	err := serverSession.Connect(endpoint)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &serverSession
}

//func GetFriends() []byte {
func GetAddresses(serverSes *ratp.Session) []string {
	// close session when finished
	defer serverSes.Conn.Close()

	// Make initial handshaking for current session
	if !ratp.Handshake(serverSes) {
		//return nil
		return nil
	}

	// Make session secure
	if !ratp.SecureSession(serverSes) {
		return nil
	}

	addrs := ratp.GetFriendsAddrs(serverSes)

	return addrs
}

// GetUserInput - gets line of input from user from stdin
func GetUserInput() string {
	in := bufio.NewReader(os.Stdin)
	str, _ := in.ReadString('\n')

	return str
}
