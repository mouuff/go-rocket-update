# go-rocket-update: Build self-updating Go programs

[![Build Status](https://github.com/mouuff/go-rocket-update/workflows/Go/badge.svg?branch=master)](https://github.com/mouuff/go-rocket-update/actions)
[![codecov](https://codecov.io/gh/mouuff/go-rocket-update/branch/master/graph/badge.svg)](https://codecov.io/gh/mouuff/go-rocket-update)
[![Go Report Card](https://goreportcard.com/badge/github.com/mouuff/go-rocket-update)](https://goreportcard.com/report/github.com/mouuff/go-rocket-update)
[![Go Reference](https://pkg.go.dev/badge/github.com/mouuff/go-rocket-update.svg)](https://pkg.go.dev/github.com/mouuff/go-rocket-update)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Enable your Golang applications to easily and safely to self update.

Here is the list of [projects using this package](https://github.com/mouuff/go-rocket-update/network/dependents?package_id=UGFja2FnZS0yMjc3OTEzNjc1).

It provides the flexibility to implement different updating user experiences like auto-updating, or manual user-initiated updates, and updates from different sources.

![Go rocket image](docs/social.png)
_The gopher in this image was created by [Takuya Ueda][tu], licensed under [Creative Commons 3.0 Attributions license][cc3-by]._

## Features

- Flexible way to provide updates (ex: using Github or Gitlab!)
- Cross platform support (Mac, Linux, Arm, and Windows)
- RSA signature verification
- Tooling to generate and verify signatures
- Background update
- Rollback feature

## QuickStart

### Install library

`go get -u github.com/mouuff/go-rocket-update/...`

### Enable your App to Self Update

Here is an example using Github releases:

``` go
u := &updater.Updater{
	Provider: &provider.Github{
		RepositoryURL: "github.com/mouuff/go-rocket-update-example",
		ArchiveName:   fmt.Sprintf("binaries_%s.zip", runtime.GOOS),
	},
	ExecutableName: fmt.Sprintf("go-rocket-update-example_%s_%s", runtime.GOOS, runtime.GOARCH),
	Version:        "v0.0.1",
}

if _, err := u.Update(); err != nil {
	log.Println(err)
}
```

For more examples, please take a look at some [code samples](./examples) and this [example project](https://github.com/mouuff/go-rocket-update-example).

### Push an update

The updater uses a `Provider` as an input source for updates. It provides files and version for the updater.

Here is few examples of providers:

- `provider.Github`: It will check for the latest release on Github with a specific archive name (zip or tar.gz)
- `provider.Gitlab`: It will check for the latest release on Gitlab with a specific archive name (zip or tar.gz)
- `provider.Local`: It will use a local folder, version will be defined in the VERSION file (can be used for testing, or in a company with a shared folder for example)
- `provider.Zip`: It will use a `zip` file. The version is defined by the file name (Example: `binaries-v1.0.0.tar.gz`). Use [GlobNewestFile](https://github.com/mouuff/go-rocket-update/blob/0cad960c4449b42726537e2c559786b3d6174868/pkg/provider/common.go#L24) to find the right file.
- `provider.Gzip`: Same as `provider.Zip` but with a `tar.gz` file.

The updater will list the files and retrieve them the same way for all the providers:

The directory should have files containing `ExecutableName`.

Example directory content with `ExecutableName: fmt.Sprintf("test_%s_%s", runtime.GOOS, runtime.GOARCH)`:

    test_windows_amd64.exe
    test_darwin_amd64
    test_linux_arm

We recommend using [goxc](https://github.com/laher/goxc) for compiling your Go application for multiple platforms.

### Important notes
- To update the binary, you need to have the right permissions for the folder where it is installed. For example, if the binary is in a folder like "Program Files", the process will need to acquire admin permissions.

### Planned features

This project is currently under construction, here is some of the things to come:

- More documentation and examples
- [Variable templating](https://github.com/mouuff/go-rocket-update/issues/14)
- Mutliple providers (enables the use of another provider if the first one is down)
- Update channels for Github provider (alpha, beta, ...)
- Validation of the executable being installed

[tu]: https://twitter.com/tenntenn
[cc3-by]: https://creativecommons.org/licenses/by/3.0/

