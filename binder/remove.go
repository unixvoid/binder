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

func remove(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	sec := r.FormValue("sec")
	filename := r.FormValue("filename")

	// make sure all params are set
	if (len(sec) == 0) || (len(filename) == 0) {
		glogger.Debug.Println("not all parameters set")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.Contains(filename, "..") {
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
				println("removing file: ", config.Binder.FileDirectory+filename)
				if err != nil {
					// could not delete, server error
					glogger.Error.Println("could not remove file")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				cleanDir(fmt.Sprintf("%s%s", config.Binder.FileDirectory, filename))
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

func checkDir(dirName string) (bool, error) {
	f, err := os.Open(dirName)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func cleanDir(filename string) {
	// clean up directory resursively if empty
	delimNum := strings.Count(filename, "/")
	var removePath string
	if delimNum != 0 {
		glogger.Debug.Println("directory detected, attempting to clean")
		// i > 1, do not remove config.Binder.FileDirectory aka root
		for i := delimNum; i > 1; i-- {
			tmpDir := strings.Split(filename, "/")
			for y := 0; y < i; y++ {
				removePath = fmt.Sprintf("%s/%s", removePath, tmpDir[y])
			}
			// TODO check config.Binder.FileDirectory to make sure it contains / (in binder.go)
			fullPath := fmt.Sprintf(".%s/", removePath)
			//glogger.Debug.Println("attempting to remove path:", fullPath)
			isEmpty, err := checkDir(fullPath)
			if err != nil {
				glogger.Error.Println(err)
			}
			if isEmpty {
				err := os.Remove(fullPath)
				if err != nil {
					glogger.Error.Println(err)
				}
			} else {
				//glogger.Debug.Println("directory not empty, stopping removal sequence")
				return
			}
			removePath = ""
		}
	}
}
