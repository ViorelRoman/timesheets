package main

import (
	"os/exec"
)

type Repo interface {
	GetLog(user string) (string, error)
}

type HGRepo struct {
	path string
}

func (r HGRepo) GetLog(user string) (string, error) {
	c := exec.Command("hg", "log", "--user", user, "--pager", "never", "--template", `'{date}|{desc}\n`)
	c.Dir = r.path
	output, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type GitRepo struct {
	path string
}

func (r GitRepo) GetLog(user string) (string, error) {
	c := exec.Command("git", "log", "--author", user, `--pretty=format:"%at|%s`)
	c.Dir = r.path
	output, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
