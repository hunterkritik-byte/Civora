# Civora

**Real GitHub Actions CI cost and performance optimizer.**

**Sponsorship & partnerships:** hunterkritik@gmail.com

Civora helps engineering teams understand where their GitHub Actions time and runner capacity are being spent, identify CI waste, and apply measurable optimizations.

> **See where your CI time and money go.**

## Why Civora?

CI can become a significant engineering and infrastructure cost as repositories grow. Slow jobs, repeated failures, inefficient caching, unnecessary work, and poorly structured workflows all consume runner time.

Civora focuses on **measuring real CI behavior first** and making recommendations from that evidence. It does not manufacture savings numbers or rely on simulated workflow data.

## Current MVP

- Real GitHub Actions workflow runs
- Real job and step timing metadata
- Job-duration aggregation
- Step-duration aggregation
- Failure-rate analysis
- Slow-job detection
- Recurring-failure detection
- Workflow critical-path approximation
- Configurable runner cost estimation
- Cache/setup signal detection without fabricated hit-rate claims
- JSON reports
- No simulated workflow data

## Quick start

Set a GitHub token with the minimum required **Actions: read** permission:

```bash
export GITHUB_TOKEN='your_token'
go run ./cmd/civora -repo OWNER/REPO -runs 30
```

For private repositories, use a fine-grained token with Actions read-only access to the repository being analyzed.

### Real-world example

Civora was run against the public TensorFlow repository:

```text
CIVORA
────────────────────
CI Waste Score       90/100

Runner time          375.0 min
Failure rate         13.3%

Top waste
  Build and test     19854s
  Initialize containers 926s

Potential improvements
  • Investigate slow job: build-windows-x86 / build-and-test
  • Investigate recurring failures: build-windows-x86 / build-and-test
  • Investigate slow job: build-linux-x86-cuda13 / build-and-test

Monthly projection
  Runs/month: 300
  Current runtime: ~62.5h
  Estimated cost: $30.00
  Potential saving: $6.00 (20% scenario)
```

The workflow, job, and step timings and failures above are observed from GitHub Actions. The Waste Score and recommendations are Civora analysis. Monthly runtime, cost, and potential savings are projections based on the supplied `-monthly-runs`, `-cost-per-minute`, and savings scenario. They are not GitHub billing statements.

Example:

```bash
go run ./cmd/civora -repo hunterkritik-byte/Civora -runs 30

To estimate cost, provide the runner rate you actually use (for example, your internal blended rate):

```bash
go run ./cmd/civora -repo OWNER/REPO -runs 30 -cost-per-minute 0.008
```

Civora labels this as an estimate; it does not claim a universal GitHub Actions price.
```

## What Civora should improve next

The next engineering priorities are:

1. **Better evidence** — collect workflow/job queue time, runner labels, retries, and more historical runs.
2. **Better cost modeling** — distinguish measured GitHub billing data from configurable estimates and support runner-type pricing.
3. **Better cache intelligence** — inspect cache restore/save behavior and dependency-install patterns instead of treating setup signals as cache hits.
4. **Better reliability analysis** — detect recurring failures, flaky tests, retries, and failure clusters.
5. **Better recommendations** — every recommendation should include the observed evidence and avoid unsupported savings claims.
6. **GitHub App ingestion** — move from on-demand API analysis toward secure event-driven analysis.
7. **Dashboard and history** — track whether an optimization actually improved runtime, reliability, or cost.

## Usage

Analyze a repository with real GitHub Actions data:

```bash
export GITHUB_TOKEN='your_token'
go run ./cmd/civora -repo OWNER/REPO -runs 30
```

Generate a cost and monthly projection using your own runner-rate assumption:

```bash
go run ./cmd/civora -repo OWNER/REPO -runs 30 -monthly-runs 300 -cost-per-minute 0.008
```

For machine-readable output:

```bash
go run ./cmd/civora -repo OWNER/REPO -runs 30 -json
```

**Important:** Civora does not claim that the configurable cost rate is your actual GitHub invoice. Use your actual blended/internal runner rate when modeling costs.

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
Contact **hunterkritik@gmail.com**.

## Security & privacy

Civora should follow least-privilege access and minimize repository data collection.

The project is designed to analyze CI metadata rather than copy entire source repositories. Production deployment will include explicit data-retention controls, secret handling, auditability, and secure GitHub App authentication.

Security reports should be handled privately rather than posted publicly when they contain sensitive details.

## Roadmap

- [x] Real GitHub Actions API analyzer
- [x] Job timing and failure analysis
- [x] JSON reporting
- [x] Step and critical-path analysis
- [x] Configurable cost estimation
- [ ] GitHub App installation flow
- [ ] Webhook ingestion
- [ ] Cache effectiveness analyzer
- [x] Workflow critical-path analysis
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

