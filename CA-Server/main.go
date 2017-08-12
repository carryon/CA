package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", CaHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
