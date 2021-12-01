package miner

import (
	"context"
	"fmt"
	"github.com/google/go-github/scrape"
	"github.com/google/go-github/v39/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io"
	"sort"
	"time"
)

type GitHubCIDetector struct {
	client  *github.Client
	history *BuildHistory
	context *context.Context
}

func NewGitHubCIDetector(token string, url string) *GitHubCIDetector {
	validateGithubClientParams(token, url)
	ctx := context.Background()
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	github.NewClient(oauth2.NewClient(ctx, src))
	return &GitHubCIDetector{
		client:  github.NewClient(oauth2.NewClient(ctx, src)),
		context: &ctx,
	}
}

func validateGithubClientParams(token string, url string) {
	if token == "" {
		log.Fatal("the github token cannot be empty")
	}

	if url == "" {
		log.Fatal("the github url cannot be empty")
	}
}

func isExpired(w *github.WorkflowRun) bool {
  days := time.Now().Sub(w.CreatedAt.Time).Hours() / 24
  return days > 90
}

func duration(w *github.WorkflowRun) int64 {
  duration := w.UpdatedAt.Sub(w.CreatedAt.Time)
  return duration.Milliseconds()
}

func (g GitHubCIDetector) RetrieveBuildHistory(owner string, project string) (*BuildHistory, error) {
	if g.history != nil {
		return g.history, nil
	}
	workflows, err := g.ListWorkflows(owner, project, 0, 100)
	if err != nil {
		return nil, err
	}
	builds := make(map[*string]map[time.Time]BuildPoint)
	var workflowRuns []*github.WorkflowRun
	for _, workflow := range workflows {
		runs, err := g.ListWorkflowRunsByID(owner, project, 0, 100, *workflow.ID)
		if err != nil {
			return nil, err
		}
    workflowRuns = append(workflowRuns, runs...)
	}
  for _, run := range workflowRuns {
    if !isExpired(run) {
      build := convertWorkFlowRunToBuildPoint(run)
      _, hasKey := builds[build.BuildFeature.Branch]
      if hasKey {
        points := builds[build.BuildFeature.Branch]
        points[build.StartAt] = *build
      } else {
        buildPoint := make(map[time.Time]BuildPoint)
        buildPoint[build.StartAt] = *build
        builds[build.BuildFeature.Branch] = buildPoint
      }
    }
  }
  
	history := &BuildHistory{
		Owner:   owner,
		Project: project,
		Builds:  builds,
	}
	return history, nil
}

func (g GitHubCIDetector) ListWorkflows(owner, repo string, page, perPage int) ([]*github.Workflow, error) {
	workflowsResponse, response, err := g.client.Actions.ListWorkflows(*g.context, owner, repo, &github.ListOptions{Page: page, PerPage: perPage})
	if err == nil {
		workflows := workflowsResponse.Workflows
		if response.NextPage != 0 {
			workflowsNext, err := g.ListWorkflows(owner, repo, response.NextPage, perPage)
			if err == nil {
				workflows = append(workflows, workflowsNext...)
			}
		}
		return workflows, nil
	}
	return nil, err
}
func (g GitHubCIDetector) ListWorkflowRunsByID(owner, repo string, page, perPage int, ID int64) ([]*github.WorkflowRun, error) {
	options := &github.ListWorkflowRunsOptions{}
	options.Page = page
	options.PerPage = perPage
	runsByID, response, err := g.client.Actions.ListWorkflowRunsByID(*g.context, owner, repo, ID, options)
  if err == nil {
    workflowRuns := runsByID.WorkflowRuns
    if response.NextPage != 0 {
      workflowRunsNext, err := g.ListWorkflowRunsByID(owner, repo, response.NextPage, perPage, ID)
      if err == nil {
        workflowRuns = append(workflowRuns, workflowRunsNext...)
      }
    }
    return workflowRuns, nil
  }
	return nil, err
}

func (g GitHubCIDetector) LinearizeBuildHistory(builds *BuildHistory) (map[time.Time]BuildPoint, error) {
	linearHistory := make(map[time.Time]BuildPoint)
	var keys []time.Time
	for _, build := range builds.Builds {
		for key, value := range build {
			linearHistory[key] = value
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})
	orderedHistory := make(map[time.Time]BuildPoint)
	for _, key := range keys {
		orderedHistory[key] = linearHistory[key]
	}
	return orderedHistory, nil
}

func convertWorkFlowRunToBuildPoint(run *github.WorkflowRun) *BuildPoint {
	return &BuildPoint{
		ID:      run.ID,
		StartAt: run.CreatedAt.Time,
		EndAt:   run.UpdatedAt.Time,
		BuildFeature: BuildFeature{
			Branch:      run.HeadBranch,
			Status:      run.Conclusion,
			StartAt:     run.CreatedAt.Time,
			Duration:    duration(run),
			EventType:   run.Event,
			BuildNumber: run.RunNumber,
			Jobs:        getJobs(run.JobsURL),
		},
	}
}

//func getJobs(jobs []github.WorkflowRunJob) []BuildJob {
//	if len(jobs) == 0 {
//		return nil
//	}
//	var buildJobs []BuildJob
//	for _, job := range jobs {
//		buildJobs = append(buildJobs, BuildJob{
//			ID:    job.ID,
//			State: job.Conclusion,
//		})
//	}
//	return buildJobs
//}

func (g GitHubCIDetector) RetrieveLogPath(owner string, project string, ID int) string {
	//logContent, err := g.client.GetLog(owner, project, ID)
	//if err != nil {
	//	log.Error(err)
	//}
	return "logContent"
}

func Sample() {
	client := scrape.NewClient(nil)
	if err := client.Authenticate("allanfvc", "", ""); err != nil {
		log.Fatal(err)
	}
	res, _ := client.Get("https://github.com/speedment/speedment/commit/5a2f8abb4d27b5ea89f6cb2d248693941f0b487d/checks/4096085917/logs/1")
	body, _ := io.ReadAll(res.Body)
	fmt.Print(string(body))
}
