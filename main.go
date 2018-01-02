package main

import (
	"fmt"
	"log"
)

func main() {
	msg, err := RunAppCommands()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(msg)
}
