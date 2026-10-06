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

type workflowRunsResponse struct { TotalCount int `json:"total_count"`; Runs []struct { ID int64 `json:"id"`; Name string `json:"name"`; Status string `json:"status"`; Conclusion string `json:"conclusion"` } `json:"workflow_runs"` }
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
  for _,job:=range jobs.Jobs {
   report.TotalJobs++; d:=duration(job.StartedAt,job.CompletedAt); report.TotalJobSeconds+=d.Seconds(); failed:=job.Conclusion=="failure"; if failed{report.FailedJobs++}
   found:=false
   for i:=range report.Jobs { if report.Jobs[i].Name==job.Name { report.Jobs[i].Runs++; report.Jobs[i].TotalSeconds+=d.Seconds(); report.Jobs[i].Failed=report.Jobs[i].Failed||failed; found=true; break } }
   if !found { report.Jobs=append(report.Jobs,JobSummary{Name:job.Name,Runs:1,TotalSeconds:d.Seconds(),Failed:failed}) }
  }
 }
 for i:=range report.Jobs { if report.Jobs[i].Runs>0 {report.Jobs[i].AverageSeconds=report.Jobs[i].TotalSeconds/float64(report.Jobs[i].Runs)} }
 report.Recommendations=Recommend(report); return report,nil
}
func duration(start,end string) time.Duration { if start==""||end==""{return 0}; a,e1:=time.Parse(time.RFC3339Nano,start); b,e2:=time.Parse(time.RFC3339Nano,end); if e1!=nil||e2!=nil||b.Before(a){return 0}; return b.Sub(a) }
