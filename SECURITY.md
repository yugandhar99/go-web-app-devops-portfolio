# Security Policy

## Supported Use

This repository is a portfolio/demo project. It should be tested and reviewed before being used in a production environment.

## Security Practices Included

- No hardcoded secrets or API keys
- `.env` files are ignored by default
- Docker production image runs as non-root user
- Minimal production image using Docker multi-stage build
- HTTP security headers in application middleware
- GitHub Actions uses least-privilege permissions
- CodeQL static analysis workflow
- Trivy filesystem and container image scanning
- CycloneDX SBOM generation
- Helm pod security context with dropped Linux capabilities

## Reporting Issues

Open a GitHub issue with:

- Description of the issue
- Steps to reproduce
- Impact
- Suggested fix if available

Do not include real credentials, tokens, private keys, or cloud account information in issues.
