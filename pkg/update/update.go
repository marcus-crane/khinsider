package update

import (
	"errors"
	"fmt"
	"github.com/marcus-crane/khinsider/v2/pkg/types"
	"github.com/marcus-crane/khinsider/v2/pkg/util"
	"github.com/pterm/pterm"
	"golang.org/x/mod/semver"
	"strings"
)

const (
	AppReleaseFeed    = "https://api.github.com/repos/marcus-crane/khinsider/releases/latest"
	AppPrereleaseFeed = "https://api.github.com/repos/marcus-crane/khinsider/releases"
	IndexReleaseFeed  = "https://api.github.com/repos/marcus-crane/khinsider-index/releases/latest"
)

func GetRemoteIndexVersion() string {
	release, err := GetGithubRelease(IndexReleaseFeed)
	if err != nil {
		return ""
	}
	return ValidateIndexVersion(release.Version, "remote")
}

func GetRemoteAppVersion() string {
	release, err := GetGithubRelease(AppReleaseFeed)
	if err != nil {
		return ""
	}
	return ValidateIndexVersion(release.Version, "app")
}

func GetRemoteAppPrerelease() string {
	release, err := GetGithubPrerelease(AppPrereleaseFeed)
	if err != nil {
		return ""
	}
	return ValidateIndexVersion(release.Version, "app")
}

func GetGithubRelease(releaseFeed string) (types.RemoteIndexMetadata, error) {
	release := types.RemoteIndexMetadata{}
	res, err := util.RequestJSON(releaseFeed)
	if err != nil {
		return release, err
	}
	defer res.Body.Close()
	if err := util.LoadJSON(res.Body, &release); err != nil {
		return release, err
	}
	return release, nil
}

func GetGithubPrerelease(releaseFeed string) (types.RemoteIndexMetadata, error) {
	var releaseList []types.RemoteIndexMetadata
	res, err := util.RequestJSON(releaseFeed)
	if err != nil {
		return types.RemoteIndexMetadata{}, err
	}
	defer res.Body.Close()
	if err := util.LoadJSON(res.Body, &releaseList); err != nil {
		return types.RemoteIndexMetadata{}, err
	}
	var latestPrerelease types.RemoteIndexMetadata
	for _, entry := range releaseList {
		if entry.Prerelease == true {
			latestPrerelease = entry
			break
		}
	}
	return latestPrerelease, nil
}

func ValidateIndexVersion(version string, indexLocation string) string {
	if !strings.HasPrefix(version, "v") {
		pterm.Error.Printfln("%s index version %s doesn't start with a v.", indexLocation, version)
		panic(errors.New(fmt.Sprintf("%s index version is invalid", indexLocation)))
	}
	versionBits := strings.Split(version, ".")
	if len(versionBits) != 3 {
		pterm.Error.Printf("expected %s version %s to have 3 parts. %s only has %d", indexLocation, version, len(versionBits))
		panic(errors.New(fmt.Sprintf("%s index version is invalid", indexLocation)))
	}
	return version
}

func IsRemoteVersionNewer(localVersion string, remoteVersion string) bool {
	result := semver.Compare(localVersion, remoteVersion)
	if result == -1 {
		return true
	}
	return false
}
