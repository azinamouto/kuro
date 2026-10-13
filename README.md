# KURO

**Hybrid Web Penetration Testing Framework — v0.1-alpha**

A Go CLI for authorized, reproducible web application security testing. This first milestone implements a local-only HTTP inspection engine; it is **not** a full penetration testing suite.

## Requirements

- Go 1.23+
- A web app you own running on `127.0.0.1` (or `::1`) on an explicit port

## Quick start

```bash
go test ./...
go run ./cmd/kuro scope -file scope.example.json
go run ./cmd/kuro scan -scope scope.example.json -out report.json
```

Edit `scope.example.json` to match the local app's exact IP and port. You can use `go build -o kuro ./cmd/kuro` to produce a binary (on Windows use `-o kuro.exe`).

## What works

- Strict, literal loopback-IP and explicit-port scope
- No HTTP proxy inheritance; no automatic redirect following
- Connection destination restricted to loopback
- Single bounded HTTP GET, with a timeout and response-body cap
- Observations for `X-Content-Type-Options`, CSP, and framing protections
- JSON report with restricted file permissions
- Automated scope and HTTP tests

## What does not work yet

- No endpoint crawler, authenticated tests, replay, active payloads, or injection checks
- No arbitrary staging or production targets
- No claim of vulnerability confirmation: header observations are informational or low-confidence signals
- No robust evidence redaction (only bounded response data is discarded, not persisted)
- No comprehensive URL normalization, path policy, or comprehensive SSRF hardening for non-loopback hosts

## Engineering principles

1. Scope validation is enforced at the transport boundary.
2. Every security finding must have evidence and a stated confidence level.
3. No destructive tests by default.
4. Tests run against synthetic fixtures and authorized local services.
5. AI is not a runtime dependency.

## Next steps

1. Add negative tests for URL normalization and scope parsing.
2. Add a shared scheduler with rate limits and cancellation.
3. Introduce a versioned module interface and evidence model.
4. Implement endpoint discovery against local fixtures.
5. Add explicit authorized authenticated testing with synthetic users.

See `docs/ROADMAP.md`.
