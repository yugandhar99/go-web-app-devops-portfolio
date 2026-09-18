#!/usr/bin/env python3
"""Generate an AI-style release summary for a Go DevOps portfolio project.

The script is safe by default: offline mode creates a deterministic summary from
local Git metadata and optional scan files. Bedrock mode is intentionally a
placeholder pattern so no API keys are committed to the repository.
"""

from __future__ import annotations

import argparse
import json
import os
import subprocess
from pathlib import Path
from typing import Iterable


def run(command: list[str]) -> str:
    try:
        return subprocess.check_output(command, text=True, stderr=subprocess.DEVNULL).strip()
    except Exception:
        return ""


def recent_commits(limit: int = 8) -> list[str]:
    output = run(["git", "log", f"-{limit}", "--pretty=format:%h %s"])
    return [line for line in output.splitlines() if line.strip()] or ["No Git commit history available in this environment."]


def summarize_scan_files(paths: Iterable[Path]) -> list[str]:
    notes: list[str] = []
    for path in paths:
        if not path.exists():
            continue
        try:
            if path.suffix == ".json":
                data = json.loads(path.read_text(encoding="utf-8"))
                if isinstance(data, dict):
                    notes.append(f"Parsed {path.name}: top-level keys={', '.join(list(data.keys())[:8])}")
                else:
                    notes.append(f"Parsed {path.name}: JSON list with {len(data)} entries")
            else:
                content = path.read_text(encoding="utf-8", errors="ignore")
                notes.append(f"Read {path.name}: {len(content.splitlines())} lines")
        except Exception as exc:  # pragma: no cover - defensive CLI behavior
            notes.append(f"Could not parse {path.name}: {exc}")
    return notes or ["No scan artifacts found. CI should attach Trivy/SBOM artifacts when available."]


def offline_summary() -> str:
    commits = recent_commits()
    scans = summarize_scan_files([Path("trivy-results.json"), Path("sbom.cdx.json"), Path("coverage.out")])
    env = os.getenv("APP_ENV", "ci")

    lines = [
        "# AI-Assisted Release Summary",
        "",
        f"Environment: `{env}`",
        "",
        "## Recent Changes",
    ]
    lines.extend(f"- {commit}" for commit in commits)
    lines.extend([
        "",
        "## Validation Signals",
    ])
    lines.extend(f"- {scan}" for scan in scans)
    lines.extend([
        "",
        "## Release Risk Summary",
        "- Risk level: Low to Medium for a static Go web application when tests, linting, Docker build, Helm validation, and vulnerability scans pass.",
        "- Review focus: endpoint changes, container image scan findings, Helm probe settings, and public ingress configuration before production promotion.",
        "",
        "## Recommended Next Action",
        "- Attach CI artifacts, review any HIGH/CRITICAL vulnerabilities, and promote the immutable container tag through the Helm values file or GitOps repo.",
    ])
    return "\n".join(lines) + "\n"


def bedrock_summary() -> str:
    # Keep this repository credential-free. In a real workflow, call Amazon Bedrock
    # with short-lived AWS credentials from GitHub Actions OIDC and pass the prompt
    # created by offline_summary().
    return offline_summary() + "\n> Bedrock mode placeholder: wire this to Amazon Bedrock in a private CI environment using OIDC-based AWS auth.\n"


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate release summary from local CI context.")
    parser.add_argument("--mode", choices=["offline", "bedrock"], default="offline")
    parser.add_argument("--output", default="release-summary.md")
    args = parser.parse_args()

    summary = bedrock_summary() if args.mode == "bedrock" else offline_summary()
    Path(args.output).write_text(summary, encoding="utf-8")
    print(f"Wrote {args.output}")


if __name__ == "__main__":
    main()
