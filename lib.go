package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func webcoreSay(content string) {
	fmt.Println("\033[0;32m[ webcore ]\033[0m : " + content)
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

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
	// never reached
	panic(true)
	return nil, nil
}

func cmdwithstdout(cmdstr string) {
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()

	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()

	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()

	err := cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
