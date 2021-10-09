package miner

import (
	"encoding/json"
	"time"
)

const BuildEventFork = "fork"
const BuildEventSync = "sync"
const BuildEventCommit = "commit"

type ICIDetector interface {
	RetrieveBuildHistory(owner string, project string) (*BuildHistory, error)
	LinearizeBuildHistory(builds *BuildHistory) (map[time.Time]BuildPoint, error)
}

type BuildHistory struct {
	Owner   string
	Project string
	Builds  map[string]map[time.Time]BuildPoint
}

type BuildPoint struct {
	ID           int          `json:"id"`
	StartAt      time.Time    `json:"start_at"`
	EndAt        time.Time    `json:"end_at"`
	BuildFeature BuildFeature `json:"build_feature"`
	VCSfeature   VCSfeature   `json:"vcs_feature"`
}

func (bp *BuildPoint) String() string {
	out, _ := json.Marshal(bp)
	return string(out)
}

func (bp *BuildPoint) Duration() time.Duration {
	duration := bp.EndAt.Sub(bp.StartAt)
	return duration
}

type BuildFeature struct {
	Branch      string    `json:"branch"`
	Status      string    `json:"status"`
	StartAt     time.Time `json:"start_at"`
	Duration    int64     `json:"duration"`
	BuildNumber int       `json:"build_number"`
	EventType   string    `json:"event_type"`
}

type BuildJob struct {
	ID     int
	Number int
	State  string
}

type BuildEvent struct {
	Date      time.Time
	Branch    string
	EventType string
}

type VCSfeature struct {
	SHA           string `json:"sha"`
	FullBranch    string `json:"full_branch"`
	Message       string `json:"message"`
	CommitterName string `json:"committer_name"`
	CommitterDate string `json:"committer_date"`
	Type          int    `json:"type"`
}
