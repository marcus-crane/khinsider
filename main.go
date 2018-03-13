package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "khinsider"
	app.Usage = "Fetch albums from download.khinsider.com"
	app.Action = func(c *cli.Context) error {
		album := c.Args().Get(0)
		if album != "" {
			fmt.Printf("You got " + album)
			return nil
		}
		return errors.New("Please enter the name of an album eg bubblegum-crisis-tokyo-2040")
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Print(err)
	}
}
