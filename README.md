# khinsider

![A screenshot of khinsider running](screenshot.png?raw=true)

A usable but fairly featureless khinsider downloader written in Go. I use it myself but I need to clean up the code and write some more docs

### Table of Contents
* [Usage](#usage)
* [Installation](#installation)
  * [Prebuilt binaries](#prebuilt-binaries)
  * [Compiling from source](#compiling-from-source)
* [Special thanks](#special-thanks)

# Usage

As mentioned, this tool is very barebones so there's only really one thing you can do with it.

Let's say we wanted to download the OST for [Persona 4 Dancing All Night](https://downloads.khinsider.com/game-soundtracks/album/persona-4-dancing-all-night).

Assuming the URL is https://downloads.khinsider.com/game-soundtracks/album/persona-4-dancing-all-night, we want to take the slug, which is the portion of the URL after `/album/` like so:

```golang
khinsider persona-4-dancing-all-night
```

It will create a folder in your downloads folder (`$HOME/Downloads/`) named after the slug and then start to download each track so eg; `~/Downloads/persona-4-dancing-all-night/17 カリステギア Karisutegia.mp3`

There are no options for providing a download directory or anything like that but feel free to submit a feature request.

# Installation

There are two options for installing `khinsider`.

While both of these options will provide a binary in the download/compile directory, I recommend moving it to somewhere in your `PATH` such as `/usr/local/bin/khinsider`.

That way, you can access it going forward by just running `khinsider` and not having to specify eg; `~/Downloads/khinsider`

## Prebuilt binaries

Personally, I don't get off on the idea of compiling software so thanks to Github Actions, each release is already prebuilt and ready to go [on the releases page](https://github.com/marcus-crane/khinsider/releases).

I've provided builds for Windows, macOS and Linux, which contains a mix of both `x86` and `arm` binaries.

I do actually have an Apple Silicon Macbook Air which I'm pretty sure I tested the macOS binaries on but honestly I'd have to double check.

Let me know if there are any other platforms you'd like supported or feel free to add them yourself [here](https://github.com/marcus-crane/khinsider/blob/master/.github/workflows/release.yaml) by submitting a pull request.

## Compiling from source

This should just be as simple as the following:

```go
> go build main.go
> ./main
Please enter the name of an album eg katamari-damacy-soundtrack-katamari-fortissimo-damacy%
```

If you want to compile for a different operating system or architecture, just use the Golang compiler flags like so:

```go
> GOOS=linux GOARCH=arm64 go build main.go
> ./main
zsh: exec format error: ./main
> uname -ms
Darwin x86_64
```

I can't run the above example of course because I'm not running an arm based Linux machine but perhaps you'd like to compile for your Raspberry Pi while offline or something.

# Special thanks

If you're looking for something more feature complete, check out [obskyr](https://github.com/obskyr)'s original which inspired this one: https://github.com/obskyr/khinsider

In general, he's a cool guy and is up to lots of interesting stuff on [Twitter](https://twitter.com/obskyr)!

He didn't pay me to say this.
