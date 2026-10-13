package scanner

import (
	"context"
	"fmt"
	"github.com/azinamouto/kuro/internal/scope"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Finding struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Severity       string `json:"severity"`
	Confidence     string `json:"confidence"`
	Evidence       string `json:"evidence"`
	Recommendation string `json:"recommendation"`
}
type Result struct {
	Target    string    `json:"target"`
	Status    int       `json:"status"`
	Findings  []Finding `json:"findings"`
	CheckedAt time.Time `json:"checked_at"`
}

func Scan(ctx context.Context, p scope.Policy, target string) (Result, error) {
	if err := p.Check(target); err != nil {
		return Result{}, err
	}
	count := 0
	transport := &http.Transport{Proxy: nil, DisableKeepAlives: true, DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return nil, fmt.Errorf("blocked non-loopback dial: %s", address)
		}
		d := net.Dialer{Timeout: 3 * time.Second}
		return d.DialContext(ctx, network, address)
	}}
	defer transport.CloseIdleConnections()
	client := &http.Client{Timeout: 5 * time.Second, Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	if count >= p.MaxRequests {
		return Result{}, fmt.Errorf("request budget exceeded")
	}
	count++
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("User-Agent", "KURO-alpha/0.1 (authorized-local-test)")
	resp, err := client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	r := Result{Target: target, Status: resp.StatusCode, CheckedAt: time.Now().UTC(), Findings: []Finding{}}
	add := func(id, title, severity, evidence, fix string) {
		r.Findings = append(r.Findings, Finding{id, title, severity, "observed", evidence, fix})
	}
	if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
		add("HDR-001", "Missing X-Content-Type-Options", "low", "Header is not set to nosniff", "Set X-Content-Type-Options: nosniff")
	}
	if resp.Header.Get("Content-Security-Policy") == "" {
		add("HDR-002", "Missing Content-Security-Policy", "informational", "No CSP header observed", "Consider a restrictive CSP suitable for this application")
	}
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "text/html") && resp.Header.Get("X-Frame-Options") == "" && !strings.Contains(strings.ToLower(resp.Header.Get("Content-Security-Policy")), "frame-ancestors") {
		add("HDR-003", "Missing framing protection", "low", "No X-Frame-Options or CSP frame-ancestors observed", "Set CSP frame-ancestors or X-Frame-Options")
	}
	if loc := resp.Header.Get("Location"); resp.StatusCode >= 300 && resp.StatusCode < 400 && loc != "" {
		u, err := url.Parse(loc)
		if err == nil {
			next := req.URL.ResolveReference(u)
			if p.Check(next.String()) != nil {
				add("SCOPE-001", "Redirect outside scope blocked", "informational", "Redirect was not followed", "Review the redirect destination separately before authorizing")
			}
		}
	}
	return r, nil
}

