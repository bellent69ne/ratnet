package ratp

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"log"
)

type Parcel struct {
	Message    []byte
	Attachment []byte
}

// Make - makes appropriate parcel given the message and attachment
//func (self *Parcel) Make(message string, attachment interface{}) error {
func NewParcel(message string, attachment interface{}) (Parcel, error) {
	var err error
	var newParcel Parcel

	switch message {
	case MsgHelloFriend:
		{
			newParcel.Message = []byte(message)
			newParcel.helloFriend(attachment)
		}
	case MsgHaveAGift:
		{
			err = newParcel.haveAgift(attachment)
			if err == nil {
				newParcel.Message = []byte(message)
			}
		}
	case MsgAppreciate:
		newParcel.Message = []byte(message)
	case MsgNeedFriends:
		newParcel.Message = []byte(message)
	case MsgYoureWelcome:
		{
			err = newParcel.youreWelcome(attachment)
			if err == nil {
				newParcel.Message = []byte(message)
			}
		}
	case MsgBeFriends:
		{
			//err = newParcel.beFriends(attachment)
			//if err == nil {
			//	newParcel.Message = []byte(message)
			//}
			newParcel.Message = []byte(message)
			newParcel.beFriends(attachment)
		}
	case MsgWereFriends:
		newParcel.Message = []byte(message)
	case ErrCantHelp:
		newParcel.Message = []byte(message)
	case ErrDontUnderstand:
		newParcel.Message = []byte(message)
	case MsgData:
		{
			err = newParcel.yourData(attachment)
			if err == nil {
				newParcel.Message = []byte(message)
			}
		}
	case MsgDone:
		newParcel.Message = []byte(message)
	default:
		err = errors.New("Wrong message " + message)
	}

	return newParcel, err
}

func (self *Parcel) beFriends(attachment interface{}) {
	//errMsg := MsgBeFriends + " doesn't have associated friends ip addresses"

	switch attachment := attachment.(type) {
	case []byte:
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

func (self *Parcel) yourData(attachment interface{}) error {
	errMsg := MsgData + " doesn't have any attached data"
	return self.byteData(attachment, errMsg)
}

func (self *Parcel) helloFriend(attachment interface{}) {
	switch attachment := attachment.(type) {
	case *rsa.PublicKey:
		{
			data, err := json.Marshal(*attachment)
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
	errMsg := MsgHaveAGift + " doesn't have attached block cipher key"
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
