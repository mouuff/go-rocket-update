

# go-rocket-update: Build self-updating Go programs

[![Build Status](https://github.com/mouuff/go-rocket-update/workflows/Go/badge.svg?branch=master)](https://github.com/mouuff/go-rocket-update/actions)
[![codecov](https://codecov.io/gh/mouuff/go-rocket-update/branch/master/graph/badge.svg)](https://codecov.io/gh/mouuff/go-rocket-update)
[![Go ReportCard](http://goreportcard.com/badge/mouuff/go-rocket-update)](http://goreportcard.com/report/mouuff/go-rocket-update)
[![Go Reference](https://pkg.go.dev/badge/github.com/mouuff/go-rocket-update.svg)](https://pkg.go.dev/github.com/mouuff/go-rocket-update)


Enable your Golang applications easily and safely to self update.

It provides the flexibility to implement different updating user experiences like auto-updating, or manual user-initiated updates, and updates from different sources.

![Go rocket image](docs/social.png)
*The gopher in this image was created by [Takuya Ueda][tu], licensed under [Creative Commons 3.0 Attributions license][cc3-by].*

## Features
* Flexible way to provide updates (ex: using Github or Gitlab!)
* Cross platform support (Mac, Linux, Arm, and Windows)
* RSA signature verification
* Tooling to generate and verify signatures
* Background update
* Rollback feature

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
		ExecutableName: "go-rocket-update-example",
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
* `provider.Github`: It will check for the latest release on Github with a specific zip name
* `provider.Gitlab`: It will check for the latest release on Gitlab with a specific zip name
* `provider.Local`: It will use a local folder, version will be defined in the VERSION file (can be used for testing, or in a company with a shared folder for example)
* `provider.Zip`: Same as provider.Local but with a `Zip` file

*In the future there will be providers for FTP servers and Google cloud storage.*

The updater will list the files and retrieve them the same way for all the providers:

The directory should contain files with the name: ExecutableName-$GOOS-$ARCH.

Example with ExecutableName `test`:

    test-windows-386
    test-darwin-amd64
    test-linux-arm

We recommend using [goxc](https://github.com/laher/goxc) for compiling your Go application for multiple platforms.

### Planned features
This project is currently under construction, here is some of the things to come:
* More documentation and examples
* Google cloud storage and FTP providers
* Mutliple providers (enables the use of another provider if the first one is down)
* Update channels for Github provider (alpha, beta, ...)



## API Breaking Changes
- **Feb 7, 2021**: The `BinaryName` variable used in `Updater` have been renamed to `ExecutableName`.
- **Feb 12, 2021**: The method `Updater.Update()` now returns `(UpdateStatus, error)` instead of just `(error)`.


[tu]: https://twitter.com/tenntenn
[cc3-by]: https://creativecommons.org/licenses/by/3.0/
