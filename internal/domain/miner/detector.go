package miner

import (
	"encoding/json"
	"time"
)

type ICIDetector interface {
	RetrieveBuildHistory(owner string, project string) (*BuildHistory, error)
}

type BuildHistory struct {
	Owner   string
	Project string
	Builds  map[string]map[time.Time]BuildPoint
}

type BuildPoint struct {
	ID           int `json:"id"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
	BuildFeature BuildFeature `json:"build_feature"`
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
	Branch      string `json:"branch"`
	Status      string `json:"status"`
	StartAt     time.Time `json:"start_at"`
	Duration    int64 `json:"duration"`
	BuildNumber int `json:"build_number"`
	EventType   string `json:"event_type"`
}

type BuildJob struct {
	ID     int
	Number int
	State  string
}