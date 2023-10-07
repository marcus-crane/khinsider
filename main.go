package main

import (
	"github.com/marcus-crane/khinsider/v3/cmd/khinsider"
)

var (
	version = "v3.0.0"
	commit  = "n/a"
	date    = "n/a"
	builtBy = "source"
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
