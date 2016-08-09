package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unixvoid/glogger"
	"gopkg.in/gcfg.v1"
	"gopkg.in/redis.v4"
)

type Config struct {
	Server struct {
		Port     int
		Loglevel string
	}
	Redis struct {
		Host     string
		Password string
	}
}

var (
	config = Config{}
)

func main() {
	// initalize conf and logs
	readConf()
	initLogger()
	redisClient, err := initRedisConnection()
	if err != nil {
		glogger.Error.Println("redis connection cannot be made.")
	} else {
		glogger.Info.Println("connection to redis succeeded.")
	}

	// all handlers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		upload(w, r, redisClient)
	})
}

func readConf() {
	// init config file
	err := gcfg.ReadFileInto(&config, "config.gcfg")
	if err != nil {
		panic(fmt.Sprintf("Could not load config.gcfg, error: %s\n", err))
	}
	return
}

func initLogger() {
	// init logger
	if config.Server.Loglevel == "debug" {
		glogger.LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else if config.Server.Loglevel == "cluster" {
		glogger.LogInit(os.Stdout, os.Stdout, ioutil.Discard, os.Stderr)
	} else if config.Server.Loglevel == "info" {
		glogger.LogInit(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
	} else {
		glogger.LogInit(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	}
}

func initRedisConnection() (*redis.Client, error) {
	// init redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       0,
	})

	_, redisErr := redisClient.Ping().Result()
	return redisClient, redisErr
}
