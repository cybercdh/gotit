## gotit

Take a list of URLs and prints the HTML title tag. This can be useful to perform quick flyovers of large lists of subdomains to yield interesting results.

## Recommended Usage

`$ cat subdomains | gotit`

or 

`$ assetfinder -subs-only example.com | gotit -c 50`

or

`$ gotit sub.example.com`

## Demo

[![asciicast](https://asciinema.org/a/J85aIYqUzqlYGV0fPNyPSa62p.svg)](https://asciinema.org/a/J85aIYqUzqlYGV0fPNyPSa62p)

## Install

If you have Go installed and configured (i.e. with $GOPATH/bin in your $PATH):

`go get -u github.com/cybercdh/gotit`

## Thanks

The concurrency techniques in this code was heavily inspired by @tomnomnom. That guy f**ks.
