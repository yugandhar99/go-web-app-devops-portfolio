# GenAI Enhancement

## Purpose

The GenAI enhancement is intentionally lightweight and practical. The script generates an AI-style release summary from local commit history and scan artifacts.

## Why it matters

Modern DevOps teams are using AI-assisted workflows for:

- Release notes
- Change-risk summaries
- Incident summaries
- Vulnerability triage summaries
- Deployment readiness reviews

## Local Mode

```bash
python3 scripts/genai_release_summary.py --mode offline --output release-summary.md
```

## Future Bedrock Mode

In a real AWS environment, this can be connected to Amazon Bedrock using short-lived AWS credentials from GitHub Actions OIDC. Do not store long-lived AWS access keys in GitHub secrets.

## Example Output

```text
Risk level: Low to Medium
Review focus: endpoint changes, container image scan findings, Helm probe settings, and public ingress configuration before production promotion.
```
