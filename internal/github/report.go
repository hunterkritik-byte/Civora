package github

import "sort"

type Report struct {
 Repository string `json:"repository"`
 RunsAnalyzed int `json:"runs_analyzed"`
 TotalJobs int `json:"total_jobs"`
 FailedJobs int `json:"failed_jobs"`
 TotalJobSeconds float64 `json:"total_job_seconds"`
 EstimatedRunnerMinutes float64 `json:"estimated_runner_minutes"`
 EstimatedCostUSD float64 `json:"estimated_cost_usd,omitempty"`
 CostRateUSDPerMinute float64 `json:"cost_rate_usd_per_minute,omitempty"`
 TotalCriticalPathSeconds float64 `json:"total_critical_path_seconds"`
 TotalParallelRuns int `json:"parallel_runs"`
 AverageCriticalPathSeconds float64 `json:"average_critical_path_seconds"`
 CacheSignals int `json:"cache_signals_observed"`
 Jobs []JobSummary `json:"jobs"`
 Steps []StepSummary `json:"steps"`
 Recommendations []Recommendation `json:"recommendations"`
}

type JobSummary struct {
 Name string `json:"name"`
 Runs int `json:"runs"`
 TotalSeconds float64 `json:"total_seconds"`
 AverageSeconds float64 `json:"average_seconds"`
 Failed bool `json:"failed"`
}

type StepSummary struct {
 Name string `json:"name"`
 Runs int `json:"runs"`
 TotalSeconds float64 `json:"total_seconds"`
 AverageSeconds float64 `json:"average_seconds"`
}

type Recommendation struct {
 Severity string `json:"severity"`
 Title string `json:"title"`
 Reason string `json:"reason"`
}

func ApplyCost(r *Report, usdPerMinute float64) {
 if usdPerMinute > 0 {
  r.CostRateUSDPerMinute = usdPerMinute
  r.EstimatedCostUSD = r.EstimatedRunnerMinutes * usdPerMinute
 }
}

func Recommend(r Report) []Recommendation {
 jobs := append([]JobSummary(nil), r.Jobs...)
 sort.Slice(jobs, func(i, j int) bool { return jobs[i].AverageSeconds > jobs[j].AverageSeconds })

 var out []Recommendation
 for i, j := range jobs {
  if i >= 5 {
   break
  }
  if j.AverageSeconds >= 600 {
   out = append(out, Recommendation{"high", "Investigate slow job: " + j.Name, "Average duration is over 10 minutes."})
  } else if j.AverageSeconds >= 300 {
   out = append(out, Recommendation{"medium", "Review job: " + j.Name, "Average duration is over 5 minutes."})
  }
  if j.Failed && j.Runs >= 3 {
   out = append(out, Recommendation{"high", "Investigate recurring failures: " + j.Name, "This job has failed in analyzed runs."})
  }
 }

 if r.TotalJobs > 0 && r.FailedJobs*5 >= r.TotalJobs {
  out = append(out, Recommendation{"high", "CI reliability issue", "At least 20% of analyzed jobs failed."})
 }

 if r.TotalParallelRuns > 0 && r.AverageCriticalPathSeconds > 0 && r.TotalJobSeconds/r.AverageCriticalPathSeconds > 1.5 {
  out = append(out, Recommendation{"medium", "Review workflow parallelism", "Aggregate job time is substantially higher than the average critical-path duration; parallelization may be hiding runner waste."})
 }

 if r.CacheSignals > 0 {
  out = append(out, Recommendation{"info", "Inspect cache configuration", "Cache/setup steps were observed. Civora does not infer cache hit rates from metadata alone; inspect cache keys and restore behavior before claiming savings."})
 }

 return out
}
