package main

import (
	"log"
	"os"
)

type commandHandler func(command string, params []string)

var handlers = map[string]commandHandler{}

func main() {

	args := os.Args

	if len(args) < 2 {
		log.Fatal("must with command")
	}

	command := os.Args[1]

	params := []string{}
	if len(args) > 2 {
		params = os.Args[2:]
	}

	if _, ok := handlers[command]; ok {
		handlers[command](command, params)
	} else {
		log.Fatal(command + " is invalid command")
	}

}
