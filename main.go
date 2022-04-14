package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Ready to listen!")
	run()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
