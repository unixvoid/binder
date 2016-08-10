package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v4"
)

func register(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	// check if instance is already registered
	_, err := client.Get("sec").Result()
	if err != redis.Nil {
		glogger.Debug.Println("instance already registered.")
		glogger.Debug.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// instance is not registered, generate key
		sec := randStr(config.Binder.SecTokenSize)
		secHash := sha3.Sum512([]byte(sec))

		// upload sec key to server
		err := client.Set("sec", fmt.Sprintf("%x", secHash), 0).Err()
		if err != nil {
			// cannot update sec key
			glogger.Error.Println("error in setting sec key in redis")
			glogger.Error.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			//return security token to client
			w.Header().Set("secToken", sec)
			fmt.Fprintf(w, "%s", sec)
		}
	}
}

func randStr(strSize int) string {
	dictionary := config.Binder.SecDictionary

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
