package github

import (
  "encoding/json"
  "fmt"
  "github.com/allanfvc/cisc/utils"
  log "github.com/sirupsen/logrus"
  "io"
)

const (
	GraphqlEndpoint = "/graphql"
)

func validateGithubClientParams(token string, url string) {
  if token == "" {
    log.Fatal("the github token cannot be empty")
  }

  if url == "" {
    log.Fatal("the github url cannot be empty")
  }
}

func (g GitHub) ListWorkflows(owner, repo string) ([]Workflow, error) {
	page := 0
	perPage := 100

	response, err := g.listPagedWorkflows(owner, repo, page, perPage)
	if err == nil {
		workflows := response.Workflows
		runs := getRuns(response.TotalCount, perPage) - 1
		for i := 0; i < runs; i++ {
			page++
			response, _ := g.listPagedWorkflows(owner, repo, page, perPage)
			workflows = append(workflows, response.Workflows...)
		}
		return workflows, err
	}
	return nil, err
}

func (g GitHub) listPagedWorkflows(owner, repo string, page, perPage int) (*WorkflowResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/workflows?page=%d&per_page=%d", g.ApiURl, owner, repo, page, perPage)
	res, err := g.request(url, "GET", nil)
	if err == nil {
		body, _ := io.ReadAll(res.Body)
		workflows := &WorkflowResponse{}
		err = json.Unmarshal(body, workflows)
		return workflows, err
	}
	return nil, err
}

func (g *GitHub) ListWorkflowRunsByID(owner, repo string, ID int) ([]WorkflowRun, error) {
	page := 0
	perPage := 100

	response, err := g.listPagedWorkflowRunsByID(owner, repo, ID, page, perPage)
	if err != nil {
		return nil, err
	}
	workflows := response.WorkflowRun
	runs := getRuns(response.TotalCount, perPage) - 1
	for i := 0; i < runs; i++ {
		page++
		response, _ := g.listPagedWorkflowRunsByID(owner, repo, ID, page, perPage)
		workflows = append(workflows, response.WorkflowRun...)
	}
	return workflows, err
}

func (g GitHub) listPagedWorkflowRunsByID(owner, repo string, ID, page, perPage int) (*WorkflowRunResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/workflows/%d/runs?page=%d&per_page=%d", g.ApiURl, owner, repo, ID, page, perPage)
	res, err := g.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(res.Body)
	response := &WorkflowRunResponse{}
	err = json.Unmarshal(body, response)
	var runs []WorkflowRun
	for _, run := range response.WorkflowRun {
		if !run.IsLogExpired() {
			//content, err := g.getLog(run.LogsUrl, run.ID)
			//if err != nil {
			//	log.Infof("removed run with ID: %d due log expiration", run.ID)
			//} else {
			//	run.LogContent = content
			//}
      jobs, _ := g.getJobs(run.JobsUrl)
      run.Jobs = jobs
      for _, job := range jobs {
        run.LogContent = g.getJobLog(job, owner, repo)
      }
      run.duration()
      runs = append(runs, run)
		}

	}
	response.WorkflowRun = runs
	return response, err
}

func (g GitHub) getLog(url string, id int) (string, error) {
	res, err := g.request(url, "GET", nil)
	if err == nil {
		body, _ := io.ReadAll(res.Body)
		if isZip(res) {
			message, err := utils.ReadZipFiles(body)
			return message, err
		} else if res.StatusCode == 410 {
			return "", fmt.Errorf("failed to get log for run with id: %v, due log expiration", id)
		}
	}
	return "", fmt.Errorf("failed to get log for run with id: %v, with error %v", id, err)
}

func (g GitHub) GetLog(owner, repo string, ID int) (string, error) {
  url := fmt.Sprintf("%s/repos/%s/%s/actions/jobs/%d/logs", g.ApiURl, owner, repo, ID)
  return g.getLog(url, ID)
}

func (g GitHub) getJobs(url string) ([]WorkflowRunJob, error) {
  res, err := g.request(url, "GET", nil)
  if err == nil {
    body, _ := io.ReadAll(res.Body)
    response := &WorkflowRunJobResponse{}
    err = json.Unmarshal(body, response)
    if err != nil {
      return nil, err
    }
    return response.Jobs, nil
  }
  return nil, fmt.Errorf("failed to get jobs for url %s, with error %v", url, err)
}

func (g GitHub) getJobLog(job WorkflowRunJob, owner, repo string) string {
  log := ""
  for _, step := range job.Steps {
    log += g.getStepLog(step.Number, job, owner, repo)
  }
  return log
}

func (g GitHub) getStepLog(ID int, job WorkflowRunJob, owner, repo string) string {
  url := fmt.Sprintf("https://github.com/%s/%s/commit/%s/checks/%d/logs/%d", owner, repo, job.HeadSha, job.ID, ID)
  res, err := g.request(url, "GET", nil)
  if err == nil {
    fmt.Print(res)
  }
  return ""
}
