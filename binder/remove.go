package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v4"
)

func remove(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	sec := r.FormValue("sec")
	filename := r.FormValue("filename")

	// make sure all params are set
	if (len(sec) == 0) || (len(filename) == 0) {
		glogger.Debug.Println("not all parameters set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.Contains(filename, "/") {
		glogger.Debug.Println("filename contains malicious characters, stopping")
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
			// check if file exists
			_, err := os.Stat(config.Binder.FileDirectory + filename)
			if err != nil {
				// file does not exist, bad request
				glogger.Debug.Println("file does not exist, cannot remove")
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				// client is authed, remove
				err := os.Remove(config.Binder.FileDirectory + filename)
				if err != nil {
					// could not delete, server error
					glogger.Error.Println("could not remove file")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		} else {
			// client auth failed
			glogger.Debug.Println("client auth failed")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		// sec not set
		glogger.Debug.Println("sec not set while trying to remove")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
