package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// your-awesome-app

var new commandHandler

type WebcoreConf struct {
	Appname string `json:"appname"`
	Port    int    `json:"port"`
	Approot string `json:""approot`
}

func init() {
	new = func(command string, params []string) {

		// params
		port := 0
		flag.IntVar(&port, "port", 7777, "your app port")
		template := "simple"
		flag.StringVar(&template, "tpl", "simple", "start template")
		flag.CommandLine.Parse(os.Args[3:])

		whole := webcoreStartAndDone("init your app")
		var done func()

		// 1, archive, unzip and rename
		done = webcoreStartAndDone("download project template")
		projectName := "your-awesome-project"
		if len(params) > 0 {
			projectName = params[0]
		}
		if projectName == "" {
			projectName = "your-awesome-project"
		}

		if template == "simple" {
			newProjectCommand := getNewProjectCommand(projectName)
			cmdexer(newProjectCommand)
			done()
		} else if template == "curd" {
			newProjectCommand := getNewProjectWithAutoCurdCommand(projectName)
			cmdexer(newProjectCommand)
			done()
		} else {
			log.Fatal(template + " is invalid template name")
		}

		// 2, rewirte appconf.json with valid approot and appname
		done = webcoreStartAndDone("rewrite project config")

		// read origin conf
		jsonDir := "./" + projectName + "/appconf.json"
		webcoreconf := WebcoreConf{}
		webcoreconfByte, err := ioutil.ReadFile(jsonDir)
		if err != nil {
			log.Fatal(err)
		}

		// parse workdir
		json.Unmarshal(webcoreconfByte, &webcoreconf)
		approot := getMyWordDir(projectName)
		webcoreconf.Approot = approot
		webcoreconf.Appname = projectName
		webcoreconf.Port = port

		// write new conf
		webcoreconfNewByte, err := json.Marshal(webcoreconf)
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(jsonDir, webcoreconfNewByte, 644)
		done()
		whole()

	}
	handlers["new"] = new
}

// normal project
func getNewProjectCommand(projectName string) string {
	return "curl -o tmp.zip https://codeload.github.com/shaomingquan/webcore-sample/zip/master && unzip tmp.zip && rm tmp.zip && mv webcore-sample-master " + projectName
}

// project with curd
func getNewProjectWithAutoCurdCommand(projectName string) string {
	return "curl -o tmp.zip https://codeload.github.com/shaomingquan/webcore-curd-sample/zip/master && unzip tmp.zip && rm tmp.zip && mv webcore-curd-sample-master " + projectName
}
