package ratp

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
)

// Handshake - makes initial handshaking between
// two peers in the session
func Handshake(curSession *Session) bool {
	// if couldn't say "hello fri3nd"
	// nothing to do
	if !SayHelloFriend(curSession) {
		return false
	}

	// if couldn't receive "hello fri3nd"
	if !ReceiveHelloFriend(curSession) {
		// then he can fuck off
		_ = SayFuckOFF(curSession)
		return false
	}

	return true
}

func printParcel(parcel *Parcel) {
	fmt.Println(string(parcel.Message))
	fmt.Println(parcel.Attachment)
}

// SayHelloFriend - sends "hello fri3nd" parcel to the current session
func SayHelloFriend(curSession *Session) bool {
	// Message hello fri3nd
	curParcel, err := NewParcel(MsgHelloFriend, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	//////////////////////////////////////////////////////
	// Send Message hello fri3nd
	err = curSession.SendParcel(&curParcel)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// ReceiveHelloFriend - receive parcel with "hello fri3nd" message
func ReceiveHelloFriend(curSession *Session) bool {
	// Receive parcel from the session
	// Should have "hello fri3nd" message
	// with remote peers public key attached
	gotParcel, err := curSession.ReceiveParcel()
	// if couldn't receive parcel
	if err != nil {
		// log why, exit
		log.Println(err)
		return false
	}
	printParcel(&gotParcel)
	// if message in received parcel is not "hello fri3nd"
	// exit
	if string(gotParcel.Message) != MsgHelloFriend {
		return false
	}

	// Decode alien rsa public key for this session
	var alien rsa.PublicKey
	if len(gotParcel.Attachment) != 0 {
		err = json.Unmarshal(gotParcel.Attachment, &alien)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	// set remote peers public key
	curSession.SetAlienKey(&alien)
	return true
}

// SayFuckOFF - sends "don't understand" parcel to the session
func SayFuckOFF(curSession *Session) bool {
	// make fuckoff parcel
	justFuckOff, err := NewParcel(ErrDontUnderstand, nil)
	// if couldn't create
	if err != nil {
		// log why, exit
		log.Println(err)
		return false
	}

	// send fuckoff parcel to the session
	err = curSession.SendParcel(&justFuckOff)
	// if couldn't send fuckoff parcel
	if err != nil {
		// log why, exit
		log.Println(err)
		return false
	}

	return true
}

// SecureSession - makes session secure negotiating aes key
// for communication
func SecureSession(curSession *Session) bool {
	// Now generate aes encryption key for this session
	err := curSession.GenerateAESkey()
	if err != nil {
		log.Println(err)
		return false
	}
	///////////////////////////////////////////////////

	// Create Message "I have a gift" with encrypted aes key
	curParcel, err := NewParcel(MsgHaveAGift, curSession.AesKey())
	if err != nil {
		log.Println(err)
		return false
	}
	/////////////////////////////////////////////////////////
	// Now encrypt the Message "I have a gift" with aes key attached
	// and then send it over the network
	err = curSession.SendParcel(&curParcel)
	if err != nil {
		log.Println(err)
		return false
	}

	if !Appreciated(curSession) {
		_ = SayFuckOFF(curSession)
		return false
	}

	return true
}

// Appreciated - receives parcel and checks whether it is
// appreciation. Returns true if it is appreciation
func Appreciated(curSession *Session) bool {
	// Receive appreciation from the ratnet server
	gotParcel, err := curSession.ReceiveParcel()
	if err != nil {
		log.Println(err)
		return false
	}
	//printParcel(&gotParcel)
	// if message in received parcel is not appreciation
	// nothing to do
	if string(gotParcel.Message) != MsgAppreciate {
		return false
	}

	return true
}

func GetFriendsAddrs(curSession *Session) []string {
	if !SayINeedFriends(curSession) {
		return nil
	}

	addrs, err := ReceiveFriends(curSession)
	if err != nil {
		log.Println(err)
		return nil
	}

	return addrs
}

func SayINeedFriends(curSession *Session) bool {
	// Now create Message "I need fri3nds"
	curParcel, err := NewParcel(MsgNeedFriends, nil)
	// if couldn't create parcel with message "I need fri3nds"
	if err != nil {
		// log why, exit
		log.Println(err)
		return false
	}
	//////////////////////////////////////////////////

	// Encrypt that message and send it over the network
	err = curSession.SendParcel(&curParcel)
	// if couldn't send parcel to the session
	if err != nil {
		// log why, exit
		log.Println(err)
		return false
	}

	return true
}

func ReceiveFriends(curSession *Session) ([]string, error) {
	// Receive parcel from the session
	// should have "You're welcome" message in it
	gotParcel, err := curSession.ReceiveParcel()
	if err != nil {
		return nil, err
	}
	if string(gotParcel.Message) != MsgYoureWelcome {
		return nil, err
	}

	friendsAddrs, err := decodeFriends(gotParcel.Attachment)
	if err != nil {
		return nil, err
	}

	return friendsAddrs, nil
}

func decodeFriends(friendsIps []byte) ([]string, error) {
	var friendsAddrs []string
	if len(friendsIps) != 0 {
		err := json.Unmarshal(friendsIps, &friendsAddrs)
		if err != nil {
			return nil, err
		}
	}

	return friendsAddrs, nil
}

func BeFriends(curSession *Session, data []byte) bool {
	newParcel, err := NewParcel(MsgBeFriends, data)
	if err != nil {
		log.Println(err)
		return false
	}

	err = curSession.SendParcel(&newParcel)
	if err != nil {
		log.Println(err)
		return false
	}

	// think about it
	if !acceptedFriendship(curSession) {
		return false
	}

	return true
}

func acceptedFriendship(curSession *Session) bool {
	gotParcel, err := curSession.ReceiveParcel()
	if err != nil {
		log.Println(err)
		return false
	}

	if string(gotParcel.Message) != MsgWereFriends {
		return false
	}

	return true

}

func WantsToBeFriends(curSession *Session) ([]byte, error) {
	gotParcel, err := curSession.ReceiveParcel()
	if err != nil {
		return nil, err
	}

	if string(gotParcel.Message) == MsgBeFriends {
		err = SayWereFriends(curSession)
		if err != nil {
			return nil, err
		}
	}

	return gotParcel.Attachment, nil
}

func SayWereFriends(curSession *Session) error {
	curParcel, err := NewParcel(MsgWereFriends, nil)
	if err != nil {
		return err
	}

	err = curSession.SendParcel(&curParcel)
	if err != nil {
		return err
	}

	return nil
}

func SayBeFriends(curSession *Session, ipAddrs []byte) error {
	parcel, err := NewParcel(MsgBeFriends, ipAddrs)
	if err != nil {
		return err
	}

	err = curSession.SendParcel(&parcel)
	if err != nil {
		return err
	}

	return nil
}

func ReceiveData(curSession *Session) ([]byte, error) {
	parcel, err := curSession.ReceiveParcel()
	if err != nil {
		return nil, err
	}

	if string(parcel.Message) != MsgData {
		return nil, err
	}

	return parcel.Attachment, nil
}

func SayDone(curSession *Session) error {
	parcel, err := NewParcel(MsgDone, nil)
	if err != nil {
		return err
	}

	err = curSession.SendParcel(&parcel)
	if err != nil {
		return err
	}

	return nil
}
