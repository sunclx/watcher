package fsnotify

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sunclx/watcher"
)

type Watcher struct {
	Runner watcher.Runner
	p      watcher.Project
	w      *fsnotify.Watcher
	events chan fsnotify.Event
}

func NewWatcher(p watcher.Project, r watcher.Runner) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	w := &Watcher{
		Runner: r,
		p:      p,
		w:      watcher,
		events: make(chan fsnotify.Event),
	}
	go func() {
		for event := range watcher.Events {
			if event.Op&fsnotify.Create == fsnotify.Create {
				fi, _ := os.Lstat(event.Name)
				if fi.IsDir() {
					watcher.Add(event.Name)
				}
			}
			if matched, _ := filepath.Match("*.go", event.Name); matched {
				w.events <- event
			}
		}
	}()
	w.add(p.Path)
	return w
}
func (w *Watcher) add(path string) {
	fi, err := os.Lstat(path)
	if err != nil {
		return
	}
	if !fi.IsDir() {
		w.w.Add(path)
		return
	}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			w.w.Add(path)
		}
		return nil
	})
}

func (w *Watcher) Watching() {
	err := w.Runner.Start()
	fmt.Println(err)

	timer := time.NewTimer(time.Second * 1)
	start := false
	for {
		if !timer.Stop() {
			<-timer.C
		}
		if start {
			timer.Reset(time.Second * 1)
		}
		select {
		case <-w.events:
			start = true
		case <-timer.C:
			err = w.Runner.Stop()
			fmt.Println(err)
			err = w.Runner.Start()
			fmt.Println(err)
			fmt.Println("restart success")
		}
	}

}
