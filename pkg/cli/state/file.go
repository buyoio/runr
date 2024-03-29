package state

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	cache map[string]string
)

func (f File) Path() string {
	return resolvePath(string(f))
}

func (f File) Content() (string, error) {
	if cache == nil {
		cache = make(map[string]string)
	}

	if _, ok := cache[f.Path()]; ok {
		return cache[f.Path()], nil
	}
	content, err := os.ReadFile(f.Path())
	if err != nil {
		return "", err
	}
	cache[f.Path()] = string(content)
	return cache[f.Path()], nil
}

func (f File) Output(envs []string) ([]byte, error) {
	cmd := exec.Command(f.Path())
	cmd.Env = []string{"PATH=" + os.Getenv("PATH")}
	for _, env := range envs {
		if !strings.Contains(env, "=") {
			env = fmt.Sprintf("%s=%s", env, os.Getenv(env))
		}
		cmd.Env = append(cmd.Env, env)
	}
	return cmd.Output()
}
