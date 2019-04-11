package main

var update commandHandler

func init() {
	update = func(command string, params []string) {

		// 1, just exe sciprt
		done := catalystStartAndDone("update runtime libs and update catalyst exe")
		cmdexer(`
			govendor fetch github.com/shaomingquan/catalyst/gene
			govendor fetch github.com/shaomingquan/catalyst/core
			govendor fetch github.com/shaomingquan/catalyst

			go get github.com/shaomingquan/catalyst
		`)
		done()

	}
	handlers["update"] = update
}
