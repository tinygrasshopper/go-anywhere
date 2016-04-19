package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err) //TODO: handle
	}
	//TODO: args

	//TODO: try to find the file recursively till root
	file, err := os.Open("package.path") //TODO: make filename configurable
	if err != nil {
		panic(err) //TODO: handle
	}

	packagePathBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err) //TODO: handle
	}
	packagePath := strings.TrimSpace(string(packagePathBytes))

	if _, err = os.Stat(".go-anywhere"); os.IsNotExist(err) {
		err = os.MkdirAll(".go-anywhere/pkg/", 0744|os.ModeDir)
		if err != nil {
			panic(err) //TODO: handle
		}

		err = os.MkdirAll(".go-anywhere/bin/", 0744|os.ModeDir)
		if err != nil {
			panic(err) //TODO: handle
		}

		err := os.MkdirAll(".go-anywhere/src/"+filepath.Dir(packagePath), 0744|os.ModeDir)
		if err != nil {
			panic(err) //TODO: handle
		}
		err = os.Symlink(pwd, pwd+"/.go-anywhere/src/"+packagePath)
		if err != nil {
			panic(err) //TODO: handle
		}
	}

	err = os.Setenv("GOPATH", pwd+"/.go-anywhere")
	if err != nil {
		panic(err) //TODO: handle
	}

	cmd := exec.Command("go", os.Args[1:]...)
	err = cmd.Start()
	if err != nil {
		panic(err) //TODO: handle
	}

	go io.Copy(cmd.Stdout, os.Stdout)
	go io.Copy(cmd.Stderr, os.Stderr)

	err = cmd.Wait()
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		} else {
			panic(err)
		}
	}

}
