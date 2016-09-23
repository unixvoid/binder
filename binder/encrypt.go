package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/unixvoid/glogger"
	"gopkg.in/redis.v4"
)

func encryptString(secHash [64]byte, value string, client *redis.Client) []byte {
	key := fmt.Sprintf("%x", secHash)

	block, err := aes.NewCipher([]byte(key[0:32]))
	if err != nil {
		glogger.Error.Println(err)
	}
	str := []byte(value)

	ciphertext := make([]byte, aes.BlockSize+len([]byte(value)))
	iv := ciphertext[:aes.BlockSize]

	// encrypt
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(str))
	encrypter.XORKeyStream(encrypted, str)

	return encrypted
}
