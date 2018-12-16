package ratcrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

const (
	ErrDontHavePrivateKey = "RSA... Don't have private key"
	ErrDontHavePublicKey  = "RSA... Don't have public key"
)

// GenerateRSAkey - generates an RSA keypair of the 4096 bit size
func GenerateRSAkey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		//handle error
		return nil, err
	}

	return privateKey, nil
}

// Encrypt - encrypts given message with the public key
func EncryptRSA(pub *rsa.PublicKey, message []byte) ([]byte, error) {
	if pub == nil {
		return nil, errors.New(ErrDontHavePublicKey)
	}
	encryptedMsg, err := rsa.EncryptOAEP(sha256.New(),
		rand.Reader, pub, message, nil)
	if err != nil {
		return nil, err
	}

	return encryptedMsg, nil
}

// Decrypt - decrypts given message with the private key
func DecryptRSA(private *rsa.PrivateKey, message []byte) ([]byte, error) {
	if private == nil {
		return nil, errors.New(ErrDontHavePrivateKey)
	}
	decryptedMsg, err := rsa.DecryptOAEP(sha256.New(), rand.Reader,
		private, message, nil)

	if err != nil {
		return nil, err
	}

	return decryptedMsg, nil
}
