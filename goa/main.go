package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	linkedPath := pwd + "/.go-anywhere/src/" + packagePath

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
		err = os.Symlink(pwd, linkedPath)
		if err != nil {
			panic(err) //TODO: handle
		}
	}

	err = os.Setenv("GOPATH", pwd+"/.go-anywhere")
	if err != nil {
		panic(err) //TODO: handle
	}

	fullName, err := exec.LookPath(os.Args[1])
	if err != nil {
		panic(err) //TODO: handle
	}

	runnerScript, err := ioutil.TempFile("", "go-anywhere-runner")
	if err != nil {
		panic(err) //TODO: handle
	}

	runnerScript.WriteString(`cd ` + linkedPath + "\n")
	runnerScript.WriteString(fullName + " " + strings.Join(os.Args[2:], " "))

	// TODO: Redcue permissions
	err = runnerScript.Chmod(0700)
	if err != nil {
		panic(err) //TODO: handle
	}

	cmd := exec.Command("/bin/sh", runnerScript.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		panic(err) //TODO: handle
	}
}
