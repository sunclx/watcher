package main

import (
	"os"

	"github.com/sunclx/watcher"
	"github.com/sunclx/watcher/command"
	"github.com/sunclx/watcher/fsnotify"
)

func main() {

	prj := watcher.Project{
		Name: "watcher",
		Path: os.Getenv("GOPATH") + "/github.com/sunclx/watcher/cmd",
	}
	runner := command.NewRunner(prj)
	watcher := fsnotify.NewWatcher(prj, runner)

	watcher.Watching()

}
