package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

// your-awesome-app

var new commandHandler

type WebcoreConf struct {
	Appname string `json:"appname"`
	Port    int    `json:"port"`
	Approot string `json:""approot`
}

func init() {
	run = func(command string, params []string) {

		// 1, archive, unzip and rename
		projectName := "your-awesome-project"
		if len(params) > 0 {
			projectName = params[0]
		}
		if projectName == "" {
			projectName = "your-awesome-project"
		}

		newProjectCommand := getNewProjectCommand(projectName)
		cmdexer(newProjectCommand)

		// 2, rewirte appconf.json with valid approot and appname

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

		// parse port, why the f**king flag dont work????
		port := 0
		flag.IntVar(&port, "port", 7777, "your app port")
		flag.Parse()
		webcoreconf.Port = port

		// write new conf
		webcoreconfNewByte, err := json.Marshal(webcoreconf)
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(jsonDir, webcoreconfNewByte, 644)

	}
	handlers["new"] = run
}

// normal project
func getNewProjectCommand(projectName string) string {
	return "curl -o tmp.zip https://codeload.github.com/shaomingquan/webcore-sample/zip/master && unzip tmp.zip && rm tmp.zip && mv webcore-sample-master " + projectName
}

// project with curd
func getNewProjectWithAutoCurdCommand(projectName string) string {
	return "curl -o tmp.zip https://codeload.github.com/shaomingquan/webcore-curd-sample/zip/master && unzip tmp.zip && rm tmp.zip && mv webcore-sample-master " + projectName
}
