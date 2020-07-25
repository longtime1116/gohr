# gohr

`gohr` is a command line tool for hot-reloading(live-reloading) Go CLI application.
Just run `gohr` in the terminal, the program is compiled every time the files are updated, and also executed.

`gohr` means "go hot-reloading".

## Installation

```shell
go get github.com/longtime1116/gohr
```

Then test the installation

```shell
gohr -h
```

## Usage

```shell
gohr <binary name>
```

For example,

```shell
gohr main
```

If you ommit binary name, the basename of current directory is used.

## Demo
![gohr demo](https://raw.githubusercontent.com/wiki/longtime1116/gohr/images/gohr.gif)


