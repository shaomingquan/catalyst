package main

var run commandHandler

func init() {
	run = func(command string, params []string) {

		withGene := false
		if len(params) > 0 && params[0] == "withgene" {
			withGene = true
		}

		// 1, gene file
		if withGene {
			gene(command, params) // what diff when from run
		}

		// 2, run project
		cmdexer("go run *.go")

	}
	handlers["run"] = run
}
