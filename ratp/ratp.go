package ratp

import (
	"crypto/rsa"
	"encoding/gob"
	"net"
)

// Message types
const (
	// initial handshaking, can have public key with it
	MsgHelloFriend = "/* hell0 fri3nd */\n"
	// gift means aes encryption key encrypted with recievers public key
	MsgIHaveAGift = "/* I have a gift */\n"
	// receiver received the aes key, everything is ok
	MsgIAppreciateThat = "/* I appreciate that */\n"
	// asking server for range of ip addresses in the network
	MsgINeedData = "/* I need data */\n"
	// Server sent range of ip addresses, means ok
	MsgYoureWelcome = "/* You're welcome */\n"
	// request for peers to create a chain of hosts(network)
	MsgINeedFriends = "/* I need fri3nds */\n"
	PORT            = 1366
)

const (
	ErrDontHaveData = "/* Can't help */"
	Err
)

type Session struct {
	conn        net.Conn
	privateKey  *rsa.PrivateKey
	aesKey      []byte
	alienPubKey *rsa.PublicKey
}

// HelloFriend - makes initial handshaking between peers,
// can have public key attached
func (self *Session) HelloFriend() {
	var newParcel Parcel
	if self.privateKey != nil {
		newParcel = createParcel(MsgHelloFriend, self.privateKey.PublicKey)
	}
	newParcel = createParcel(MsgHelloFriend, nil)
	SendParcel(self.conn, &newParcel)

}

// IHaveAGift - sends encrypted aes key for futher communication
// has to have encrypted aes key attached
func (self *Session) IHaveAGift() bool {
	// you should encrypt aesKey with alien public key
	newParcel := createParcel(MsgIHaveAGift, self.aesKey)

	SendParcel(self.conn, &newParcel)

	answer := ReceiveParcel(self.conn)

	if answer.Message == MsgIAppreciateThat {
		return true
	}

	return false
}

func (self *Session) IAppreciateThat() {
	newParcel := createParcel(MsgIAppreciateThat, nil)

	SendParcel(self.conn, &newParcel)
}

type Parcel struct {
	Message    string
	Attachment interface{}
}

func (self *Session) Connect(endpoint string) (err error) {
	self.conn, err = net.Dial("tcp", endpoint)
	if err != nil {
		return err
	}

	return nil
}

func createParcel(msgType string, attachment interface{}) Parcel {
	return Parcel{msgType, attachment}
}

// SendParcel - sends parcel to the connection
func SendParcel(conn net.Conn, parcel *Parcel) {
	// Create package to send over network
	//parcelToSend := createParcel(MsgHelloFriend, *publicKey)
	parcelEncoder := gob.NewEncoder(conn)

	parcelEncoder.Encode(*parcel)
}

// ReceiveParcel - receives parcel from the connection
func ReceiveParcel(conn net.Conn) (receivedParcel Parcel) {
	parcelDecoder := gob.NewDecoder(conn)
	parcelDecoder.Decode(&receivedParcel)

	return
}
