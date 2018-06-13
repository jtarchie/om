package acceptance

import (
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("interpolate command", func() {
	Context("when given a valid YAML file", func() {
		createFile := func(contents string) *os.File {
			file, err := ioutil.TempFile("", "")
			Expect(err).ToNot(HaveOccurred())
			_, err = file.WriteString(contents)
			Expect(err).ToNot(HaveOccurred())
			return file
		}

		It("outputs a YAML file", func() {
			yamlFile := createFile("---\nname: bob\nage: 100")
			command := exec.Command(pathToMain,
				"interpolate", "-c", yamlFile.Name(),
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5).Should(gexec.Exit(0))
			Eventually(session.Out, 5).Should(gbytes.Say("age: 100\nname: bob"))
		})

		Context("when an output file is provided", func() {
			It("writes the YAML to that file", func() {
				yamlFile := createFile("---\nname: bob\nage: 100")
				outputFile, err := ioutil.TempFile("", "")
				Expect(err).ToNot(HaveOccurred())

				err = os.Remove(outputFile.Name())
				Expect(err).ToNot(HaveOccurred())

				command := exec.Command(pathToMain,
					"interpolate",
					"-c", yamlFile.Name(),
					"-o", outputFile.Name(),
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session, 5).Should(gexec.Exit(0))

				fileInfo, err := os.Stat(outputFile.Name())
				Expect(err).ToNot(HaveOccurred())
				Expect(fileInfo.Mode()).To(Equal(os.FileMode(0600)))

				contents, err := ioutil.ReadFile(outputFile.Name())
				Expect(err).ToNot(HaveOccurred())
				Expect(contents).To(MatchYAML("age: 100\nname: bob"))
				Consistently(session.Out, 5).ShouldNot(gbytes.Say("age: 100\nname: bob"))
			})

			It("errors when the file does not exist", func() {
				yamlFile := createFile("---\nname: bob\nage: 100")
				outputFile, err := ioutil.TempFile("", "")
				Expect(err).ToNot(HaveOccurred())
				outputFile.Close()
				Expect(os.Remove(outputFile.Name())).ToNot(HaveOccurred())

				command := exec.Command(pathToMain,
					"interpolate",
					"-c", yamlFile.Name(),
					"-o", outputFile.Name(),
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session, 5).Should(gexec.Exit(0))

				contents, err := ioutil.ReadFile(outputFile.Name())
				Expect(err).ToNot(HaveOccurred())
				Expect(contents).To(MatchYAML("age: 100\nname: bob"))
				Consistently(session.Out, 5).ShouldNot(gbytes.Say("age: 100\nname: bob"))
			})
		})

		Context("with vars defined in the manifest", func() {
			It("successfully replaces the vars", func() {
				varsFile := createFile("---\nname1: moe\nage1: 500")
				yamlFile := createFile("---\nname: ((name1))\nage: ((age1))")
				command := exec.Command(pathToMain,
					"interpolate",
					"-c", yamlFile.Name(),
					"-l", varsFile.Name(),
				)
				defer varsFile.Close()
				defer yamlFile.Close()

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5).Should(gexec.Exit(0))
				Eventually(session.Out, 5).Should(gbytes.Say("age: 500\nname: moe"))
			})

			It("replaces the vars based on the order precedence of the vars file", func() {
				vars1File := createFile("---\nname1: moe\nage1: 500")
				vars2File := createFile("---\nname1: bob")
				yamlFile := createFile("---\nname: ((name1))\nage: ((age1))")
				command := exec.Command(pathToMain,
					"interpolate",
					"-c", yamlFile.Name(),
					"-l", vars1File.Name(),
					"-l", vars2File.Name(),
				)
				defer vars1File.Close()
				defer vars2File.Close()
				defer yamlFile.Close()

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5).Should(gexec.Exit(0))
				Eventually(session.Out, 5).Should(gbytes.Say("age: 500\nname: bob"))
			})

			It("errors when no vars are provided", func() {
				yamlFile := createFile("---\nname: ((name1))\nage: ((age1))")
				command := exec.Command(pathToMain,
					"interpolate",
					"-c", yamlFile.Name(),
				)
				defer yamlFile.Close()

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5).Should(gexec.Exit(1))
				Eventually(session.Err, 5).Should(gbytes.Say("Expected to find variables: age1"))
			})
		})
	})
})
