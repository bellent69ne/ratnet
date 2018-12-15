package ratp

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"log"
)

type Parcel struct {
	Message    string
	Attachment []byte
}

// Make - makes appropriate parcel given the message and attachment
func (self *Parcel) Make(message string, attachment interface{}) error {
	var err error

	switch message {
	case MsgHelloFriend:
		{
			self.Message = message
			self.helloFriend(attachment)
		}
	case MsgHaveAGift:
		{
			err = self.haveAgift(attachment)
			if err == nil {
				self.Message = message
			}
		}
	case MsgAppreciate:
		self.Message = message
	case MsgNeedFriends:
		self.Message = message
	case MsgYoureWelcome:
		{
			err = self.youreWelcome(attachment)
			if err == nil {
				self.Message = message
			}
		}
	case MsgBeFriends:
		{
			err = self.beFriends(attachment)
			if err == nil {
				self.Message = message
			}
		}
	case MsgWereFriends:
		self.Message = message
	case ErrCantHelp:
		self.Message = message
	case ErrDontUnderstand:
		self.Message = message
	case MsgData:
		{
			err = self.yourData(attachment)
			if err == nil {
				self.Message = message
			}
		}
	default:
		err = errors.New("Wrong message " + message)
	}

	return err
}

func (self *Parcel) beFriends(attachment interface{}) error {
	errMsg := MsgBeFriends + " has to have associated ip addresses in []byte"
	return self.byteData(attachment, errMsg)
}

func (self *Parcel) yourData(attachment interface{}) error {
	errMsg := MsgData + " should have []byte of data attached"
	return self.byteData(attachment, errMsg)
}

func (self *Parcel) helloFriend(attachment interface{}) {
	switch attachment := attachment.(type) {
	case rsa.PublicKey:
		{
			data, err := json.Marshal(attachment)
			if err != nil {
				log.Println(err)
			}
			self.Attachment = data
		}
	default:
		self.Attachment = nil
	}
}

func (self *Parcel) haveAgift(attachment interface{}) error {
	errMsg := MsgHaveAGift + " has to have block cipher key attached"
	return self.byteData(attachment, errMsg)
}

func (self *Parcel) byteData(attachment interface{}, errMsg string) error {
	switch attachment := attachment.(type) {
	case []byte:
		self.Attachment = attachment
	default:
		return errors.New(errMsg)
	}

	return nil
}

func (self *Parcel) youreWelcome(attachment interface{}) error {
	errMsg := MsgYoureWelcome + " has to have []byte of ip addresses attached"
	return self.byteData(attachment, errMsg)
}
