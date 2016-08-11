package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v4"
)

func upload(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	sec := r.FormValue("sec")
	filename := r.FormValue("filename")
	file, handler, err := r.FormFile("file")

	// make sure all params are set
	if (len(sec) == 0) || (file == nil) {
		glogger.Debug.Println("not all parameters set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// set filename equal to the uploaded name if it is not set.
	if len(filename) == 0 {
		glogger.Debug.Println("filename not set, setting to " + handler.Filename)
		filename = handler.Filename
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
			// client is authed, upload
			f, err := os.Create(config.Binder.FileDirectory + filename)
			if err != nil {
				glogger.Error.Println("could not write file to filesystem")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = io.Copy(f, file)
			if err != nil {
				glogger.Error.Println("could not write file to filesystem")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			// client auth failed
			glogger.Debug.Println("client auth failed")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		// sec not set
		glogger.Debug.Println("sec not set while trying to upload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
