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
	r.port = port
	r.randomPort = randomPort
	return r
}

func ParseRoute(route string, randomPort bool) (*Route, error) {
	r := new(Route)
	url, err := url.Parse("https://" + route)
	if err != nil {
		return nil, err
	}

	if len(url.Port()) > 0 {
		p, err := strconv.Atoi(url.Port())
		if err != nil {
			return nil, err
		}
		r.port = int32(p)
	}

	uri := url.RequestURI()
	if uri != "/" {
		r.path = strings.Trim(uri, "/")
	}

	if r.port > 0 || randomPort {
		r.domain = url.Hostname()
	} else {
		fqdn := strings.SplitN(url.Hostname(), ".", 2)
		r.host = fqdn[0]
		r.domain = fqdn[1]
	}

	r.randomPort = randomPort
	return r, nil
}

func (route *Route) String() string {
	var s string

	if len(route.host) > 0 {
		s = route.host + "." + route.domain
	} else {
		s = route.domain
	}

	if route.port > 0 {
		s = s + ":" + strconv.Itoa(int(route.port))
	}

	if len(route.path) > 0 {
		s = s + "/" + route.path
	}

	return s
}
