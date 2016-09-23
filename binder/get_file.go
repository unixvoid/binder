package main

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v4"
)

func getfile(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	//r.ParseForm()
	sec := r.FormValue("sec")
	key := r.FormValue("key")

	// make sure all params are set
	if (len(sec) == 0) || (len(key) == 0) {
		glogger.Debug.Println("not all parameters set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// find the hash of the provided security token
	secHash := sha3.Sum512([]byte(sec))

	// make sure sec has been set
	storedSecHash, err := client.Get("sec").Result()
	if err != redis.Nil {
		// sec exists, check auth
		// check if auth is valid
		if fmt.Sprintf("%x", secHash) == storedSecHash {
			// client is authed, upload
			keyHash := sha3.Sum512([]byte(key))

			encryptedValue, err := client.Get(fmt.Sprintf("hkey:file:%x", keyHash)).Result()
			if err != nil {
				glogger.Error.Println("error getting hkey in redis")
			}

			decryptedValue := decryptString(secHash, encryptedValue)
			decodeVal, _ := base64.StdEncoding.DecodeString(decryptedValue)
			fmt.Fprintf(w, "%s", decodeVal)
		} else {
			// client auth failed
			glogger.Debug.Println("client auth failed")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		// sec not set
		glogger.Debug.Println("sec not set while trying to set key")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
