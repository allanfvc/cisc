package miner

import (
  "github.com/allanfvc/cisc/internal/domain/github"
  "sort"
  "time"
)

type GitHubCIDetector struct {
	client *github.GitHub
}

func NewGitHubCIDetector(token string, url string) *GitHubCIDetector {
	return &GitHubCIDetector{client: github.NewGithubClient(token, url)}
}

func (g GitHubCIDetector) RetrieveBuildHistory(owner string, project string) (*BuildHistory, error) {
	workflows, err := g.client.ListWorkflows(owner, project)
	if err != nil {
		return nil, err
	}
	builds := make(map[string]map[time.Time]BuildPoint)
	var runs []github.WorkflowRun
	for _, workflow := range workflows {
		workflowRuns, err := g.client.ListWorkflowRunsByID(owner, project, workflow.ID)
		if err != nil {
			return nil, err
		}
		runs = append(runs, workflowRuns...)
	}

	for _, run := range runs {
		if !run.IsLogExpired() {
			build := convertWorkFlowRunToBuildPoint(run)
			branch := builds[build.BuildFeature.Branch]
			if branch != nil {
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

func convertWorkFlowRunToBuildPoint(run github.WorkflowRun) *BuildPoint {
	return &BuildPoint{
		ID:      run.ID,
		StartAt: run.CreatedAt,
		EndAt:   run.UpdatedAt,
		BuildFeature: BuildFeature{
			Branch:      run.Branch,
			Status:      run.Conclusion,
			StartAt:     run.CreatedAt,
			Duration:    run.Duration,
			EventType:   run.Event,
			BuildNumber: run.ID,
		},
	}
}
