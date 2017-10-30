package main

import (
	"fmt"
	"log"
)

func main() {
	msg, err := ParseCLICommands()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(msg)
}
