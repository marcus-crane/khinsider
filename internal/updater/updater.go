// Most of this package was taken from
// github.com/superfly/flyctl

package updater

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cli/safeexec"
)

func isUnderHomebrew() bool {
	flyBinary, err := os.Executable()
	if err != nil {
		return false
	}

	brewExe, err := safeexec.LookPath("brew")
	if err != nil {
		return false
	}

	brewPrefixBytes, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}

	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBytes)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(flyBinary, brewBinPrefix)
}

func updateCommand(version string) string {
	switch {
	case isUnderHomebrew():
		return "brew upgrade khinsider"
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
