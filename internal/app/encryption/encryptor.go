package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
)

type Encryptor struct {
	key []byte
}

func NewEncryptor(key []byte) *Encryptor {
	return &Encryptor{
		key: key,
	}
}

type EncryptorMethodes interface {
	Decode(userID []byte) (string, error)
	Encode() (string, error)
}

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (e *Encryptor) Encode(userID []byte) (string, error) {
	c, err := aes.NewCipher(e.key)
	if err != nil {
		return "Cipher", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "GCM", err
	}

	nonce, err := GenerateRandom(gcm.NonceSize())
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	out := gcm.Seal(nonce, nonce, userID, nil)
	
	return b64.StdEncoding.EncodeToString([]byte(out)), nil
}

func (e *Encryptor) Decode(userID string) (string, error) {
	b64Decode, _ := b64.StdEncoding.DecodeString(userID)

	aesblock, err := aes.NewCipher(e.key)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	nonce, b64UserID := b64Decode[:nonceSize], b64Decode[nonceSize:]

	decrypted, err := aesgcm.Open(nil, nonce, b64UserID, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return "", err
	}

	return string(decrypted), nil
}
