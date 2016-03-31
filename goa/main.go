package main

import (
	"fmt"
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

	if _, err := os.Stat(".go-anywhere"); err == os.ErrNotExist {
		err = os.MkdirAll(".go-anywhere/src/"+filepath.Dir(packagePath), 0744|os.ModeDir)
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
	err = cmd.Wait()
	data, err := cmd.CombinedOutput()
	fmt.Print(string(data))
	fmt.Print(err)
	if err != nil {
		panic(err) //TODO: handle
	}

}
