package command

import (
	"errors"
	"os"
	"os/exec"

	"github.com/sunclx/watcher"
)

type Runner struct {
	Project watcher.Project
	run     *exec.Cmd
}

func NewRunner(p watcher.Project) watcher.Runner {
	return &Runner{
		Project: p,
	}
}
func (c *Runner) Start() error {
	if c.Project.Name == "" || c.Project.Path == "" {
		return errors.New("empty Project")
	}

	path := c.Project.Path
	runCmd := func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		cmd.Dir = path
		return cmd.Run()
	}

	// err := runCmd("git", "pull")
	// if err != nil {
	// 	return err
	// }

	// err = runCmd("go", "get")
	// if err != nil {
	// 	return err
	// }

	filepath := os.TempDir() + c.Project.Name
	err := runCmd("go", "build", "-o", filepath)
	if err != nil {
		return err
	}
	os.Remove(filepath)

	filepath = os.Getenv("GOPATH") + "/bin" + c.Project.Name
	err = runCmd("go", "build", "-o", filepath)
	if err != nil {
		return err
	}

	run := exec.Command(filepath)
	err = run.Start()
	if err != nil {
		return err
	}
	c.run = run
	return nil
}
func (c *Runner) Stop() error {
	return c.run.Process.Kill()
}
