package peerutil

import (
	"encoding/json"
	"errors"
	"github.com/bellent69ne/ratnet/ratPeer/copycat"
	"github.com/bellent69ne/ratnet/ratServe/serveutil"
	"github.com/bellent69ne/ratnet/ratp"
	"log"
	"net"
)

func Serve(serverAddr string) {
	ln, err := net.Listen("tcp", ratp.PeerPort)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var newSession ratp.Session
		newSession.Conn, err = ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleSession(&newSession, serverAddr)
	}

}

func handleSession(curSession *ratp.Session, serverAddr string) {
	defer curSession.Conn.Close()

	//parcel, err := curSession.ReceiveParcel()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	err := serveutil.InitiateSession(curSession, nil)
	if err != nil {
		log.Println(err)
		return
	}

	chainAddrs, err := receiveFriends(curSession)
	if err != nil {
		log.Println(err)
		return
	}
	err = ratp.SayWereFriends(curSession)
	if err != nil {
		log.Println(err)
		return
	}

	if len(chainAddrs)+1 < ratp.NODES {
		beLikeARouter(curSession, chainAddrs, serverAddr)
	} else {
		beTheLooter(curSession)
		err = ratp.SayDone(curSession)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func beTheLooter(curSession *ratp.Session) {
	data, err := ratp.ReceiveData(curSession)
	if err != nil {
		log.Println(err)
		return
	}

	url := string(data)
	size, err := copycat.Inspect(&url)

	if err != nil {
		log.Println(err)
	}

	// something wrong here
	stream := make(chan []byte)
	var nextChunk int64

	// Should reconsider how to output filename
	for nextChunk != int64(size) {
		go copycat.LootChunk(&url, nextChunk, stream)

		received := <-stream
		parcel, err := ratp.NewParcel(ratp.MsgData, received)
		err = curSession.SendParcel(&parcel)
		if err != nil {
			log.Println(err)
			return
		}

		nextChunk += copycat.CHUNK

		if nextChunk > int64(size) {
			nextChunk = int64(size)
		}
	}
}

func beLikeARouter(curSession *ratp.Session, frIpAddrs []string, serverAddr string) error {
	serverSes := ServerSession(serverAddr)
	if serverSes == nil {
		return errors.New("Couldn't create session with server")
	}
	addrs := GetAddresses(serverSes)
	if addrs == nil {
		return errors.New("Couldn't get peer addresses")
	}
	ses := getSecureSession(addrs)
	if ses == nil {
		return errors.New("Couldn't make secure session with any peers")
	}
	defer ses.Conn.Close()

	ip := justIpAddr(curSession.Conn.RemoteAddr().String())
	frIpAddrs = append(frIpAddrs, ip)

	data, err := json.Marshal(frIpAddrs)
	if err != nil {
		return err
	}

	err = ratp.SayBeFriends(ses, data)
	if err != nil {
		return err
	}

	go redirect(ses, curSession)
	redirect(curSession, ses)
	return nil
}

func justIpAddr(addr string) string {
	var data []byte
	for i, val := range addr {
		if val == ':' {
			data = []byte(addr)
			data = data[:i]
			break
		}
	}

	return string(data)
}

func redirect(sessionA, sessionB *ratp.Session) {
	for {
		parcel, err := sessionA.ReceiveParcel()
		if err != nil {
			log.Println(err)
			break
		}
		if len(parcel.Message) == 0 {
			continue
		}
		err = sessionB.SendParcel(&parcel)
		if err != nil {
			log.Println(err)
			break
		}
		if string(parcel.Message) == ratp.MsgDone {
			break
		}
	}
}

func receiveFriends(curSession *ratp.Session) ([]string, error) {
	data, err := ratp.WantsToBeFriends(curSession)
	if err != nil {
		return nil, err
	}

	var friends []string
	err = json.Unmarshal(data, &friends)
	if err != nil {
		return nil, err
	}

	return friends, nil
}
