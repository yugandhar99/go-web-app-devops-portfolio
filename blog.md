# From Basic Go App to DevOps-Ready Portfolio Project

This project started as a simple Go web application serving static pages. I enhanced it into a production-style DevOps project by adding containerization, CI/CD, security scanning, SBOM generation, Kubernetes packaging, and observability endpoints.

## What changed

The application now exposes `/healthz`, `/readyz`, and `/metrics` endpoints. These are important because Kubernetes and monitoring platforms need clear signals to understand whether an application is alive, ready to receive traffic, and producing operational metrics.

The Dockerfile uses a multi-stage build so build tools stay out of the production runtime image. This reduces image size and lowers the attack surface.

GitHub Actions now validates the project through Go tests, linting, Docker build, Helm checks, Trivy scanning, SBOM generation, CodeQL, and dependency review.

## Why it matters

This project demonstrates the full path from application code to production delivery. It is not only about writing a Go app; it shows how to package, test, scan, monitor, and deploy it using practices that are common in DevOps and platform engineering teams.
