package main

import (
	"bytes"
	"log"
	"os/exec"
)

func cmdexer(cmdstr string) string {
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	log.Fatal(err.Error())
	return out.String()
}

func cmdarrexer(c string, args []string) string {
	cmd := exec.Command(c, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	log.Fatal(err.Error())
	return out.String()
}
