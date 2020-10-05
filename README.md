# promcurl

![Build](https://github.com/danielb42/promcurl/workflows/Build/badge.svg)
![Tag](https://img.shields.io/github/v/tag/danielb42/promcurl)
![Go Version](https://img.shields.io/github/go-mod/go-version/danielb42/promcurl)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/danielb42/promcurl)](https://pkg.go.dev/github.com/danielb42/promcurl)
[![Go Report Card](https://goreportcard.com/badge/github.com/danielb42/promcurl)](https://goreportcard.com/report/github.com/danielb42/promcurl)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Colorize Prometheus metric output on the terminal.
|&nbsp;|&nbsp;|
|-|-|
| `promcurl -u http://prmths/metrics`    | ![screenshot1](screen1.png) |
| `promcurl -u http://prmths/metrics -n` | ![screenshot2](screen2.png) |

## Install

### Either download the binary ...

Precompiled amd64-binaries are available for Linux and MacOS: [Latest Release](https://github.com/danielb42/promcurl/releases/latest)

### ... or build it yourself

`go get github.com/danielb42/promcurl`
