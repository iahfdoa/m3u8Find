package main

import (
	"github.com/projectdiscovery/gologger"
	"m3u8Find/pkg/runner"
)

func main() {

	userOptions := runner.ParseUserOptions()

	run, err := runner.NewRunner(userOptions)

	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	run.Run()
}
