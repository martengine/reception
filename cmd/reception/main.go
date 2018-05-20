package main

import (
	"log"
	"net/http"

	"github.com/martengine/reception"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", reception.Server())
	mux.Handle("/services", reception.Info())

	log.Fatal(http.ListenAndServeTLS(":8080", "cert/server.crt", "cert/server.key", mux))
}
