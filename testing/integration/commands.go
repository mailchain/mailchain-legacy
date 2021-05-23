package integration

import (
	"os"
	"os/exec"
	"syscall"
)

var path = "../../cmd/mailchain/main.go"

func bundle(items ...string) []string {
	return items
}

func createCommand(configLocation string, commands []string, args []string) *exec.Cmd {
	combined := []string{
		"run", path, "--config", configLocation,
	}
	combined = append(combined, append(commands, args...)...)

	c := exec.Command(
		"go",
		combined...,
	)
	c.Env = os.Environ()
	c.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	return c
}
