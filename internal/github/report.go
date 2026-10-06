package github

import "sort"

type Report struct {
 Repository string `json:"repository"`
 RunsAnalyzed int `json:"runs_analyzed"`
 TotalJobs int `json:"total_jobs"`
 FailedJobs int `json:"failed_jobs"`
 TotalJobSeconds float64 `json:"total_job_seconds"`
 Jobs []JobSummary `json:"jobs"`
 Recommendations []Recommendation `json:"recommendations"`
}
type JobSummary struct { Name string `json:"name"`; Runs int `json:"runs"`; TotalSeconds float64 `json:"total_seconds"`; AverageSeconds float64 `json:"average_seconds"`; Failed bool `json:"failed"` }
type Recommendation struct { Severity string `json:"severity"`; Title string `json:"title"`; Reason string `json:"reason"` }

func Recommend(r Report) []Recommendation {
 jobs := append([]JobSummary(nil), r.Jobs...)
 sort.Slice(jobs, func(i,j int) bool { return jobs[i].AverageSeconds > jobs[j].AverageSeconds })
 var out []Recommendation
 for i,j := range jobs {
  if i >= 5 { break }
  if j.AverageSeconds >= 600 { out=append(out, Recommendation{"high","Investigate slow job: "+j.Name,"Average duration is over 10 minutes."}) } else if j.AverageSeconds >= 300 { out=append(out, Recommendation{"medium","Review job: "+j.Name,"Average duration is over 5 minutes."}) }
  if j.Failed && j.Runs >= 3 { out=append(out, Recommendation{"high","Investigate recurring failures: "+j.Name,"This job has failed in analyzed runs."}) }
 }
 if r.TotalJobs > 0 && r.FailedJobs*5 >= r.TotalJobs { out=append(out, Recommendation{"high","CI reliability issue","At least 20% of analyzed jobs failed."}) }
 return out
}
