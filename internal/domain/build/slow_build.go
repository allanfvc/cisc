package build

import (
	"github.com/allanfvc/cisc-action/internal/domain/miner"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

type SlowBuildRunner struct {
	detector miner.ICIDetector
}

func NewSlowBuildRunner(token string, url string) *SlowBuildRunner {
	return &SlowBuildRunner{
		detector: miner.NewGitHubCIDetector(token, url),
	}
}

func (s SlowBuildRunner) Run(owner, project string) {
	history, err := s.detector.RetrieveBuildHistory(owner, project)
	if err != nil {
		log.Fatal(err)
	}
	s.visitHistory(history)
}

func (s SlowBuildRunner) visitHistory(history *miner.BuildHistory) {
	log.Info("key,build_id,duration")
	linearHistory := make(map[time.Time]miner.BuildPoint)
	var keys []time.Time
	for _, build := range history.Builds {
		for key, value := range build {
			linearHistory[key] = value
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})
	for _, key := range keys {
		bp := linearHistory[key]
		log.Info(key, bp.ID, bp.Duration())
	}
	//increase=100*((lm.fit$coefficients[1]+lm.fit$coefficients[2]*nWeeks)/(lm.fit$coefficients[1]+lm.fit$coefficients[2]*1)-1)
}
