# Civora

**Real GitHub Actions CI cost and performance optimizer.**

Civora connects to the real GitHub Actions API, analyzes workflow runs and jobs, and reports measured CI usage and actionable findings.

## Current MVP

- Real workflow runs from GitHub Actions
- Real job and step timing metadata
- Job duration aggregation
- Failure-rate analysis
- Slow-job detection
- Recurring-failure detection
- JSON output
- No simulated workflow data

## Run locally

Set a GitHub token with **Actions: read** access:

    export GITHUB_TOKEN=your_token
    go run ./cmd/civora -repo OWNER/REPO -runs 30

For private repositories, use a fine-grained token with Actions read-only permission.

Civora reports measured CI time. It does **not** invent dollar savings. Cost estimation will use explicit runner-pricing configuration because billing depends on runner type and account context.

## Roadmap

- GitHub App installation flow
- Cache effectiveness analysis
- Workflow YAML optimization suggestions
- Configurable runner cost estimation
- HTML dashboard
- Optional PR generation for approved fixes
