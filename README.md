<div align="center">

# KURO

**A hybrid web penetration-testing framework, built in Go.**

*Inspect the surface. Verify the weakness. Document the evidence.*

[![Go](https://img.shields.io/badge/Go-1.x-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
![Status](https://img.shields.io/badge/status-v0.1--alpha-222222?style=flat-square)
![Scope](https://img.shields.io/badge/scope-loopback%20only-555555?style=flat-square)

</div>

---

## Overview

**KURO** is an open-source project exploring the engineering behind **hybrid dynamic application security testing (DAST)**: automated checks, human-led investigation, reproducible evidence, and developer-friendly security reports.

The long-term goal is a modular penetration-testing workbench for **explicitly authorized** web applications and APIs. The current alpha is intentionally much smaller: a **single-request, loopback-only HTTP inspector** with basic response-header observations and JSON output.

> **Alpha notice:** KURO is experimental software. It is **not** a comprehensive vulnerability scanner, exploit framework, or substitute for a professional security assessment.

## Current capabilities

| Capability | v0.1-alpha |
| :--- | :--- |
| Scope policy loaded from JSON | Available |
| Explicit loopback IP and port enforcement | Available |
| Single HTTP GET inspection | Available |
| Out-of-scope redirect prevention | Available (redirects are not followed) |
| Basic HTTP security-header observations | Available |
| JSON report generation | Available |
| Automated scope and scanner tests | Available |
| Endpoint crawling and discovery | Planned |
| Authenticated testing and manual replay | Planned |
| Injection and authorization test modules | Planned |
| HTML/PDF reporting | Planned |

### Checks implemented

- **HDR-001:** `X-Content-Type-Options` is not `nosniff`.
- **HDR-002:** `Content-Security-Policy` is missing.
- **HDR-003:** HTML responses have neither `X-Frame-Options` nor a CSP `frame-ancestors` directive.
- **SCOPE-001:** An HTTP redirect points outside the permitted origin; KURO records that it did not follow it.

These are **observations**, not automatic proof of exploitable vulnerabilities. Context and manual validation matter.

## Quick start

**Requirements:** Go installed and a local HTTP application running on `127.0.0.1:8080` (or another explicitly configured loopback port).

```bash
go test ./...
go build -o kuro ./cmd/kuro
```

On Windows PowerShell, build with:

```powershell
go build -o kuro.exe ./cmd/kuro
```

Example policy (`scope.example.json`):

```json
{
  "version": 1,
  "origin": "http://127.0.0.1:8080",
  "max_requests": 10
}
```

Validate the policy:

```powershell
.\kuro.exe scope -file scope.example.json
```

Scan a running local service:

```powershell
.\kuro.exe scan -scope scope.example.json -target http://127.0.0.1:8080 -out report.json
```

Inspect the report:

```powershell
Get-Content .\report.json
```

**No server listening?** For a disposable local test target, open a second terminal and run:

```powershell
python -m http.server 8080 --bind 127.0.0.1
```

Then rerun KURO in the first terminal. A missing-header observation on this simple server does not necessarily indicate an exploitable issue.

## CLI reference

```text
kuro scope -file <policy.json>
kuro scan  -scope <policy.json> [-target <url>] [-out <report.json>]
```

If `-target` is omitted, KURO scans the `origin` defined in the policy. The current alpha makes one GET request and does **not** crawl links or follow redirects.

## Architecture

```text
kuro/
├── cmd/kuro/          # CLI entry point
├── internal/
│   ├── scope/         # Loopback-only target policy
│   ├── scanner/       # HTTP request and header checks
│   └── report/        # JSON report writer
├── docs/              # Engineering roadmap
├── testdata/          # Test fixtures
├── scope.example.json
└── go.mod
```

**Design direction:** separate scope enforcement, HTTP transport, discovery, test modules, evidence validation, and reporting. Every future network-capable module should use the same scope-enforcing transport rather than creating its own unrestricted client.

## Roadmap

| Release | Focus |
| :--- | :--- |
| **v0.1-alpha** | Loopback policy, single GET, header observations, JSON report |
| **v0.1-beta** | Stronger scope tests, bounded crawler, request budget, evidence fingerprints |
| **v0.2** | Manual request replay, synthetic identities, authorization assertions |
| **v0.3** | Controlled lab-only injection, XSS, and path-handling validation |
| **v0.4** | HTML reports, retesting, redacted evidence manifests |
| **v0.5** | CI integration, benchmarks, coverage and recovery |
| **v1.0** | Stable interfaces, threat model, reproducible builds, independent scope review |

See [`docs/ROADMAP.md`](docs/ROADMAP.md) for the repository roadmap.

## Security and responsible use

KURO is intended **only for systems you own or have explicit permission to assess**.

The current alpha deliberately restricts connections to **literal loopback IP addresses**, disables proxy use, and does not follow redirects. These controls are useful early guardrails, **not a fully audited authorization boundary**. Do not expand the scope to public or production targets without further review and testing.

- Prefer isolated labs, synthetic data, and disposable services.
- Do not commit private scan results, tokens, cookies, or credentials.
- Verify findings manually before treating them as vulnerabilities.
- Treat application responses as untrusted input.

## Development

```bash
go test ./...
go test -race ./...
go vet ./...
```

KURO is being developed as an **engineering-first, AI-independent** project: networking, parsing, test execution, evidence handling, and reporting are implemented in Go without an AI runtime dependency.

Contributions should include a clear security rationale, deterministic tests, and a description of possible side effects. For now, keep changes focused on the alpha architecture and authorized lab testing.

---

<div align="center">

**KURO** — *Findings need evidence. Evidence needs context.*

Built for learning, testing, and responsible security engineering.

</div>
