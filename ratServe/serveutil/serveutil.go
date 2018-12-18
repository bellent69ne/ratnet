package serveutil

import (
	//"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/bellent69ne/ratnet/ratp"
	"log"
	"net"
)

const (
	friendsNum = 3
)

func Serve() {
	fmt.Println("Greetings from ratnet server :)...")
	ln, err := net.Listen("tcp", ":1366") //+string(ratp.PORT))
	if err != nil {
		log.Fatal(err)
	}

	addresses := make([]string, 0)
	addrChan := make(chan string)

	go func() {
		for {
			addr := <-addrChan
			//addr = justIpAddr(addr)
			if addrDoesntExist(addr, addresses) {
				addresses = append(addresses, addr)
			}
		}
	}()

	for {
		var newSession ratp.Session
		newSession.Conn, err = ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(addresses)
		go handleSession(&newSession, addrChan, addresses)
	}
}

func addrDoesntExist(addr string, addresses []string) bool {
	for _, val := range addresses {
		if val == addr {
			return false
		}
	}

	return true
}

func justIpAddr(addr string) string {
	var tmp []byte
	for i, val := range addr {
		if val == ':' {
			tmp = []byte(addr)
			tmp = tmp[:i]
			break
		}
	}

	return string(tmp)
}

func handleSession(curSession *ratp.Session, addrChan chan string,
	addrs []string) {
	// Close session when finished
	defer curSession.Conn.Close()
	// Receive parcel from the peer
	parcel, err := curSession.ReceiveParcel()
	// if couldn't receive parcel
	if err != nil {
		// log why, and close session
		log.Println(err)
		return
	}
	// print contents
	printParcel(&parcel)
	// if that parcel is not "hello fri3nd" handshaking
	if !isHelloFriend(&parcel, curSession, addrChan) {
		// have nothing to do, close session
		return
	}

	// Generate RSA key for this session
	err = curSession.GenerateRSAkey()
	// if couldn't generate rsa key for this session
	if err != nil {
		// log why, close session
		log.Println(err)
		return
	}

	// if we couldn't answer to "hello fri3nd" handshaking
	// have nothing to do, close session
	if !toldHelloFriend(curSession) {
		return
	}

	// Receive parcel, should receive "I have a gift"
	// with associated aes key
	parcel, err = curSession.ReceiveParcel()
	// if couldn't receive parcel
	if err != nil {
		// log why, close session
		log.Println(err)
		return
	}
	// print the contents fo the parcel
	printParcel(&parcel)

	// if received parcel is not "I have a gift"
	// nothing to do, close session
	if !isHaveAGift(curSession, &parcel) {
		return
	}

	// Receive new parcel, should receive "I need fri3nds"
	parcel, err = curSession.ReceiveParcel()
	// if couldn't receive any parcel
	if err != nil {
		// log why, close session
		log.Println(err)
		return
	}
	// print the contents of the parcel
	printParcel(&parcel)
	// if received parcel is not "I need fri3nds"
	// nothing to do, close session
	if !isNeedFriends(&parcel) {
		return
	}

	// Send ip addresses of potential fri3nds
	err = sendFriends(curSession, addrs)
	// if couldn't send potential fri3nds
	if err != nil {
		// log why, close session
		log.Println(err)
		return
	}
}

func printParcel(newParcel *ratp.Parcel) {
	fmt.Println(string(newParcel.Message))
	fmt.Println(newParcel.Attachment)
}

func isHelloFriend(newParcel *ratp.Parcel, curSession *ratp.Session,
	addrChan chan string) bool {
	if string(newParcel.Message) != ratp.MsgHelloFriend {

		justFuckOf, err := ratp.NewParcel(ratp.ErrDontUnderstand, nil)
		if err != nil {
			log.Println(err)
			return false
		}
		err = curSession.SendParcel(&justFuckOf)
		if err != nil {
			log.Println(err)
		}
		return false
	}

	addrChan <- curSession.Conn.RemoteAddr().String()

	return true
}

func toldHelloFriend(newSession *ratp.Session) bool {
	greeting, err := ratp.NewParcel(ratp.MsgHelloFriend, newSession.PublicKey())
	if err != nil {
		log.Println(err)
		return false
	}
	err = newSession.SendParcel(&greeting)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func isHaveAGift(curSession *ratp.Session, parcel *ratp.Parcel) bool {
	if string(parcel.Message) == ratp.MsgHaveAGift {
		curSession.SetAES(parcel.Attachment)
		newParcel, err := ratp.NewParcel(ratp.MsgAppreciate, nil)
		if err != nil {
			log.Println(err)
			return false
		}
		err = curSession.SendParcel(&newParcel)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}

	return false
}

func isNeedFriends(parcel *ratp.Parcel) bool {
	if string(parcel.Message) == ratp.MsgNeedFriends {
		return true
	}

	return false
}

func sendFriends(curSession *ratp.Session, addrs []string) error {
	friendsSlice := make([]string, 0)
	num := 0
	//curAlienIp := justIpAddr(curSession.Conn.RemoteAddr().String())
	for _, val := range addrs {
		if num == friendsNum {
			break
		}
		//	if val == curAlienIp {
		if val == curSession.Conn.RemoteAddr().String() {
			continue
		}
		friendsSlice = append(friendsSlice, val)
		num++
	}

	var err error
	if len(friendsSlice) == 0 {
		err = makeAndSend(curSession, ratp.ErrCantHelp, nil)
	} else {
		data, err := encodeFriends(friendsSlice)
		if err != nil {
			return err
		}
		err = makeAndSend(curSession, ratp.MsgYoureWelcome, data)
	}

	return err
}

func makeAndSend(curSession *ratp.Session, message string, attachment interface{}) error {
	newParcel, err := ratp.NewParcel(message, attachment)
	if err != nil {
		return err
	}
	err = curSession.SendParcel(&newParcel)
	if err != nil {
		return err
	}

	return nil
}

func encodeFriends(friendsSlice []string) ([]byte, error) {
	data, err := json.Marshal(friendsSlice)
	if err != nil {
		return nil, err
	}

	return data, nil
}
