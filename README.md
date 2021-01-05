

# go-rocket-update: Build self-updating Go programs

[![Build Status](https://github.com/mouuff/go-rocket-update/workflows/Go/badge.svg?branch=master)](https://github.com/mouuff/go-rocket-update/actions)
[![codecov](https://codecov.io/gh/mouuff/go-rocket-update/branch/master/graph/badge.svg)](https://codecov.io/gh/mouuff/go-rocket-update)

Enable your Golang applications easily and safely to self update.

It provides the flexibility to implement different updating user experiences like auto-updating, or manual user-initiated updates, and updates from different sources.

![Go rocket image](docs/social.png)
*The gopher in this image was created by [Takuya Ueda][tu], licensed under [Creative Commons 3.0 Attributions license][cc3-by].*

## Features
* Flexible way to provide updates (ex: using Github!)
* Cross platform support (Mac, Linux, Arm, and Windows)
* RSA signature verification
* Tooling to generate and verify signatures
* Background update

## QuickStart

### Install library

`go get -u github.com/mouuff/go-rocket-update/...`

### Enable your App to Self Update

Here is an example using Github releases:

	u := &updater.Updater{
		Provider: &provider.Github{
			RepositoryURL: "github.com/mouuff/go-rocket-update-example",
			ZipName:       "binaries_" + runtime.GOOS + ".zip",
		},
		BinaryName: "go-rocket-update-example",
		Version:    "v0.1",
	}
	log.Println(u.Version)
	err := u.Update()
	if err != nil {
		log.Error(err)
	}

Check this project for a complete example: https://github.com/mouuff/go-rocket-update-example

### Push an update

The updater uses a `Provider` as an input source for updates. It provides files and version for the updater.

Here is few examples of providers:
* `provider.Github`: It will check for the lastest release on Github with a specific zip name
* `provider.Local`: It will use a local folder, version will be defined in the VERSION file (can be used for testing, or in a company with a shared folder for example)
* `provider.Zip`: Same as provider.Local but with a `Zip` file

*In the future there will be providers for FTP servers and Gitlab.*

The updater will list the files and retrieve them the same way for all the providers:

The directory should contain files with the name: BinaryName-$GOOS-$ARCH.

Example with BinaryName `progname`:

    progname-windows-386
    progname-darwin-amd64
    progname-linux-arm

We recommend using [goxc](https://github.com/laher/goxc) for compiling your Go application for multiple platforms.

### Planned features
This project is currently under construction, here is some of the things to come:
* More documentation and examples
* Gitlab and FTP providers
* Mutliple providers (enables the use of another provider if the first one is down)
* Update channels for Github provider (alpha, beta, ...)
* Rollback feature

[tu]: https://twitter.com/tenntenn
[cc3-by]: https://creativecommons.org/licenses/by/3.0/
