package main

import (
	"log"
	"net/http"

	"github.com/martengine/reception"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", reception.Server())

	log.Fatal(http.ListenAndServe(":8080", mux))
}
