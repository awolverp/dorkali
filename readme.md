# Dorkali
A mini and simple tool to search in search engines written in golang.

- Supported engines:
    - Google

# Installation
```bash
GO111MODULE=off go get github.com/awolverp/dorkali/cmd/dorkali/
```
**Or**
```bash
git clone https://github.com/awolverp/dorkali/ && cd dorkali && go install ./cmd/dorkali/
```

## Usage
Use `dorkali help` to see usage:
```
Dorkali a program written in golang to dorks queries in search engines

Usage:
        dorkali [list | version [engineName] | help [engineName]]
        dorkali engineName [OPTIONS]

*Commands:
        version [engineName]   print version, or engine version if pass engineName, and exit
        list                   print list of engines and exit
        help [engineName]      print this help, or print engine help if pass engineName, and exit
```

For example if you want to see google engine help, you use `dorkali help google` command. you will see that:
```
Usage: dorkali google [OPTIONS] QUERY

*Output Options:
        ...

*Request Options:
        ...

*Search Options:
        ...

*Query Helpers:
        ...
```

## Example
```bash
$ dorkali google -n 5 "github"

https://desktop.github.com/
https://github.blog/2022-11-10-introducing-github-actions-importer/
http://en.wikipedia.org/wiki/C_(programming_language)
https://github.blog/
https://play.google.com/store/apps/details?id=com.github.android&hl=en&gl=US
```
