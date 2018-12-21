package peerutil

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bellent69ne/ratnet/ratPeer/copycat"
	"github.com/bellent69ne/ratnet/ratp"
	"log"
	"os"
	"strings"
)

const (
	UpdatePEERs = "update peers"
	ShowPEERS   = "show peers"
	LOOT        = "loot"
)

func JustDoIt() {
	var serverSes *ratp.Session
	var serverAddr string
	for {
		fmt.Println("Server address: ")
		serverAddr = GetUserInput()
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
		var serverSes *ratp.Session

		switch str {
		case UpdatePEERs:
			{
				fmt.Println(serverAddr)
				serverSes = ServerSession(serverAddr)
				if serverSes == nil {
					fmt.Println("Couldn't update peers...")
					break
				}
				addrs = GetAddresses(serverSes)
			}
		case ShowPEERS:
			printPeers(addrs)
		case "":
			continue

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
				fmt.Printf("command \"%s\" has to have only one argument\n",
					splitted[0])
				return
			}
			// Loot
			Loot(addrs, splitted[1])
		}
	default:
		fmt.Printf("command \"%s\" doesn't exist\n\n", splitted[0])
	}
}

func Loot(addrs []string, url string) {
	peerSes := getSecureSession(addrs)
	if peerSes == nil {
		log.Println("Couldn't initiate secure connection with existing peers")
		return
	}

	err := MakeFriendShip(peerSes)
	if err != nil {
		log.Println(err)
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
	fmt.Println("Sent parcel with data in Loot")

	stream := make(chan []byte)

	go getData(peerSes, stream)
	for {
		err = copycat.WriteToFile(copycat.Filename(&url), stream)
		if err != nil {
			log.Println(err)
			return
		}
	}
	// should think about sending done message
}

func MakeFriendShip(peerSes *ratp.Session) error {
	parcel, err := ratp.NewParcel(ratp.MsgBeFriends, nil)
	if err != nil {
		return err
	}
	err = peerSes.SendParcel(&parcel)
	if err != nil {
		return err
	}

	gotParcel, err := peerSes.ReceiveParcel()
	if err != nil {
		return err
	}

	if string(gotParcel.Message) != ratp.MsgWereFriends {
		return errors.New("Peer didn't verified friendship")
	}

	return nil
}

func getData(peerSes *ratp.Session, stream chan []byte) {
	//var doneParcel ratp.Parcel
	//var err error
	for { //string(doneParcel.Message) != ratp.MsgDone {
		doneParcel, err := peerSes.ReceiveParcel()
		if err != nil {
			log.Println("Are we in get data?")
			log.Println(err)
			stream <- nil
		}
		if string(doneParcel.Message) == ratp.MsgDone {
			break
		}

		log.Println("Data = ", len(doneParcel.Attachment))
		if len(doneParcel.Attachment) != 0 {
			stream <- doneParcel.Attachment
		}
	}
	stream <- nil
}

//type DataChunk struct {
//	Data    []byte
//	CanLoot bool
//}

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
	data := []byte(str)
	data = data[:len(str)-1]

	return string(data)
}
