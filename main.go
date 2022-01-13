package main

import (
	"github.com/marcus-crane/khinsider/v2/cmd/khinsider"
)

var (
	version = "2.0.0"
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
