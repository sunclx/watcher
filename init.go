package watcher

import "os/exec"

func init() {
	w := New([]string{"*.go"}, func() {
		exec.Command("go", "build")
	})
	w.watch(".")
	w.Run()
}
