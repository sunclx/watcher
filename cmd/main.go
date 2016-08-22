package main

import (
	"fmt"
	"os"

	"github.com/sunclx/watcher"
	"github.com/sunclx/watcher/command"
	"github.com/sunclx/watcher/fsnotify"
)

func main() {
	fmt.Println("starting...")
	prj := watcher.Project{
		Name: "watcher",
		Path: os.Getenv("GOPATH") + "/src/github.com/sunclx/watcher",
		Main: "cmd",
	}
	runner := command.NewRunner(prj)
	watcher := fsnotify.NewWatcher(prj, runner)
	fmt.Println("watching...")
	watcher.Watching()

}
