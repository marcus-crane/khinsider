// Most of this package was taken from
// github.com/superfly/flyctl

package updater

import (
	"fmt"
	"github.com/cli/safeexec"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func isUnderHomebrew() bool {
	brewExe, err := safeexec.LookPath("brew")
	if err != nil {
		return false
	}

	_, err = exec.Command(brewExe, "list", "khinsider").Output()
	if err != nil {
		return false
	}
	return true
}

func updateCommand(version string) string {
	switch {
	case isUnderHomebrew():
		return "brew update && brew install marcus-crane/tap/khinsider"
	case runtime.GOOS == "windows":
		cmd := "iwr https://utf9k.net/khinsider/install.ps1 -useb | iex"
		if version != "" {
			cmd = fmt.Sprintf("$v=\"%s\"; ", version) + cmd
		}
		return cmd
	default:
		return "curl -L \"https://utf9k.net/khinsider/install.sh\" | sh -s " + version
	}
}

func UpgradeInPlace(stdout io.Writer, stderr io.Writer, versionToFetch string) error {
	if runtime.GOOS == "windows" {
		if err := renameCurrentBinaries(); err != nil {
			return err
		}
	}

	shellToUse, ok := os.LookupEnv("SHELL")
	switchToUse := "-c"

	if !ok {
		if runtime.GOOS == "windows" {
			shellToUse = "powershell.exe"
			switchToUse = "-Command"
		} else {
			shellToUse = "/bin/bash"
		}
	}
	fmt.Println(shellToUse, switchToUse)

	command := updateCommand(versionToFetch)

	fmt.Fprintf(stderr, "Running automatic update [%s]\n", command)

	cmd := exec.Command(shellToUse, switchToUse, command)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

// can't replace binary on windows, need to move
func renameCurrentBinaries() error {
	binary, err := os.Executable()
	if err != nil {
		return err
	}

	if err := os.Rename(binary, binary+".old"); err != nil {
		return err
	}
	return nil
}
