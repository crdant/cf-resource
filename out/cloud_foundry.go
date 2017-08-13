package out

import (
	"os"
	"os/exec"
	"strconv"
)

type PAAS interface {
	Login(api string, username string, password string, insecure bool) error
	Target(organization string, space string) error
	CreateRoute(space string, domain string, host string, path string, port int32, randomPort bool) error
	MapRoute(app string, domain string, host string, path string, port int32) error
	UnmapRoute(app string, domain string, host string, path string, port int32) error
}

type CloudFoundry struct{}

func NewCloudFoundry() *CloudFoundry {
	return &CloudFoundry{}
}

func (cf *CloudFoundry) Login(api string, username string, password string, insecure bool) error {
	args := []string{"api", api}
	if insecure {
		args = append(args, "--skip-ssl-validation")
	}

	err := cf.cf(args...).Run()
	if err != nil {
		return err
	}

	return cf.cf("auth", username, password).Run()
}

func (cf *CloudFoundry) Target(organization string, space string) error {
	return cf.cf("target", "-o", organization, "-s", space).Run()
}

func (cf *CloudFoundry) CreateRoute(space string, domain string, host string, path string, port int32, randomPort bool) error {
	args := []string{"create-route", space, domain}
	if len(host) > 0 {
		args = append(args, "--hostname", host)

	}
	if len(path) > 0 {
		args = append(args, "--path", path)
	}
	if randomPort {
		args = append(args, "--random-port")
	} else if port > 0 {
		args = append(args, "--port", strconv.Itoa(int(port)))
	}

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) MapRoute(app string, domain string, host string, path string, port int32) error {
	args := []string{"map-route", app, domain}
	if len(host) > 0 {
		args = append(args, "--hostname", host)

	}
	if len(path) > 0 {
		args = append(args, "--path", path)
	}
	if port > 0 {
		args = append(args, "--port", strconv.Itoa(int(port)))
	}

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) UnmapRoute(app string, domain string, host string, path string, port int32) error {
	args := []string{"unmap-route", app, domain}
	if len(host) > 0 {
		args = append(args, "--hostname", host)

	}
	if len(path) > 0 {
		args = append(args, "--path", path)
	}
	if port > 0 {
		args = append(args, "--port", strconv.Itoa(int(port)))
	}

	return cf.cf(args...).Run()
}

func (cf *CloudFoundry) cf(args ...string) *exec.Cmd {
	cmd := exec.Command("cf", args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CF_COLOR=true")

	return cmd
}
