package ratp

import (
	//	"crypto/cipher"
	"crypto/rsa"
	"encoding/gob"
	"encoding/json"
	//"errors"
	//"fmt"
	"github.com/bellent69ne/ratnet/ratcrypt"
	//"io/ioutil"
	"log"
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
	MsgDone        = "/* We're done */"
	ServerPort     = ":1366"
	PeerPort       = ":1367"
	NODES          = 1
)

const (
	ErrCantHelp       = "/* Can't help */\n"
	ErrDontUnderstand = "/* I don't understand you */\n"
)

type Session struct {
	Conn        net.Conn
	privateKey  *rsa.PrivateKey
	aesKey      []byte
	alienPubKey *rsa.PublicKey
}

func (self *Session) SetAES(key []byte) {
	self.aesKey = key
}

func (self *Session) SetAlienKey(alien *rsa.PublicKey) {
	self.alienPubKey = alien
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
	self.Conn, err = net.Dial("tcp", endpoint)
	if err != nil {
		return err
	}

	return nil
}

// SendParcel - sends parcel to the connection
func (self *Session) SendParcel(newParcel *Parcel) error {
	var err error
	var encryptedParcel Parcel
	switch string(newParcel.Message) {
	case MsgHelloFriend:
		// don't encrypt parcel
		encryptedParcel = *newParcel
	case MsgHaveAGift:
		// encrypt Parcel with rsa
		encryptedParcel, err = encryptRSA(self, newParcel)
	case ErrDontUnderstand:
		{
			// don't encrypt parcel
			if self.aesKey == nil {
				encryptedParcel = *newParcel
			} else {
				encryptedParcel, err = encryptAES(self.aesKey, newParcel)
			}
		}
	default:
		//encrypt parcel with aes
		encryptedParcel, err = encryptAES(self.aesKey, newParcel)
	}

	if err != nil {
		return err
	}

	parcelEncoder := gob.NewEncoder(self.Conn)

	parcelEncoder.Encode(encryptedParcel)

	return nil
}

//////////////////////////////////////////////////////////////////////////////
////////////////////    Encryption routines  /////////////////////////////////
func encryptRSA(curSession *Session, newParcel *Parcel) (Parcel, error) {
	var encryptedParcel Parcel

	data, err := ratcrypt.EncryptRSA(curSession.alienPubKey,
		newParcel.Message)
	if err != nil {
		return Parcel{nil, nil}, err
	}

	encryptedParcel.Message = data
	data, err = ratcrypt.EncryptRSA(curSession.alienPubKey,
		newParcel.Attachment)
	if err != nil {
		return Parcel{nil, nil}, err
	}

	encryptedParcel.Attachment = data
	return encryptedParcel, nil
}

// Something wrong happens here
func encryptAES(key []byte, newParcel *Parcel) (Parcel, error) {
	newEnvelope, err := ratcrypt.EncryptAES(key,
		newParcel.Message)
	if err != nil {
		return Parcel{nil, nil}, err
	}

	// json encoding
	data, err := json.Marshal(newEnvelope)
	if err != nil {
		return Parcel{nil, nil}, err
	}

	var encryptedParcel Parcel
	encryptedParcel.Message = data

	newEnvelope, err = ratcrypt.EncryptAES(key,
		newParcel.Attachment)
	if err != nil {
		return Parcel{nil, nil}, err
	}
	// json encoding
	data, err = json.Marshal(newEnvelope)
	if err != nil {
		return Parcel{nil, nil}, err
	}

	encryptedParcel.Attachment = data
	return encryptedParcel, nil
}

////////////////////////  Encryption routines  /////////////////////////
////////////////////////////////////////////////////////////////////////

// ReceiveParcel - receives parcel from the connection
func (self *Session) ReceiveParcel() (Parcel, error) {
	var receivedParcel Parcel
	parcelDecoder := gob.NewDecoder(self.Conn)
	parcelDecoder.Decode(&receivedParcel)

	var decryptedParcel Parcel
	var err error

	if string(receivedParcel.Message) == MsgHelloFriend {
		return receivedParcel, nil
	} else if string(receivedParcel.Message) == ErrDontUnderstand {
		return receivedParcel, nil
	} else {
		if self.aesKey == nil {
			decryptedParcel, err = decryptRSA(self.privateKey, &receivedParcel)
		} else {
			decryptedParcel, err = decryptAES(self.aesKey, &receivedParcel)
		}
	}

	return decryptedParcel, err
}

////////////////////////  Decryption routines  //////////////////////////
/////////////////////////////////////////////////////////////////////////

func decryptRSA(private *rsa.PrivateKey, newParcel *Parcel) (Parcel, error) {
	data, err := ratcrypt.DecryptRSA(private, newParcel.Message)
	if err != nil {
		return Parcel{nil, nil}, err
	}

	var decryptedParcel Parcel
	decryptedParcel.Message = data

	data, err = ratcrypt.DecryptRSA(private, newParcel.Attachment)
	if err != nil {
		return Parcel{nil, nil}, err
	}

	decryptedParcel.Attachment = data

	return decryptedParcel, nil
}

func decryptAES(key []byte, newParcel *Parcel) (Parcel, error) {
	var newEnvelope ratcrypt.Envelope
	var data []byte

	if len(newParcel.Message) != 0 {
		err := json.Unmarshal(newParcel.Message, &newEnvelope)
		if err != nil {
			log.Println("Are we in decrypt aes1 ?")
			return Parcel{nil, nil}, err
		}
		data, err = ratcrypt.DecryptAES(key, newEnvelope)
		if err != nil {
			log.Println("Are we in decrypt aes2 ?")
			return Parcel{nil, nil}, err
		}
	}

	var decryptedParcel Parcel
	decryptedParcel.Message = data

	if len(newParcel.Attachment) != 0 {
		err := json.Unmarshal(newParcel.Attachment, &newEnvelope)
		if err != nil {
			log.Println("Are we in decrypt aes3 ?")
			return Parcel{nil, nil}, err
		}

		data, err = ratcrypt.DecryptAES(key, newEnvelope)
		if err != nil {
			log.Println("Are we in decrypt aes4 ?")
			return Parcel{nil, nil}, err
		}
	}

	decryptedParcel.Attachment = data

	return decryptedParcel, nil
}

////////////////////////  Decryption routines  //////////////////////////
/////////////////////////////////////////////////////////////////////////

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
