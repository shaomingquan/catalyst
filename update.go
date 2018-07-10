package main

var update commandHandler

func init() {
	update = func(command string, params []string) {

		// 1, just exe sciprt
		done := webcoreStartAndDone("update runtime libs and update webcore exe")
		cmdexer(`
			govendor fetch github.com/shaomingquan/webcore/gene
			govendor fetch github.com/shaomingquan/webcore/core
			govendor fetch github.com/shaomingquan/webcore

			go get github.com/shaomingquan/webcore
		`)
		done()

	}
	handlers["update"] = update
}
