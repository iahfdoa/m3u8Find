package main

import (
	"github.com/iahfdoa/m3u8Find/pkg/runner"
	"github.com/projectdiscovery/gologger"
)

func main() {

	userOptions := runner.ParseUserOptions()

	run, err := runner.NewRunner(userOptions)

	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	run.Run()
}
