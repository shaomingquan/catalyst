package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/shaomingquan/webcore/core"
)

var _conf core.Conf

func getConf() *core.Conf {
	approot := rootRelatedString()
	if _conf != (core.Conf{}) {
		return &_conf
	} else {
		// read from exe dir.
		filebyted, err := ioutil.ReadFile(approot + "appconf.json")
		if err != nil {
			log.Fatal(err.Error())
		}

		json.Unmarshal(filebyted, &_conf)

		if _conf.AppName == "" {
			log.Fatal("appconf appname is required")
		}

		if _conf.Port == 0 {
			log.Fatal("appconf port is required")
		}

		if _conf.AppRoot == "" {
			log.Fatal("appconf approot is required")
		}

		return &_conf
	}
}
