# Architecture

## Application Architecture

```text
Browser
  |
  v
Go net/http server
  |
  +-- Static HTML pages
  +-- Health endpoint
  +-- Readiness endpoint
  +-- Prometheus-style metrics endpoint
```

## DevOps Architecture

```text
Developer Push
  |
  v
GitHub Actions
  |-- gofmt / go vet / tests / coverage
  |-- lint
  |-- Helm lint/template
  |-- Docker build
  |-- Trivy scan
  |-- SBOM generation
  |-- CodeQL analysis
  v
Container Registry / Kubernetes Deployment
```

## Kubernetes Readiness

The Helm chart includes:

- Deployment
- Service
- Optional Ingress
- ServiceAccount
- HPA
- Liveness probe
- Readiness probe
- Pod security context
- Container security context
- Prometheus scrape annotations
