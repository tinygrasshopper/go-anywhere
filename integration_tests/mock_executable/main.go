package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Printf(os.Getenv("RETURN_VALUE"))

	exitCode, err := strconv.Atoi(os.Getenv("RETURN_EXIT_CODE"))
	if err != nil {
		exitCode = 0
	}

	os.Exit(exitCode)
}
