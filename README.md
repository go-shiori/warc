WARC
===

[![GoDoc](https://godoc.org/github.com/go-shiori/warc?status.png)](https://godoc.org/github.com/go-shiori/warc)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-shiori/warc)](https://goreportcard.com/report/github.com/go-shiori/warc)

**This project is now archived**. If you want to archive, consider checking out [`obelisk`](https://github.com/go-shiori/obelisk). It has better output format (plain HTML) and IMHO better written than this.

WARC is a Go package that archive a web page and its resources into a single [`bolt`](https://github.com/etcd-io/bbolt) database file. Developed as part of [Shiori](https://github.com/go-shiori/shiori) bookmarks manager.

It still in development phase but should be stable enough to use. The `bolt` database that used by this project is also stable both in API and file format. Unfortunately, right now WARC will disable Javascript when archiving a page so it still doesn't not work in SPA site like Twitter or Reddit.

## Installation

To install this package, just run `go get` :

```
go get -u -v github.com/go-shiori/warc
```

## Licenses

WARC is distributed under [MIT license](https://choosealicense.com/licenses/mit/), which means you can use and modify it however you want. However, if you make an enhancement for it, if possible, please send a pull request.
