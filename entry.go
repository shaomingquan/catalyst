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

	handlers[command](command, params)

}
