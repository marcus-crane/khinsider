package main

import (
	"github.com/marcus-crane/khinsider/v2/cmd/khinsider"
)

var (
	version = "dev"
	commit  = "n/a"
	date    = "n/a"
	builtBy = "dev"
)

func main() {
	buildInfo := khinsider.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	}
	khinsider.Execute(buildInfo)
}
