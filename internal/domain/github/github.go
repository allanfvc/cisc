package github

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/allanfvc/cisc/utils"
	"github.com/hasura/go-graphql-client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

const (
	GraphqlEndpoint = "/graphql"
)

type GitHub struct {
	ApiURl        string
	GraphqlClient *graphql.Client
	RestClient    *http.Client
}

func NewGithubClient(token string, url string) *GitHub {
  validateGithubClientParams(token, url)
  src := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: token},
  )
  httpClient := oauth2.NewClient(context.Background(), src)
  client := graphql.NewClient(url+GraphqlEndpoint, httpClient)
  return &GitHub{ApiURl: url, GraphqlClient: client, RestClient: httpClient}
}

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
	perPage := 4

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
  body2 := make(map[string]interface{})
  err = json.Unmarshal(body, &body2)
	response := &WorkflowRunResponse{}
	err = json.Unmarshal(body, response)
	var runs []WorkflowRun
	for _, run := range response.WorkflowRun {
		if !run.IsLogExpired() {
			content, err := g.getLog(run.LogsUrl, run.ID)
			if err != nil {
				log.Infof("removed run with ID: %d due log expiration", run.ID)
			} else {
				run.LogContent = content
			}
		}
		run.duration()
		runs = append(runs, run)
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
