package main

import (
	"flag"
	"os"
	"strconv"
)

var dev commandHandler

func init() {
	dev = func(command string, params []string) {

		withoutrun := false
		port := 0
		flag.BoolVar(&withoutrun, "withoutrun", false, "on genefile without run")
		flag.IntVar(&port, "port", 0, "app start port")
		flag.CommandLine.Parse(os.Args[2:])

		// 1, clear genfiles
		done := webcoreStartAndDone("clear genfiles")
		cmdexer(`
			rm boot.go
			touch boot.go
			rm ./imports/*
		`)
		done()

		// 2, genfile
		done = webcoreStartAndDone("generate boot.go and importfiles")
		cmdexer(`
			go generate ./apps/...
		`)
		done()

		if !withoutrun {
			// 3, start dev
			done = webcoreStartAndDone("start dev app")
			if port == 0 {
				cmdwithstdout(`
					DEBUG=true go run *.go
				`)
			} else {
				cmdwithstdout(`
					DEBUG=true go run *.go --port=` + strconv.Itoa(port) + `
				`)
			}
			done()
		}

	}
	handlers["dev"] = dev
}
