package out_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crdant/cf-route-resource"
	"github.com/crdant/cf-route-resource/out"
	"github.com/crdant/cf-route-resource/out/fakes"
)

var _ = Describe("Out Command", func() {
	var (
		cloudFoundry *fakes.FakePAAS
		request      out.Request
		command      *out.Command
	)

	BeforeEach(func() {
		cloudFoundry = &fakes.FakePAAS{}
		command = out.NewCommand(cloudFoundry)

		request = out.Request{
			Source: resource.Source{
				API:          "https://api.run.pivotal.io",
				Username:     "awesome@example.com",
				Password:     "hunter2",
				Organization: "secret",
				Space:        "volcano-base",
			},
			Params: out.Params{
				Create: []string{"foo.example.com"},
			},
		}
	})

	Describe("running the command", func() {
		It("creates a new route in cloud foundry", func() {
			response, err := command.Run(request)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.Version.Timestamp).To(BeTemporally("~", time.Now(), time.Second))
			Expect(response.Metadata[0]).To(Equal(
				resource.MetadataPair{
					Name:  "organization",
					Value: "secret",
				},
			))
			Expect(response.Metadata[1]).To(Equal(
				resource.MetadataPair{
					Name:  "space",
					Value: "volcano-base",
				},
			))

			By("logging in")
			Expect(cloudFoundry.LoginCallCount()).To(Equal(1))

			api, username, password, insecure := cloudFoundry.LoginArgsForCall(0)
			Expect(api).To(Equal("https://api.run.pivotal.io"))
			Expect(username).To(Equal("awesome@example.com"))
			Expect(password).To(Equal("hunter2"))
			Expect(insecure).To(Equal(false))

			By("targetting the organization and space")
			Expect(cloudFoundry.TargetCallCount()).To(Equal(1))

			org, space := cloudFoundry.TargetArgsForCall(0)
			Expect(org).To(Equal("secret"))
			Expect(space).To(Equal("volcano-base"))

			By("creating the route")
			Expect(cloudFoundry.CreateRouteCallCount()).To(Equal(1))

			routes := cloudFoundry.CreateRouteArgsForCall(0)
			Expect(routes).To(Equal([]string{"foo.example.com"}))
		})

		Describe("handling any errors", func() {
			var expectedError error

			BeforeEach(func() {
				expectedError = errors.New("it all went wrong")
			})

			It("from logging in", func() {
				cloudFoundry.LoginReturns(expectedError)

				_, err := command.Run(request)
				Expect(err).To(MatchError(expectedError))
			})

			It("from targetting an org and space", func() {
				cloudFoundry.TargetReturns(expectedError)

				_, err := command.Run(request)
				Expect(err).To(MatchError(expectedError))
			})

			It("from creating the route", func() {
				cloudFoundry.CreateRouteReturns(expectedError)

				_, err := command.Run(request)
				Expect(err).To(MatchError(expectedError))
			})
		})

		It("lets people skip the certificate check", func() {
			request.Source.SkipCertCheck = false

			_, err := command.Run(request)
			Expect(err).NotTo(HaveOccurred())

			By("logging in")
			Expect(cloudFoundry.LoginCallCount()).To(Equal(1))

			_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
			Expect(insecure).To(Equal(true))
		})

		// It("lets people create multiple routes", func() {
		// 	request = out.Request{
		// 		Source: resource.Source{
		// 			API:           "https://api.run.pivotal.io",
		// 			Username:      "awesome@example.com",
		// 			Password:      "hunter2",
		// 			Organization:  "secret",
		// 			Space:         "volcano-base",
		// 			SkipCertCheck: true,
		// 		},
		// 		Params: out.Params{
		// 			Create: []string{"foo.example.com"},
		// 		},
		// 	}
		//
		// 	_, err := command.Run(request)
		// 	Expect(err).NotTo(HaveOccurred())
		//
		// 	By("logging in")
		// 	Expect(cloudFoundry.LoginCallCount()).To(Equal(1))
		//
		// 	_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
		// 	Expect(insecure).To(Equal(true))
		// })

		It("creates a new path route in cloud foundry", func() {
			request.Params.Create = []string{"foo.example.com/bar"}

			response, err := command.Run(request)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.Version.Timestamp).To(BeTemporally("~", time.Now(), time.Second))
			Expect(response.Metadata[0]).To(Equal(
				resource.MetadataPair{
					Name:  "organization",
					Value: "secret",
				},
			))
			Expect(response.Metadata[1]).To(Equal(
				resource.MetadataPair{
					Name:  "space",
					Value: "volcano-base",
				},
			))

			By("logging in")
			Expect(cloudFoundry.LoginCallCount()).To(Equal(1))

			api, username, password, insecure := cloudFoundry.LoginArgsForCall(0)
			Expect(api).To(Equal("https://api.run.pivotal.io"))
			Expect(username).To(Equal("awesome@example.com"))
			Expect(password).To(Equal("hunter2"))
			Expect(insecure).To(Equal(false))

			By("targetting the organization and space")
			Expect(cloudFoundry.TargetCallCount()).To(Equal(1))

			org, space := cloudFoundry.TargetArgsForCall(0)
			Expect(org).To(Equal("secret"))
			Expect(space).To(Equal("volcano-base"))

			By("creating the route")
			Expect(cloudFoundry.CreateRouteCallCount()).To(Equal(1))

			routes := cloudFoundry.CreateRouteArgsForCall(0)
			Expect(routes).To(Equal([]string{"foo.example.com/bar"}))
		})

		It("creates a new tcp route in cloud foundry", func() {
			request.Params.Create = []string{"foo.example.com:1202"}

			response, err := command.Run(request)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.Version.Timestamp).To(BeTemporally("~", time.Now(), time.Second))
			Expect(response.Metadata[0]).To(Equal(
				resource.MetadataPair{
					Name:  "organization",
					Value: "secret",
				},
			))
			Expect(response.Metadata[1]).To(Equal(
				resource.MetadataPair{
					Name:  "space",
					Value: "volcano-base",
				},
			))

			By("logging in")
			Expect(cloudFoundry.LoginCallCount()).To(Equal(1))

			api, username, password, insecure := cloudFoundry.LoginArgsForCall(0)
			Expect(api).To(Equal("https://api.run.pivotal.io"))
			Expect(username).To(Equal("awesome@example.com"))
			Expect(password).To(Equal("hunter2"))
			Expect(insecure).To(Equal(false))

			By("targetting the organization and space")
			Expect(cloudFoundry.TargetCallCount()).To(Equal(1))

			org, space := cloudFoundry.TargetArgsForCall(0)
			Expect(org).To(Equal("secret"))
			Expect(space).To(Equal("volcano-base"))

			By("creating the route")
			Expect(cloudFoundry.CreateRouteCallCount()).To(Equal(1))

			routes := cloudFoundry.CreateRouteArgsForCall(0)
			Expect(routes).To(Equal([]string{"foo.example.com:1202"}))
		})

		It("lets people use a random port", func() {
			response, err := command.Run(request)
			Expect(err).NotTo(HaveOccurred())

			Expect(response.Version.Timestamp).To(BeTemporally("~", time.Now(), time.Second))
			Expect(response.Metadata[0]).To(Equal(
				resource.MetadataPair{
					Name:  "organization",
					Value: "secret",
				},
			))
			Expect(response.Metadata[1]).To(Equal(
				resource.MetadataPair{
					Name:  "space",
					Value: "volcano-base",
				},
			))

			By("logging in")
			Expect(cloudFoundry.LoginCallCount()).To(Equal(1))

			api, username, password, insecure := cloudFoundry.LoginArgsForCall(0)
			Expect(api).To(Equal("https://api.run.pivotal.io"))
			Expect(username).To(Equal("awesome@example.com"))
			Expect(password).To(Equal("hunter2"))
			Expect(insecure).To(Equal(false))

			By("targetting the organization and space")
			Expect(cloudFoundry.TargetCallCount()).To(Equal(1))

			org, space := cloudFoundry.TargetArgsForCall(0)
			Expect(org).To(Equal("secret"))
			Expect(space).To(Equal("volcano-base"))

			By("creating the route")
			Expect(cloudFoundry.CreateRouteCallCount()).To(Equal(1))

			routes := cloudFoundry.CreateRouteArgsForCall(0)
			Expect(routes).To(Equal("foo.example.com"))
		})

		It("maps a route in cloud foundry", func() {
			request = out.Request{
				Source: resource.Source{
					API:           "https://api.run.pivotal.io",
					Username:      "awesome@example.com",
					Password:      "hunter2",
					Organization:  "secret",
					Space:         "volcano-base",
					SkipCertCheck: true,
				},
				Params: out.Params{
					Create: []string{"foo.example.com"},
				},
			}

			_, err := command.Run(request)
			Expect(err).NotTo(HaveOccurred())

			By("logging in")
			Expect(cloudFoundry.LoginCallCount()).To(Equal(1))

			_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
			Expect(insecure).To(Equal(true))
		})

		// It("lets people map multiple routes", func() {
		// 	request = out.Request{
		// 		Source: resource.Source{
		// 			API:           "https://api.run.pivotal.io",
		// 			Username:      "awesome@example.com",
		// 			Password:      "hunter2",
		// 			Organization:  "secret",
		// 			Space:         "volcano-base",
		// 			SkipCertCheck: true,
		// 		},
		// 		Params: out.Params{
		// 			Create: []string{"foo.example.com"},
		// 		},
		// 	}
		//
		// 	_, err := command.Run(request)
		// 	Expect(err).NotTo(HaveOccurred())
		//
		// 	By("logging in")
		// 	Expect(cloudFoundry.LoginCallCount()).To(Equal(1))
		//
		// 	_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
		// 	Expect(insecure).To(Equal(true))
		// })

		It("unmaps a route in cloud foundry", func() {
			request = out.Request{
				Source: resource.Source{
					API:           "https://api.run.pivotal.io",
					Username:      "awesome@example.com",
					Password:      "hunter2",
					Organization:  "secret",
					Space:         "volcano-base",
					SkipCertCheck: true,
				},
				Params: out.Params{
					Create: []string{"foo.example.com"},
				},
			}

			_, err := command.Run(request)
			Expect(err).NotTo(HaveOccurred())

			By("logging in")
			Expect(cloudFoundry.LoginCallCount()).To(Equal(1))

			_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
			Expect(insecure).To(Equal(true))
		})

		// It("lets people unmap multiple routes", func() {
		// 	request = out.Request{
		// 		Source: resource.Source{
		// 			API:           "https://api.run.pivotal.io",
		// 			Username:      "awesome@example.com",
		// 			Password:      "hunter2",
		// 			Organization:  "secret",
		// 			Space:         "volcano-base",
		// 			SkipCertCheck: true,
		// 		},
		// 		Params: out.Params{
		// 			Create: []string{"foo.example.com"},
		// 		},
		// 	}
		//
		// 	_, err := command.Run(request)
		// 	Expect(err).NotTo(HaveOccurred())
		//
		// 	By("logging in")
		// 	Expect(cloudFoundry.LoginCallCount()).To(Equal(1))
		//
		// 	_, _, _, insecure := cloudFoundry.LoginArgsForCall(0)
		// 	Expect(insecure).To(Equal(true))
		// })

	})
})
