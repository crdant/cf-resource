package out_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/crdant/cf-route-resource"
	"github.com/crdant/cf-route-resource/out"
)

var _ = Describe("Out", func() {
	var (
		tmpDir  string
		cmd     *exec.Cmd
		request out.Request
	)

	BeforeEach(func() {
		var err error

		tmpDir, err = ioutil.TempDir("", "cf_resource_out")
		Expect(err).NotTo(HaveOccurred())

		err = os.Mkdir(filepath.Join(tmpDir, "project"), 0755)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(tmpDir, "project", "manifest.yml"), []byte{}, 0555)
		Expect(err).NotTo(HaveOccurred())

		err = os.Mkdir(filepath.Join(tmpDir, "another-project"), 0555)
		Expect(err).NotTo(HaveOccurred())

		request = out.Request{
			Source: resource.Source{
				API:           "https://api.run.pivotal.io",
				Username:      "awesome@example.com",
				Password:      "hunter2",
				Organization:  "org",
				Space:         "space",
				SkipCertCheck: true,
			},
			Params: out.Params{
				Create: []string{"foo.example.com"},
			},
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when my manifest and paths do not contain a glob", func() {
		It("pushes an application to cloud foundry", func() {
			session, err := gexec.Start(
				cmd,
				GinkgoWriter,
				GinkgoWriter,
			)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			var response out.Response
			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.Version.Timestamp).To(BeTemporally("~", time.Now(), time.Second))

			// shim outputs arguments
			// Expect(session.Err).To(gbytes.Say("cf api https://api.run.pivotal.io --skip-ssl-validation"))
			// Expect(session.Err).To(gbytes.Say("cf auth awesome@example.com hunter2"))
			// Expect(session.Err).To(gbytes.Say("cf target -o org -s space"))
			// Expect(session.Err).To(gbytes.Say("cf zero-downtime-push awesome-app -f %s",
			// 	filepath.Join(tmpDir, "project/manifest.yml"),
			// ))
			// Expect(session.Err).To(gbytes.Say(filepath.Join(tmpDir, "another-project")))

			// color should be always
			Expect(session.Err).To(gbytes.Say("CF_COLOR=true"))
		})
	})

})
