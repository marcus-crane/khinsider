# khinsider

![](https://img.shields.io/badge/version-v2.0.0-green)

> Easily fetch videogame soundtracks from [downloads.khinsider.com](https://downloads.khinsider.com)

* [Installation](#installation)
  * [Homebrew](#homebrew)
  * [Docker](#docker)
  * [Go](#go)
  * [Binaries](#binaries)
* [Usage](#usage)
  * [Flags](#flags)
  * [Environment variables](#environment-variables)
* [Special Thanks](#special-thanks)

## Installation

V2 of khinsider offers a wide variety of distribution formats over the original:

### Homebrew

For macOS and Linux users who prefer Homebrew, there is a Homebrew formula installable like so:

```shell
brew install marcus-crane/tap/khinsider
```

### Docker

If you prefer to not run anything on your machine, you can run `khinsider` as a Docker image. It's available from both [Docker Hub](https://hub.docker.com/r/utf9k/khinsider) and [ghcr.io](https://github.com/marcus-crane/khinsider/pkgs/container/khinsider)

```shell
docker run utf9k/khinsider
```

### Go

If you have a relatively new version of Go, you can install `khinsider` like so:

```shell
go install github.com/marcus-crane/khinsider/v2@latest
```

### Binaries

There are a wide variety of binaries available under the [releases tab](https://github.com/marcus-crane/khinsider/releases):

- `apk`
- `deb`
- `rpm`
- `exe`

Windows, macOS and Linux are all supported (x86 and arm64) although at the time of writing, I haven't tested most of those platforms personally.

## Usage

> khinsider [global options] command [command options] [arguments...]

When you run `khinsider` by itself, you'll be presented with the help menu. There are a few subcommands you can choose from:

- `search`: Interactively filter through all of the albums on khinsider. It's powered by a [prebuilt index](https://github.com/marcus-crane/khinsider-index) so you can search the entire site at once.
  - The site is checked hourly for updates so the index will be no less than 1 hour out of date at any given time. If a new index is available, it'll be automatically downloaded.
- `album`: If you know the particular album you're already, you can provide a slug and download it straight away.
- `update`: As it says on the tin, you can automatically update `khinsider` to the latest version with this command.
- `index`: This command is unlisted but if you prefer to build a copy of the index locally, you can do so with this command. The prebuilt index is simply running this command and uploading the index generated if it differs from the previous copy.

At present, all albums are downloaded to `$HOME/Downloads`. It also isn't possible to download an album if a folder with the same name exists (ie; a previous download) in which case, you'll be prompted to move the existing folder in order to carry on.

### Flags

You can run `--debug` before any command to see some more detailed information in the case of an issue but it isn't as fleshed out as I would like it to be. This should be expanding in future releases.

When updating, you can use the `--prerelease` flag to request the latest prerelease over the latest stable release.

### Environment variables

By default, `khinsider` will check if there are any new updates in the background when run. It won't download them but it will prompt the user to consider updating. If you want to disable this check, you can set `KHINSIDER_NO_UPDATE=true` in your shell environment to disable this functionality.

## Special thanks

I wouldn't have originally made this project without being inspired by [obskyr](https://github.com/obskyr)'s original [python-based downloader](https://github.com/obskyr/khinsider).

In general, he's a cool guy and is up to lots of interesting stuff on [Twitter](https://twitter.com/obskyr).

Also a shout out to [Terin Stock](https://github.com/terinjokes) for his feedback on polishing up v2.0.0.
