package watcher

import (
	"io"
	"os"
)

var (
	Stdin  io.Reader = os.Stdin  // panic(interface{})
	Stdout io.Writer = os.Stdout // panic(interface{})
	Stderr io.Writer = os.Stderr // panic(interface{})
)

type Project struct {
	Name string
	Path string
	Main string
}

type ProjectService interface {
	Project(name string) (Project, error)
	Projects() ([]Project, error)
	Add(Project) error
	Remove(Project) error
}

type Runner interface {
	Start() error
	Stop() error
}
