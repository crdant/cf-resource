package out_test

import (
	"encoding/json"
	"os/exec"
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
		cmd     *exec.Cmd
		request out.Request
	)

	BeforeEach(func() {
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

	Context("when a route to create is provided to the resource", func() {
		It("creates a route in cloud foundry", func() {
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
			Expect(session.Err).To(gbytes.Say("cf api https://api.run.pivotal.io --skip-ssl-validation"))
			Expect(session.Err).To(gbytes.Say("cf auth awesome@example.com hunter2"))
			Expect(session.Err).To(gbytes.Say("cf target -o org -s space"))
			Expect(session.Err).To(gbytes.Say("cf create-route volcano-base example.com --hostname foo"))

			// color should be always
			Expect(session.Err).To(gbytes.Say("CF_COLOR=true"))
		})
	})

})
