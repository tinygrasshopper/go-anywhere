package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	. "github.com/onsi/gomega"
)

func mockExecutableReturns(exitCode int, contents string) {
	os.Setenv("RETURN_VALUE", contents)
	os.Setenv("RETURN_EXIT_CODE", strconv.Itoa(exitCode))
}
func mockExecutableHadEnvironment(key string) string {
	file, err := os.Open(os.Getenv("TEST_ENV_PASSED"))
	Expect(err).NotTo(HaveOccurred())
	decoder := json.NewDecoder(file)
	var env []string
	decoder.Decode(&env)

	for _, v := range env {
		if strings.HasPrefix(v, key+"=") {
			return strings.Split(v, "=")[1]
		}
	}
	return ""
}
func setupEnviroment() {
	os.Setenv("TEST_ARGUMENTS_PASSED", tempFile())
	os.Setenv("TEST_ENV_PASSED", tempFile())
}

func teardownEnvironment() {
	Expect(os.RemoveAll(os.Getenv("TEST_ENV_PASSED")))
	Expect(os.RemoveAll(os.Getenv("TEST_ARGUMENTS_PASSED")))
}

func tempFile() string {
	path, err := ioutil.TempFile("", "mock_executable")
	Expect(err).NotTo(HaveOccurred())
	return path.Name()
}
