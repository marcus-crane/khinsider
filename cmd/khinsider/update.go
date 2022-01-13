package khinsider

import (
	"github.com/marcus-crane/khinsider/v2/internal/updater"
	"github.com/pterm/pterm"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/marcus-crane/khinsider/v2/pkg/update"
)

func CheckForUpdates(c *cli.Context, currentVersion string, prerelease bool) (bool, string) {
	for _, arg := range c.Args().Slice() {
		if strings.Contains(arg, "update") {
			return false, ""
		}
	}
	if os.Getenv("CI") == "true" || os.Getenv("KHINSIDER_NO_UPDATE") == "true" {
		return false, ""
	}
	remoteVersion := ""
	if !prerelease {
		remoteVersion = update.GetRemoteAppVersion()
	} else {
		remoteVersion = update.GetRemoteAppPrerelease()
	}
	isUpdateAvailable := update.IsRemoteVersionNewer(currentVersion, remoteVersion)
	if isUpdateAvailable {
		return true, remoteVersion
	}
	return false, ""
}

func UpdateAction(c *cli.Context, currentVersion string, prerelease bool) error {
	releaseAvailable, remoteVersion := CheckForUpdates(c, currentVersion, prerelease)
	if !releaseAvailable {
		pterm.Info.Printfln("Sorry, no updates are available. The latest version is %s and you're on %s", remoteVersion, currentVersion)
	}
	return updater.UpgradeInPlace(c.App.Writer, c.App.ErrWriter, remoteVersion)
}
