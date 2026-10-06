package main

import (
 "encoding/json"
 "flag"
 "fmt"
 "os"
 "sort"
 "strings"

 "github.com/hunterkritik-byte/Civora/internal/github"
)

func main() {
 repo := flag.String("repo", "", "GitHub repository (owner/name)")
 runs := flag.Int("runs", 30, "workflow runs to inspect")
 costRate := flag.Float64("cost-per-minute", 0, "estimated runner cost in USD per minute")
 monthlyRuns := flag.Float64("monthly-runs", 30, "estimated workflow runs per month")
 savingsPct := flag.Float64("savings-scenario", 20, "potential savings scenario percentage")
 jsonOutput := flag.Bool("json", false, "print the full JSON report")
 flag.Parse()

 if *repo == "" {
  fmt.Fprintln(os.Stderr, "usage: civora -repo owner/name [-runs 30]")
  os.Exit(2)
 }
 token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
 if token == "" {
  fmt.Fprintln(os.Stderr, "GITHUB_TOKEN is required")
  os.Exit(2)
 }
 report, err := github.New(token).AnalyzeRepository(*repo, *runs)
 if err != nil {
  fmt.Fprintln(os.Stderr, "civora:", err)
  os.Exit(1)
 }
 github.ApplyCost(&report, *costRate)
 github.BuildBusinessMetrics(&report, *monthlyRuns, *costRate, *savingsPct)

 if *jsonOutput {
  b, _ := json.MarshalIndent(report, "", "  ")
  fmt.Println(string(b))
  return
 }

 fmt.Println("CIVORA")
 fmt.Println("────────────────────")
 fmt.Printf("CI Waste Score       %d/100\n", report.WasteScore)
 fmt.Printf("\nRunner time          %.1f min\n", report.EstimatedRunnerMinutes)
 fmt.Printf("Failure rate         %.1f%%\n", report.FailureRatePercent)

 fmt.Println("\nTop waste")
 steps := append([]github.StepSummary(nil), report.Steps...)
 sort.Slice(steps, func(i, j int) bool { return steps[i].TotalSeconds > steps[j].TotalSeconds })
 for i, step := range steps {
  if i >= 2 { break }
  fmt.Printf("  %-18s %.0fs\n", step.Name, step.TotalSeconds)
 }

 fmt.Println("\nPotential improvements")
 for i, rec := range report.Recommendations {
  if i >= 3 { break }
  fmt.Printf("  • %s\n", rec.Title)
 }

 fmt.Println("\nMonthly projection")
 fmt.Printf("  Runs/month: %.0f\n", report.RunsPerMonth)
 fmt.Printf("  Current runtime: ~%.1fh\n", report.MonthlyRuntimeHours)
 if report.MonthlyCostUSD > 0 {
  fmt.Printf("  Estimated cost: $%.2f\n", report.MonthlyCostUSD)
  fmt.Printf("  Potential saving: $%.2f (%.0f%% scenario)\n", report.PotentialMonthlySavingsUSD, report.SavingsScenarioPercent)
 } else {
  fmt.Println("  Estimated cost: set -cost-per-minute")
 }
}

