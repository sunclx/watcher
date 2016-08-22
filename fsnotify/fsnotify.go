package fsnotify

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sunclx/watcher"
)

type Watcher struct {
	watcher.ProjectService
	w *fsnotify.Watcher
}

func NewWatcher() *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	w := &Watcher{w: watcher}
	go func() {
		for event := range watcher.Events {
			if event.Op&fsnotify.Create == fsnotify.Create {
				fi, _ := os.Lstat(event.Name)
				if fi.IsDir() {
					watcher.Add(event.Name)
				}
			}
		}
	}()
	return w
}
func (w *Watcher) Add(path string) {
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

type Filter struct {
	patterns []string
	Events   chan []fsnotify.Event
	Errors   chan error
}

func NewFilter(w *Watcher, patterns ...string) *Filter {
	f := &Filter{
		patterns: patterns,
		Events:   make(chan []fsnotify.Event),
		Errors:   w.Errors,
	}
	events := make([]fsnotify.Event, 0, 4)
	go func() {
		timer := time.NewTimer(time.Second * 1)
		for {
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time.Second * 1)

			select {
			case event := <-w.Events:
				for _, pattern := range patterns {
					if matched, err := filepath.Match(pattern, event.Name); matched && err == nil {
						events = append(events, event)
					}
				}

			case <-timer.C:
				f.Events <- events
				events = events[:0]
			}
		}
	}()
	return f
}
