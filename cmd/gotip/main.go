package main

import (
	"log"
	"os"

	"github.com/btwiuse/pub"
)

func main() {
	err := pub.RunTip(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
}
