package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unixvoid/glogger"
	"gopkg.in/gcfg.v1"
	"gopkg.in/redis.v4"
)

type Config struct {
	Binder struct {
		Port          int
		Loglevel      string
		SecDictionary string
		SecTokenSize  int
		FileDirectory string
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
	// make sure FileDirectory exists
	_, err = os.Stat(config.Binder.FileDirectory)
	if err != nil {
		glogger.Debug.Println(config.Binder.FileDirectory + " does not exist, creating")
		os.Mkdir(config.Binder.FileDirectory, 0777)
	}

	// all handlers
	// TODO make async
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, redisClient)
	}).Methods("GET")
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		upload(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		remove(w, r, redisClient)
	}).Methods("POST")

	// TODO SSL here
	// https://github.com/unixvoid/beacon/blob/develop/beacon/beacon.go#L76-L94
	glogger.Info.Println("binder running http on", config.Binder.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Binder.Port), router))

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
	if config.Binder.Loglevel == "debug" {
		glogger.LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else if config.Binder.Loglevel == "cluster" {
		glogger.LogInit(os.Stdout, os.Stdout, ioutil.Discard, os.Stderr)
	} else if config.Binder.Loglevel == "info" {
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
