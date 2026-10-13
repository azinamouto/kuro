package scope

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// Policy is intentionally loopback-only for the alpha release.
type Policy struct {
	Version     int    `json:"version"`
	Origin      string `json:"origin"`
	MaxRequests int    `json:"max_requests"`
}

func Load(path string) (Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return Policy{}, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	var p Policy
	if err := d.Decode(&p); err != nil {
		return p, err
	}
	if p.Version != 1 {
		return p, errors.New("unsupported policy version")
	}
	if p.MaxRequests < 1 || p.MaxRequests > 1000 {
		return p, errors.New("max_requests must be 1..1000")
	}
	if err := p.Check(p.Origin); err != nil {
		return p, err
	}
	return p, nil
}
func (p Policy) Check(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	base, err := url.Parse(p.Origin)
	if err != nil {
		return err
	}
	if u.Scheme != "http" || u.User != nil || u.Fragment != "" || u.Opaque != "" {
		return errors.New("only plain HTTP URLs without userinfo or fragments are supported")
	}
	if u.Host == "" || u.Hostname() == "" || u.Port() == "" {
		return errors.New("explicit loopback host and port required")
	}
	ip := net.ParseIP(u.Hostname())
	if ip == nil || !ip.IsLoopback() {
		return errors.New("only literal loopback IP targets allowed")
	}
	if u.Host != base.Host || u.Scheme != base.Scheme {
		return fmt.Errorf("out of scope: %s", u.Host)
	}
	if strings.Contains(strings.ToLower(u.EscapedPath()), "%2f") || strings.Contains(strings.ToLower(u.EscapedPath()), "%5c") {
		return errors.New("ambiguous path encoding rejected")
	}
	return nil
}

