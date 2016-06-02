package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	pwd := projectRoot()
	//TODO: args

	packagePath := packagePath()
	linkedPath := pwd + "/.go-anywhere/src/" + packagePath
	command := os.Args[1]

	//TODO: remove script
	var runner *os.File

	switch command {
	case "exec":
		runner = runnerScript(linkedPath, os.Args[2], os.Args[3:]...)
	case "build", "clean", "doc", "env", "fix", "fmt", "generate", "get", "install", "list", "run", "test", "tool", "version", "vet":
		//TODO: remove get from this list
		runner = runnerScript(linkedPath, "go", os.Args[1:]...)
	default:
		panic(fmt.Sprintf("unknown command %s", command))
	}

	ensurePath(pwd, packagePath, linkedPath)
	setGoPath(pwd)
	prependGoaToPath(os.Args[0])

	cmd := exec.Command("/bin/sh", runner.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		panic(err) //TODO: handle
	}
}

func projectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err) //TODO: handle
	}
	return cwd
}

func packagePath() string {
	//TODO: try to find the file recursively till root
	file, err := os.Open("package.path") //TODO: make filename configurable
	if err != nil {
		panic(err) //TODO: handle
	}

	packagePathBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err) //TODO: handle
	}
	return strings.TrimSpace(string(packagePathBytes))
}

func ensurePath(pwd, packagePath, linkedPath string) {
	if _, err := os.Stat(".go-anywhere"); os.IsNotExist(err) {
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
}

func setGoPath(pwd string) {
	err := os.Setenv("GOPATH", pwd+"/.go-anywhere")
	if err != nil {
		panic(err) //TODO: handle
	}
}

func runnerScript(linkedPath, executable string, args ...string) *os.File {
	fullName, err := exec.LookPath(executable)
	if err != nil {
		panic(err) //TODO: handle
	}

	runnerScript, err := ioutil.TempFile("", "go-anywhere-runner")
	if err != nil {
		panic(err) //TODO: handle
	}

	runnerScript.WriteString(`cd ` + linkedPath + "\n")
	runnerScript.WriteString(fullName + " " + strings.Join(args, " "))

	// TODO: Redcue permissions
	err = runnerScript.Chmod(0700)
	if err != nil {
		panic(err) //TODO: handle
	}
	return runnerScript
}

func prependGoaToPath(goa string) {
	dir, err := ioutil.TempDir("", "go-path")
	if err != nil {
		panic(err) //TODO: handle
	}

	fullName, err := exec.LookPath(goa)
	if err != nil {
		panic(err) //TODO: handle
	}

	err = os.Symlink(fullName, path.Join(dir, "go"))
	if err != nil {
		panic(err) //TODO: handle
	}

	path := os.Getenv("PATH")
	fmt.Printf("dir: %s\n", dir)
	os.Setenv("PATH", dir+":"+path)
}
