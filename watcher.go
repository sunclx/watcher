package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	patterns []string
	action   func()
	w        *fsnotify.Watcher
	ec       chan fsnotify.Event

	errc chan error
}

func New(patterns []string, action func()) *Watcher {
	var w Watcher
	w.patterns = patterns
	w.action = action

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	w.w = watcher

	w.ec = make(chan fsnotify.Event, 4)

	w.errc = make(chan error)
	go func() {
		for {
			fmt.Println(<-w.errc)
		}
	}()
	return &w
}

func (w *Watcher) watch(path string) {
	go func() {
		defer w.w.Close()
		for {
			select {
			case event := <-w.w.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Lstat(event.Name)
					w.errc <- err
					if fi.IsDir() {
						w.w.Add(event.Name)
					}
				}
				w.ec <- event
			case err := <-w.w.Errors:
				w.errc <- err
				return
			}
		}
	}()

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			w.w.Add(path)
		}
		return nil
	})
}

func (w *Watcher) Run() {
	timer := time.NewTimer(time.Second * 1)
	if !timer.Stop() {
		<-timer.C
	}
	for {
		select {
		case event := <-w.ec:
			if !w.match(event) {
				break
			}
			//ens = append(ens, event)
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time.Second * 1)
		case <-timer.C:
			w.action()
		}
	}
}

func (w *Watcher) match(e fsnotify.Event) bool {
	for _, pattern := range w.patterns {
		if matched, err := filepath.Match(pattern, e.Name); matched && err == nil {
			return true
		}
	}
	return false
}

type Command struct {
	Name   string
	Args   []string
	cmd    *exec.Cmd
	mutex  *sync.Mutex
	exited chan struct{}
}

func NewCommand(name string, args ...string) *Command {
	return &Command{
		Name:  name,
		Args:  args,
		cmd:   exec.Command(name, args...),
		mutex: new(sync.Mutex),
	}
}
func (c *Command) Start() {

}
