package main

import (
	"fmt"
	"log"
	"os"

	"github.com/btwiuse/pub"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(pub.Usage)
		return
	}
	err := pub.Run(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
}
