# `todos`, a simple task management tool

Quickly find and read comments in the form `TODO(yourname): brief summary`.

## Installing

Run `go get github.com/benburwell/todos`.

## Usage

Basic usage:

    $ todos

You'll need to either pass the `--name=<yourname>` argument each time you run
`todos`, or create a `~/.todorc` file that just contains `name = "yourname"`.
