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

type catalystConf struct {
	Appname string `json:"appname"`
	Port    int    `json:"port"`
	Approot string `json:""approot`
}

func init() {
	new = func(command string, params []string) {

		// params
		port := 0
		flag.IntVar(&port, "port", 7777, "your default app port")
		template := "simple"
		flag.StringVar(&template, "tpl", "simple", "start template")
		flag.CommandLine.Parse(os.Args[3:])

		whole := catalystStartAndDone("init your app")
		var done func()

		// 1, archive, unzip and rename
		done = catalystStartAndDone("download project template")
		projectName := "your-awesome-project"
		if len(params) > 0 {
			projectName = params[0]
		}
		if projectName == "" {
			projectName = "your-awesome-project"
		}

		// 1.5, if dir exsit
		if _, err := os.Stat("./" + projectName); err == nil {
			log.Fatal("directory exsit: " + projectName)
		}

		if template == "simple" {
			newProjectCommand := getNewProjectCommand(projectName)
			cmdexer(newProjectCommand)
			done()
		} else if template == "crud" {
			newProjectCommand := getNewProjectWithAutoCrudCommand(projectName)
			cmdexer(newProjectCommand)
			done()
		} else {
			log.Fatal(template + " is invalid template name")
		}

		// 2, rewirte appconf.json with valid approot and appname
		done = catalystStartAndDone("rewrite project config")

		// read origin conf
		jsonDir := "./" + projectName + "/appconf.json"
		catalystconf := catalystConf{}
		catalystconfByte, err := ioutil.ReadFile(jsonDir)
		if err != nil {
			log.Fatal(err)
		}

		// parse workdir
		json.Unmarshal(catalystconfByte, &catalystconf)
		approot := getMyWordDir(projectName)
		catalystconf.Approot = approot
		catalystconf.Appname = projectName
		catalystconf.Port = port

		// write new conf
		catalystconfNewByte, err := json.Marshal(catalystconf)
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(jsonDir, catalystconfNewByte, 644)
		done()
		whole()

	}
	handlers["new"] = new
}

// normal project
func getNewProjectCommand(projectName string) string {
	return "curl -o tmp.zip https://codeload.github.com/shaomingquan/catalyst-sample/zip/master && unzip tmp.zip && rm tmp.zip && mv catalyst-sample-master " + projectName
}

// project with crud
func getNewProjectWithAutoCrudCommand(projectName string) string {
	return "curl -o tmp.zip https://codeload.github.com/shaomingquan/catalyst-crud-sample/zip/master && unzip tmp.zip && rm tmp.zip && mv catalyst-crud-sample-master " + projectName
}
