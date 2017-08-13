package out

import (
	"net/url"
	"strconv"
	"strings"
)

type Route struct {
	domain     string
	host       string
	path       string
	port       int32
	randomPort bool
}

func NewRoute(domain string, host string, path string, port int32, randomPort bool) *Route {
	r := new(Route)
	r.domain = domain
	r.host = host
	r.path = path
	r.randomPort = randomPort
	return r
}

func ParseRoute(route string, randomPort bool) (*Route, error) {
	r := new(Route)
	url, err := url.Parse("https://" + route)
	if err != nil {
		return nil, err
	}

	p, err := strconv.ParseInt(url.Port(), 0, 32)
	if err != nil {
		r.port = int32(p)
	}

	r.path = url.RequestURI()
	fqdn := strings.SplitN(url.Hostname(), ".", 2)
	r.host = fqdn[0]
	r.domain = fqdn[1]
	return r, nil
}

func (route *Route) String() string {
	s := route.host + "." + route.domain

	if len(route.path) > 0 {
		s = s + "/" + route.path
	}
	if route.port > 0 {
		s = s + ":" + strconv.Itoa(int(route.port))
	}

	return s
}
