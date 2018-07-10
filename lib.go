package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func webcoreSay(content string) {
	fmt.Println("[ webcore ] : " + content)
}

func webcoreStartAndDone(content string) func() {
	webcoreSay(content + " ...start")
	return func() {
		webcoreSay(content + " ...done")
		fmt.Println("")
	}
}

func cmdexer(cmdstr string) string {
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	return out.String()
}

func cmdarrexer(c string, args []string) string {
	cmd := exec.Command(c, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	return out.String()
}

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func myGopath() []string {
	gopathString := os.Getenv("GOPATH")
	gopathes := strings.Split(gopathString, ":")
	return gopathes
}

func getMyWordDir(projectName string) string {
	pwd := pwd()
	gopathes := myGopath()

	workdirPrefix := ""

	for _, gopath := range gopathes {

		srcdir := ""
		if endWithSlash(gopath) {
			srcdir = gopath + "src/"
		} else {
			srcdir = gopath + "/src/"
		}

		if strings.HasPrefix(pwd, srcdir) {
			workdirPrefix = pwd[len(srcdir):]
			break
		}
	}

	if workdirPrefix == "" {
		log.Fatal(errors.New("play under your gopath src dir please"))
	}

	if !endWithSlash(workdirPrefix) {
		workdirPrefix = workdirPrefix + "/"
	}

	return workdirPrefix + projectName
}

func endWithSlash(str string) bool {
	if str == "" {
		return false
	}
	if str[len(str)-1] == '/' {
		return true
	}
	return false
}
