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

	defer os.Remove(runnerScript.Name())
	runnerScript.Write([]byte(fmt.Sprintf(`cd %s
%s "$@"`, linkedPath, fullName)))
	runnerScript.Close()

	os.Chdir(linkedPath)
	fmt.Printf("Switched to %s\n", linkedPath)
	wd, err := os.Getwd()
	fmt.Printf("WD %v %v \n", wd, err)
	fmt.Printf("Running %v %v\n", os.Args[1], os.Args[2:])

	cmd := exec.Command(runnerScript.Name(), os.Args[2:]...)

	err = cmd.Wait()
	data, err := cmd.CombinedOutput()
	fmt.Print(string(data))
	fmt.Print(err)
	if err != nil {
		panic(err) //TODO: handle
	}

}
