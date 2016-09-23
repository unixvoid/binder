package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/unixvoid/glogger"
)

func decryptString(secHash [64]byte, encrypted string) string {
	key := fmt.Sprintf("%x", secHash)

	block, err := aes.NewCipher([]byte(key[0:32]))
	if err != nil {
		glogger.Error.Println(err)
	}

	str := []byte(encrypted)
	ciphertext := make([]byte, aes.BlockSize+len([]byte(encrypted)))
	iv := ciphertext[:aes.BlockSize]

	// decrypt test
	decrypter := cipher.NewCFBDecrypter(block, iv)

	decrypted := make([]byte, len(str))
	decrypter.XORKeyStream(decrypted, []byte(encrypted))

	return string(decrypted)
}
