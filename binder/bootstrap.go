package main

import (
	"fmt"

	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v4"
)

func bootstrapCheck(client *redis.Client) {
	// check if instance is already registered
	_, err := client.Get("sec").Result()
	if err != redis.Nil {
		glogger.Debug.Println("instance already registered while bootstrapping")
		return
	} else {
		// instance is not registered, generate key
		sec := randStr(config.Binder.SecTokenSize)
		secHash := sha3.Sum512([]byte(sec))

		// upload sec key to server
		err := client.Set("sec", fmt.Sprintf("%x", secHash), 0).Err()
		if err != nil {
			// cannot update sec key
			glogger.Error.Println("error in setting sec key in redis while bootstrapping")
			glogger.Error.Println(err)
			return
		} else {
			//return security token
			glogger.Info.Println("environment bootstrapped with: " + sec)
		}
	}
}
