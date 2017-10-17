package main

import (
	"log"
	"time"

	"github.com/itsbth/glox/scanner"
)

func main() {
	log.Printf("Hello, World!")
	scan := scanner.NewScanner("\"string")
	go scan.Scan()
	go func() {
		// this is not a workable solution
		for err := range scan.Errors {
			log.Println(err.Error())
		}
	}()
	for token := range scan.Tokens {
		log.Printf("token: %+v\n", token)
	}
	time.Sleep(2 * time.Second)
}
