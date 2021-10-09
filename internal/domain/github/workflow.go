package github

import (
	"encoding/json"
	"fmt"
	"time"
)

type Workflow struct {
	ID     int    `json:"id"`
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
	State  string `json:"state"`
}

func (w *Workflow) String() string {
	out, _ := json.Marshal(w)
	return string(out)
}

type WorkflowResponse struct {
	TotalCount int        `json:"total_count"`
	Workflows  []Workflow `json:"workflows"`
}

type WorkflowRun struct {
	ID            int        `json:"id"`
	NodeId        string     `json:"node_id"`
	Branch        string     `json:"head_branch"`
	WorkflowID    int        `json:"workflow_id"`
	LogsUrl       string     `json:"logs_url"`
	LogContent    string     `json:"-"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Conclusion    string     `json:"conclusion"`
	Duration      int64      `json:"duration"`
	DurationHuman string     `json:"readable_duration"`
	Event         string     `json:"event"`
	HeadCommit    HeadCommit `json:"head_commit"`
}

func (w *WorkflowRun) IsLogExpired() bool {
	days := time.Now().Sub(w.CreatedAt).Hours() / 24
	return days > 90
}

func (w *WorkflowRun) duration() {
	duration := w.UpdatedAt.Sub(w.CreatedAt)
	w.Duration = duration.Milliseconds()
	w.DurationHuman = fmt.Sprintf("%v", duration)

}

type WorkflowRunJobStep struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	Number      int       `json:"number"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type WorkflowRunJob struct {
	ID          int              `json:"id"`
	RunID       int              `json:"run_id"`
	Status      string           `json:"status"`
	Conclusion  string           `json:"conclusion"`
	StartedAt   time.Time        `json:"started_at"`
	CompletedAt time.Time        `json:"completed_at"`
	Name        string           `json:"name"`
	Steps       []WorkflowRunJob `json:"steps"`
}

func (w *WorkflowRun) String() string {
	out, _ := json.Marshal(w)
	return string(out)
}

type WorkflowRunResponse struct {
	TotalCount  int           `json:"total_count"`
	WorkflowRun []WorkflowRun `json:"workflow_runs"`
}

type HeadCommit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Author    GitUser    `json:"author"`
	Committer  GitUser    `json:"committer"`
}

type GitUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
