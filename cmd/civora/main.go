package main

import (
 "encoding/json"
 "flag"
 "fmt"
 "os"
 "strings"
 "time"
 "github.com/hunterkritik-byte/Civora/internal/github"
)

func main() {
 repo := flag.String("repo", "", "GitHub repository (owner/name)")
 runs := flag.Int("runs", 30, "workflow runs to inspect")
 costRate := flag.Float64("cost-per-minute", 0, "estimated runner cost in USD per minute; 0 disables cost estimate")
 monthlyRuns := flag.Float64("monthly-runs", 30, "estimated workflow runs per month")
 savingsPct := flag.Float64("savings-scenario", 20, "scenario percentage used for potential savings")
 flag.Parse()
 if *repo == "" { fmt.Fprintln(os.Stderr, "usage: civora -repo owner/name [-runs 30]"); os.Exit(2) }
 token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
 if token == "" { fmt.Fprintln(os.Stderr, "GITHUB_TOKEN is required"); os.Exit(2) }
 report, err := github.New(token).AnalyzeRepository(*repo, *runs)
 if err != nil { fmt.Fprintln(os.Stderr, "civora:", err); os.Exit(1) }
 github.ApplyCost(&report, *costRate)\n github.BuildBusinessMetrics(&report, *monthlyRuns, *costRate, *savingsPct)
 b, _ := json.MarshalIndent(report, "", "  ")
 fmt.Println(string(b))
 fmt.Printf("\nAnalyzed %d workflow runs at %s\n", report.RunsAnalyzed, time.Now().UTC().Format(time.RFC3339))
 fmt.Println("\nCIVORA")\n fmt.Println("────────────────────")\n fmt.Printf("CI Waste Score       %d/100\n", report.WasteScore)\n fmt.Printf("\nRunner time          %.1f min\n", report.EstimatedRunnerMinutes)\n fmt.Printf("Failure rate         %.1f%%\n", report.FailureRatePercent)\n fmt.Println("\nTop waste")\n for i, s := range report.Steps { if i >= 2 { break }; fmt.Printf("  %-18s %.0fs\n", s.Name, s.TotalSeconds) }\n fmt.Println("\nPotential improvements")\n for i, rec := range report.Recommendations { if i >= 3 { break }; fmt.Printf("  • %s\n", rec.Title) }\n fmt.Println("\nMonthly projection")\n fmt.Printf("  Runs/month: %.0f\n", report.RunsPerMonth)\n fmt.Printf("  Current runtime: ~%.1fh\n", report.MonthlyRuntimeHours)\n if report.MonthlyCostUSD > 0 { fmt.Printf("  Estimated cost: $%.2f\n", report.MonthlyCostUSD); fmt.Printf("  Potential saving: $%.2f (%.0f%% scenario)\n", report.PotentialMonthlySavingsUSD, report.SavingsScenarioPercent) } else { fmt.Println("  Estimated cost: set -cost-per-minute") }\n fmt.Printf("\nAnalyzed %d workflow runs at %s\n", report.RunsAnalyzed, time.Now().UTC().Format(time.RFC3339))
 fmt.Printf("Average critical path: %.1f minutes\n", report.AverageCriticalPathSeconds/60)
 if report.CacheSignals > 0 { fmt.Printf("Cache/setup signals observed: %d\n", report.CacheSignals) }
}
