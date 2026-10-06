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
 flag.Parse()
 if *repo == "" { fmt.Fprintln(os.Stderr, "usage: civora -repo owner/name [-runs 30]"); os.Exit(2) }
 token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
 if token == "" { fmt.Fprintln(os.Stderr, "GITHUB_TOKEN is required"); os.Exit(2) }
 report, err := github.New(token).AnalyzeRepository(*repo, *runs)
 if err != nil { fmt.Fprintln(os.Stderr, "civora:", err); os.Exit(1) }
 github.ApplyCost(&report, *costRate)
 b, _ := json.MarshalIndent(report, "", "  ")
 fmt.Println(string(b))
 fmt.Printf("\nAnalyzed %d workflow runs at %s\n", report.RunsAnalyzed, time.Now().UTC().Format(time.RFC3339))
 fmt.Printf("Measured runner time: %.1f minutes\n", report.EstimatedRunnerMinutes)
 fmt.Printf("Average critical path: %.1f minutes\n", report.AverageCriticalPathSeconds/60)
 if report.CacheSignals > 0 { fmt.Printf("Cache/setup signals observed: %d\n", report.CacheSignals) }
}
