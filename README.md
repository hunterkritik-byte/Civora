# Civora

**Real GitHub Actions CI cost and performance optimizer.**

Civora helps engineering teams understand where their GitHub Actions time and runner capacity are being spent, identify CI waste, and apply measurable optimizations.

> **See where your CI time and money go.**

## Why Civora?

CI can become a significant engineering and infrastructure cost as repositories grow. Slow jobs, repeated failures, inefficient caching, unnecessary work, and poorly structured workflows all consume runner time.

Civora focuses on **measuring real CI behavior first** and making recommendations from that evidence. It does not manufacture savings numbers or rely on simulated workflow data.

## Current MVP

- Real GitHub Actions workflow runs
- Real job and step timing metadata
- Job-duration aggregation
- Failure-rate analysis
- Slow-job detection
- Recurring-failure detection
- JSON reports
- No simulated workflow data

## Quick start

Set a GitHub token with the minimum required **Actions: read** permission:

```bash
export GITHUB_TOKEN='your_token'
go run ./cmd/civora -repo OWNER/REPO -runs 30
```

For private repositories, use a fine-grained token with Actions read-only access to the repository being analyzed.

Example:

```bash
go run ./cmd/civora -repo hunterkritik-byte/Civora -runs 30
```

## Product direction

Civora is being built as a real developer platform, not a static CI report generator:

```text
GitHub App
    ↓
Repository authorization
    ↓
Actions events + workflow data
    ↓
Civora analysis engine
    ↓
Cost / performance / reliability findings
    ↓
Dashboard + reports
    ↓
Approved optimization PRs
```

### Planned capabilities

#### GitHub App
- One-click repository installation
- Least-privilege permissions
- Webhook-based workflow ingestion
- Repository and organization support
- Secure installation/token handling

#### CI performance
- Job and step bottleneck analysis
- Workflow critical-path analysis
- Queue/wait-time visibility
- Repeated-work detection
- Flaky/retry waste detection

#### Cache intelligence
- Cache hit/miss analysis
- Dependency-install waste detection
- Cache-key recommendations
- Cache effectiveness trends

#### Cost intelligence
- Runner-type-aware cost calculations
- Configurable pricing
- Usage trends
- Cost allocation by repository/workflow/job
- Explicit estimated-vs-measured labeling

#### Automated optimization
- Workflow YAML recommendations
- Explainable proposed changes
- Optional GitHub PR generation
- Human approval before repository changes

#### Dashboard
- CI cost and usage overview
- Performance trends
- Top waste sources
- Repository/workflow drill-down
- Findings history

## Sponsorship & partners

Civora is intended to become a production-grade open-source developer tool for teams running CI at scale.

We are interested in sponsorships and technical partnerships with:

- CI/CD platforms
- Cloud and runner providers
- Developer infrastructure companies
- GitHub ecosystem companies
- Build and caching infrastructure vendors
- Engineering productivity platforms

Sponsorship can support continued open-source development, GitHub App infrastructure, testing across larger repositories, documentation, and new integrations.

**Interested in sponsoring or partnering?**  
Open an issue in this repository or contact the maintainer listed on the GitHub profile.

## Security & privacy

Civora should follow least-privilege access and minimize repository data collection.

The project is designed to analyze CI metadata rather than copy entire source repositories. Production deployment will include explicit data-retention controls, secret handling, auditability, and secure GitHub App authentication.

Security reports should be handled privately rather than posted publicly when they contain sensitive details.

## Roadmap

- [x] Real GitHub Actions API analyzer
- [x] Job timing and failure analysis
- [x] JSON reporting
- [ ] GitHub App installation flow
- [ ] Webhook ingestion
- [ ] Cache effectiveness analyzer
- [ ] Workflow critical-path analysis
- [ ] Runner-aware cost engine
- [ ] HTML/web dashboard
- [ ] Optimization recommendation engine
- [ ] Optional automated PR generation
- [ ] Organization/team support
- [ ] Production security hardening
- [ ] Sponsor/enterprise deployment documentation

## Contributing

Contributions are welcome, especially around GitHub Actions analysis, CI performance, caching, cost modeling, security, and developer experience.

Before submitting changes, run:

```bash
go test ./...
```

## License

See the repository license for the current project terms.
