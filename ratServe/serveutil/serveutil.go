package serveutil

import (
	//"crypto/rsa"
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
			addresses = append(addresses, addr)
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

func handleSession(curSession *ratp.Session, addrChan chan string,
	addrs []string) {
	defer curSession.Conn.Close()
	parcel, err := curSession.ReceiveParcel()
	if err != nil {
		log.Println(err)
		return
	}
	if !isHelloFriend(&parcel, curSession, addrChan) {
		return
	}

	err = curSession.GenerateRSAkey()
	if err != nil {
		log.Println(err)
		return
	}

	if !toldHelloFriend(curSession) {
		return
	}

	parcel, err = curSession.ReceiveParcel()
	if err != nil {
		log.Println(err)
		return
	}

	if !isHaveAGift(&parcel) {
		return
	}
	curSession.SetAES(parcel.Attachment)

	parcel, err = curSession.ReceiveParcel()
	if err != nil {
		log.Println(err)
		return
	}

	if !isNeedFriends(&parcel) {
		return
	}

	err = sendFriends(curSession, addrs)
	if err != nil {
		log.Println(err)
		return
	}
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
		newParcel, err := ratp.NewParcel(ratp.MsgAppreciate, nil)
		if err != nil {
			log.Println(err)
			return false
		}
		err = curSession.SendParcel(newParcel)
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
	for _, val := range addrs {
		if num == friendsNum {
			break
		}
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
		err = makeAndSend(curSession, ratp.MsgYoureWelcome, addrs)
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
