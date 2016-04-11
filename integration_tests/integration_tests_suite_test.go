package integration_tests

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestSystemTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SystemTests Suite")
}

var testDirectory string

var _ = BeforeEach(func() {
	testFiles = []string{}
	var err error
	testDirectory, err = ioutil.TempDir("", "go-anywhere")
	Expect(err).NotTo(HaveOccurred())
})
var _ = AfterEach(func() {
	Expect(os.RemoveAll(testDirectory)).To(Succeed())
})

var pathToGoaCli string
var pathToMockGoExecutable string
var testFiles []string

var _ = BeforeSuite(func() {
	var err error
	pathToGoaCli, err = gexec.Build("github.com/tinygrasshopper/go-anywhere/goa")
	Expect(err).ShouldNot(HaveOccurred())

	pathToMockGoExecutable, err = gexec.Build("github.com/tinygrasshopper/go-anywhere/integration_tests/mock_executable")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
