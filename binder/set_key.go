package main

import (
	"fmt"
	"net/http"

	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v4"
)

func setkey(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	sec := r.FormValue("sec")
	key := r.FormValue("key")
	value := r.FormValue("value")

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

			encryptedValue := encryptString(secHash, value, client)
			//err := client.Set(fmt.Sprintf("hkey:%x", keyHash), fmt.Sprintf("%x", encryptedValue), 0).Err()
			err := client.Set(fmt.Sprintf("hkey:%x", keyHash), encryptedValue, 0).Err()
			if err != nil {
				glogger.Error.Println("error setting hkey in redis")
			}
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
