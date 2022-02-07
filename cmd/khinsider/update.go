package khinsider

import (
	"github.com/marcus-crane/khinsider/v2/internal/updater"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"strings"

	"github.com/marcus-crane/khinsider/v2/pkg/update"
)

func isUpdaterDisabled(c *cli.Context) bool {
	return c.Bool("no-updates")
}

func CheckForUpdates(c *cli.Context, currentVersion string, prerelease bool) (bool, string) {
	if isUpdaterDisabled(c) && c.Command.Name != "update" || strings.Contains(currentVersion, "-DEV") {
		pterm.Debug.Println("Updater is disabled. Skipping update check.")
		return false, ""
	}
	remoteVersion := ""
	if !prerelease {
		remoteVersion = update.GetRemoteAppVersion()
		pterm.Debug.Printfln("Found remote version: %s", remoteVersion)
	} else {
		pterm.Debug.Printfln("Found remote prerelease version: %s", remoteVersion)
		remoteVersion = update.GetRemoteAppPrerelease()
	}
	if remoteVersion == "" {
		// Assume we were rate limited so skip update check for now
		return false, remoteVersion
	}
	isUpdateAvailable := update.IsRemoteVersionNewer(currentVersion, remoteVersion)
	pterm.Debug.Printfln("Current is %s while remote is %s. Update is available: %t", currentVersion, remoteVersion, isUpdateAvailable)
	if isUpdateAvailable {
		return true, remoteVersion
	}
	return false, remoteVersion
}

func UpdateAction(c *cli.Context, currentVersion string, prerelease bool) error {
	releaseAvailable, remoteVersion := CheckForUpdates(c, currentVersion, prerelease)
	pterm.Debug.Printfln("Release is available: %t. Remote version is %s", releaseAvailable, remoteVersion)
	if strings.Contains(currentVersion, "-DEV") {
		pterm.Error.Println("You can't run updates when running a dev build")
		return nil
	}
	if !releaseAvailable && !isUpdaterDisabled(c) {
		pterm.Info.Printfln("Sorry, no updates are available. The latest version is %s and you're on %s", remoteVersion, currentVersion)
		return nil
	}
	if isUpdaterDisabled(c) {
		return nil
	}
	return updater.UpgradeInPlace(c.App.Writer, c.App.ErrWriter, remoteVersion)
}
