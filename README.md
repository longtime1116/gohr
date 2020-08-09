# gohr

`gohr` is a command line tool for hot-reloading(live-reloading) Go CLI application.
Just run `gohr` in the terminal, the program is compiled every time the files are updated, and also executed.

`gohr` means "go hot-reloading".

## Installation

```shell
go get github.com/longtime1116/gohr
```

Then test the installation.

```shell
gohr -h
Usage: gohr [OPTIONS] [<output binary name>]
When you ommit output binary name, the basename of current directory is used.
  -b    alias of --build-only
  -build-only
        Just only build and not execute command.
```

## Usage

```shell
gohr <output binary name>
```

For example,

```shell
gohr main
```

If you ommit output binary name, the basename of current directory is used.

Also you can add `--build-only` or `-b` option.
With this option, `gohr` just build the program and don't execute it.

```shell
gohr --build-only main
```

## Demo

![gohr demo](https://raw.githubusercontent.com/wiki/longtime1116/gohr/images/gohr.gif)
