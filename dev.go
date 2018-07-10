package main

import (
	"flag"
	"os"
)

var dev commandHandler

func init() {
	dev = func(command string, params []string) {

		withoutrun := false
		flag.BoolVar(&withoutrun, "withoutrun", false, "on genefile without run")
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
			cmdwithstdout(`
				DEBUG=true go run *.go
			`)
			done()
		}

	}
	handlers["dev"] = dev
}
