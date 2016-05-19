package integration_tests

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Goa", func() {
	BeforeEach(func() {
		setupEnviroment()
	})
	AfterEach(func() {
		teardownEnvironment()
	})
	Context("package.path file present", func() {
		var session *gexec.Session
		BeforeEach(func() {
			createPackageFile("github.com/tinygrasshopper/x")

			Expect(os.Chdir(testDirectory)).To(Succeed())
			command := exec.Command(pathToGoaCli, "go", "help")
			var err error
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			GinkgoWriter.Write(session.Out.Contents())
		})

		It("creates a gopath directory structure in .go-anywhere/", func() {
			Expect(filepath.Join(testDirectory, ".go-anywhere")).To(BeADirectory())
			Expect(filepath.Join(testDirectory, ".go-anywhere", "src")).To(BeADirectory())
		})
		It("make pkg and bin directories", func() {
			Expect(filepath.Join(testDirectory, ".go-anywhere", "pkg")).To(BeADirectory())
			Expect(filepath.Join(testDirectory, ".go-anywhere", "bin")).To(BeADirectory())
		})

		It("creates the path specified in package.path in .go-anywhere/", func() {
			Expect(filepath.Join(testDirectory, ".go-anywhere", "src", "github.com/tinygrasshopper")).To(BeADirectory())
		})
		It("symlinks the current directory to the path specified by package.path", func() {
			linkPath, err := os.Readlink(filepath.Join(testDirectory, ".go-anywhere", "src", "github.com/tinygrasshopper/x"))
			Expect(err).NotTo(HaveOccurred())
			Expect(linkPath).To(HaveSuffix(testDirectory))
		})

		It("sets the GOPATH to .go-anywhere/", func() {
			Expect(mockExecutableHadEnvironment("GOPATH")).To(HaveSuffix(filepath.Join(testDirectory, ".go-anywhere")))
		})
	})
})

func absPath(path string) string {
	stat, err := filepath.Abs(path)
	Expect(err).NotTo(HaveOccurred())
	return stat
}

func createPackageFile(contents string) {
	file, err := os.OpenFile(filepath.Join(testDirectory, "package.path"), os.O_RDWR|os.O_CREATE, 0744)
	Expect(err).NotTo(HaveOccurred())
	file.Write([]byte(contents))
	Expect(file.Close()).To(Succeed())
	testFiles = append(testFiles, file.Name())
}
