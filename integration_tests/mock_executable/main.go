package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	fmt.Printf(os.Getenv("RETURN_VALUE"))

	exitCode, _ := strconv.Atoi(os.Getenv("RETURN_EXIT_CODE"))

	data, _ := json.Marshal(os.Args)
	ioutil.WriteFile(os.Getenv("TEST_ARGUMENTS_PASSED"), data, 0666)

	data, _ = json.Marshal(os.Environ())
	ioutil.WriteFile(os.Getenv("TEST_ENV_PASSED"), data, 0666)

	os.Exit(exitCode)
}
