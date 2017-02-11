package main

import (
	"time"

	"github.com/unixvoid/glogger"
)

func garbageCollectDaemon(gcTime time.Duration) {
	// garbage collect forever
	for {
		glogger.Debug.Printf("sleeping for %v minutes.\n", (gcTime * time.Minute))
		time.Sleep(gcTime * time.Minute)
		glogger.Debug.Println("running garbage collection.")
		go garbageCollect()
	}
}

func garbageCollect() {
	glogger.Debug.Println("IMA COLLECTIN!!!")
}
