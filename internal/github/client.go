package github

import (
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strconv"
 "strings"
 "time"
)

type Client struct { token string; http *http.Client }
func New(token string) *Client { return &Client{token:token,http:&http.Client{Timeout:30*time.Second}} }

type workflowRunsResponse struct { TotalCount int `json:"total_count"`; Runs []WorkflowRun `json:"workflow_runs"` }
type WorkflowRun struct { ID int64 `json:"id"`; Name string `json:"name"`; Status string `json:"status"`; Conclusion string `json:"conclusion"`; CreatedAt string `json:"created_at"`; RunStartedAt string `json:"run_started_at"` }
type jobsResponse struct { Jobs []Job `json:"jobs"` }
type Job struct { ID int64 `json:"id"`; Name string `json:"name"`; Status string `json:"status"`; Conclusion string `json:"conclusion"`; StartedAt string `json:"started_at"`; CompletedAt string `json:"completed_at"`; Steps []Step `json:"steps"` }
type Step struct { Name string `json:"name"`; Status string `json:"status"`; Conclusion string `json:"conclusion"`; StartedAt string `json:"started_at"`; CompletedAt string `json:"completed_at"` }

func (c *Client) get(path string, out any) error {
 req,err:=http.NewRequest(http.MethodGet,"https://api.github.com"+path,nil); if err!=nil{return err}
 req.Header.Set("Accept","application/vnd.github+json"); req.Header.Set("Authorization","Bearer "+c.token); req.Header.Set("X-GitHub-Api-Version","2026-03-10")
 resp,err:=c.http.Do(req); if err!=nil{return err}; defer resp.Body.Close()
 if resp.StatusCode<200||resp.StatusCode>=300 { body,_:=io.ReadAll(io.LimitReader(resp.Body,2048)); return fmt.Errorf("GitHub API %s: %s",resp.Status,strings.TrimSpace(string(body))) }
 return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) AnalyzeRepository(repo string, limit int) (Report,error) {
 parts:=strings.Split(repo,"/"); if len(parts)!=2||parts[0]==""||parts[1]=="" {return Report{},fmt.Errorf("repo must be owner/name")}
 if limit<1{limit=1}; if limit>100{limit=100}
 var runs workflowRunsResponse
 path:="/repos/"+url.PathEscape(parts[0])+"/"+url.PathEscape(parts[1])+"/actions/runs?per_page="+strconv.Itoa(limit)
 if err:=c.get(path,&runs);err!=nil{return Report{},err}
 report:=Report{Repository:repo,RunsAnalyzed:len(runs.Runs)}
 for _,run:=range runs.Runs {
  var jobs jobsResponse
  p:=fmt.Sprintf("/repos/%s/%s/actions/runs/%d/jobs?per_page=100",url.PathEscape(parts[0]),url.PathEscape(parts[1]),run.ID)
  if err:=c.get(p,&jobs);err!=nil{return Report{},err}
  var runCritical float64
  runJobs := len(jobs.Jobs)
  if runJobs > 1 { report.TotalParallelRuns++ }
  for _,job:=range jobs.Jobs {
   report.TotalJobs++
   d:=duration(job.StartedAt,job.CompletedAt)
   seconds:=d.Seconds()
   report.TotalJobSeconds+=seconds
   if seconds>runCritical { runCritical=seconds }
   failed:=job.Conclusion=="failure"
   if failed{report.FailedJobs++}
   for _,step:=range job.Steps {
    sd:=duration(step.StartedAt,step.CompletedAt).Seconds()
    if sd<=0 { continue }
    s:=findStep(report.Steps,step.Name)
    if s==nil { report.Steps=append(report.Steps,StepSummary{Name:step.Name,Runs:1,TotalSeconds:sd,AverageSeconds:sd}); s=&report.Steps[len(report.Steps)-1] } else { s.Runs++; s.TotalSeconds+=sd; s.AverageSeconds=s.TotalSeconds/float64(s.Runs) }
    n:=strings.ToLower(step.Name)
    if strings.Contains(n,"cache") || strings.Contains(n,"setup-node") || strings.Contains(n,"setup-python") || strings.Contains(n,"setup-go") { report.CacheSignals++ }
   }
   found:=false
   for i:=range report.Jobs { if report.Jobs[i].Name==job.Name { report.Jobs[i].Runs++; report.Jobs[i].TotalSeconds+=seconds; report.Jobs[i].AverageSeconds=report.Jobs[i].TotalSeconds/float64(report.Jobs[i].Runs); report.Jobs[i].Failed=report.Jobs[i].Failed||failed; found=true; break } }
   if !found { report.Jobs=append(report.Jobs,JobSummary{Name:job.Name,Runs:1,TotalSeconds:seconds,AverageSeconds:seconds,Failed:failed}) }
  }
  report.TotalCriticalPathSeconds += runCritical
 }
 report.EstimatedRunnerMinutes=report.TotalJobSeconds/60
 report.AverageCriticalPathSeconds=report.TotalCriticalPathSeconds/float64(max(1,report.RunsAnalyzed))
 report.Recommendations=Recommend(report)
 return report,nil
}

func findStep(steps []StepSummary,name string)*StepSummary { for i:=range steps { if steps[i].Name==name{return &steps[i]} }; return nil }
func duration(start,end string) time.Duration { if start==""||end==""{return 0}; a,e1:=time.Parse(time.RFC3339Nano,start); b,e2:=time.Parse(time.RFC3339Nano,end); if e1!=nil||e2!=nil||b.Before(a){return 0}; return b.Sub(a) }
func max(a,b int) int { if a>b{return a}; return b }
