package main

import (
	"fmt"
	"net/http"

	"gopkg.in/redis.v4"
)

func upload(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	r.ParseForm()
	fmt.Println(r.Form)
}
