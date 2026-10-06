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
 flag.Parse()
 if *repo == "" { fmt.Fprintln(os.Stderr, "usage: civora -repo owner/name [-runs 30]"); os.Exit(2) }
 token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
 if token == "" { fmt.Fprintln(os.Stderr, "GITHUB_TOKEN is required"); os.Exit(2) }
 report, err := github.New(token).AnalyzeRepository(*repo, *runs)
 if err != nil { fmt.Fprintln(os.Stderr, "civora:", err); os.Exit(1) }
 b, _ := json.MarshalIndent(report, "", "  ")
 fmt.Println(string(b))
 fmt.Printf("\nAnalyzed %d workflow runs at %s\n", report.RunsAnalyzed, time.Now().UTC().Format(time.RFC3339))
}
