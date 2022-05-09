package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

type Encryptor struct {
	key []byte
}

func NewEncryptor(key []byte) *Encryptor {
	return &Encryptor{
		key:   key,
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
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}


	nonce, err := GenerateRandom(gcm.NonceSize())
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	out := gcm.Seal(nil, nonce, userID, nil)
	out = append(out, nonce...)

	return string(out), nil
}

func (e *Encryptor) Decode(userID string) ([]byte, error) {
	aesblock, err := aes.NewCipher(e.key)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	b := []byte(userID)
	nonce := b[len(b) - aesgcm.NonceSize():] 
	userIDByte := b[:len(b) - aesgcm.NonceSize()]  

	decrypted, err := aesgcm.Open(nil, nonce, userIDByte, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	return decrypted, nil
}
