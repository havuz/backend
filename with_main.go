// +build with_main

package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	ADDR := ":8090"
	if addrFromEnv := os.Getenv("ADDR"); addrFromEnv != "" {
		ADDR = addrFromEnv
	}

	http.Handle("/", http.HandlerFunc(Handler))

	log.Printf("Reverse Proxy is listening at: %+q\n\n", ADDR)
	log.Fatal(http.ListenAndServe(ADDR, nil))
}
