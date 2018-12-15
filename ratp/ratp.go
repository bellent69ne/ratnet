package ratp

import (
	//	"crypto/cipher"
	"crypto/rsa"
	"encoding/gob"
	//	"encoding/json"
	//	"errors"
	//"fmt"
	"github.com/bellent69ne/ratnet/ratcrypt"
	//"io/ioutil"
	//"log"
	"net"
)

// Message types
const (
	// initial handshaking, can have public key with it
	MsgHelloFriend = "/* hell0 fri3nd */\n"
	// gift means aes encryption key encrypted with recievers public key
	MsgHaveAGift = "/* I have a gift */\n"
	// receiver received the aes key, everything is ok
	MsgAppreciate = "/* I appreciate that */\n"
	// asking server for range of ip addresses in the network
	MsgNeedFriends = "/* I need fri3nds */\n"
	// Server sent range of ip addresses, means ok
	MsgYoureWelcome = "/* You're welcome */\n"
	// request for peers to create a chain of hosts(network)
	MsgBeFriends   = "/* Be fri3nds */\n"
	MsgWereFriends = "/* We're friends */\n"
	MsgData        = "/* Your data */\n"
	PORT           = 1366
)

const (
	ErrCantHelp       = "/* Can't help */\n"
	ErrDontUnderstand = "/* I don't understand you */\n"
)

type Session struct {
	conn        net.Conn
	privateKey  *rsa.PrivateKey
	aesKey      []byte
	alienPubKey *rsa.PublicKey
}

// PublicKey - returns public key for session
func (self *Session) PublicKey() *rsa.PublicKey {
	return &self.privateKey.PublicKey
}

// AesKey - returns aes encryption for session
func (self *Session) AesKey() []byte {
	return self.aesKey
}

// GenerateRSAkey - generates rsa key pair for session
func (self *Session) GenerateRSAkey() error {
	private, err := ratcrypt.GenerateRSAkey()
	if err != nil {
		return err
	}

	self.privateKey = private
	return nil
}

// GenerateAeskey - generates aes key for the session
func (self *Session) GenerateAESkey() error {
	key, err := ratcrypt.GenerateAESkey()
	if err != nil {
		return err
	}

	self.aesKey = key
	return nil
}

// Connect - connects to the specified endpoint
func (self *Session) Connect(endpoint string) (err error) {
	self.conn, err = net.Dial("tcp", endpoint)
	if err != nil {
		return err
	}

	return nil
}

// SendParcel - sends parcel to the connection
func (self *Session) SendParcel(parcel *Parcel) {
	parcelEncoder := gob.NewEncoder(self.conn)

	parcelEncoder.Encode(*parcel)
}

// ReceiveParcel - receives parcel from the connection
func ReceiveParcel(conn net.Conn) (receivedParcel Parcel) {
	parcelDecoder := gob.NewDecoder(conn)
	parcelDecoder.Decode(&receivedParcel)

	return
}

/*func SendParcel(conn net.Conn, parcel *Parcel) {
	data, err := json.Marshal(*parcel)

	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func ReceiveParcel(conn net.Conn) (receivedParcel Parcel) {
	received, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	var newParcel Parcel

	err = json.Unmarshal(received, &newParcel)
	if err != nil {
		log.Fatal(err)
	}

	return newParcel
}*/
