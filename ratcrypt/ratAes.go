package ratcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// GenerateAESkey - generates AES key of 32 bytes size
func GenerateAESkey() ([]byte, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// EncryptAES - encrypts message with given key with aes in GCM mode.
// Returns encrypted message, nonce for encrypted message, and error
// in case of errors
func EncryptAES(key []byte, message []byte) (Envelope, error) {
	aesGCM, err := makeAesGCM(key)
	if err != nil {
		return Envelope{nil, nil}, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return Envelope{nil, nil}, err
	}

	encryptedMsg := aesGCM.Seal(nil, nonce, message, nil)
	if err != nil {
		return Envelope{nil, nil}, err
	}

	secretEnvelope := Envelope{encryptedMsg, nonce}

	return secretEnvelope, nil
}

// DecryptAES - decrypts message with given key and nonce
// with aes in GCM mode. Returns decrypted message and error
// in case of errors
func DecryptAES(key []byte, secretEnvelope Envelope) ([]byte, error) {
	aesGCM, err := makeAesGCM(key)
	if err != nil {
		return nil, err
	}

	decryptedMsg, err := aesGCM.Open(nil, secretEnvelope.nonce,
		secretEnvelope.message, nil)
	if err != nil {
		return nil, err
	}

	return decryptedMsg, nil
}

func makeAesGCM(key []byte) (cipher.AEAD, error) {
	aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(aes)

	return aesGCM, err
}

type Envelope struct {
	message []byte
	nonce   []byte
}
