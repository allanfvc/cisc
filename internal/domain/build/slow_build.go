package build

import (
  "github.com/allanfvc/cisc/internal/domain/miner"
  log "github.com/sirupsen/logrus"
)

type SlowBuildRunner struct {
	detector miner.ICIDetector
}

func NewSlowBuildRunner(token string, url string) *SlowBuildRunner {
	return &SlowBuildRunner{
		detector: miner.NewGitHubCIDetector(token, url),
	}
}

func (s SlowBuildRunner) Run(config RunConfig) {
	builds, err := s.detector.RetrieveBuildHistory(config.Owner, config.Project)
	if err != nil {
		log.Fatal(err)
	}
	s.visitHistory(builds)
}

func (s SlowBuildRunner) visitHistory(builds *miner.BuildHistory) {
	log.Info("key,build_id,duration")
  history, _ :=s.detector.LinearizeBuildHistory(builds)
  for key, bp := range history {
    log.Info(key, bp.ID, bp.Duration())
  }
	//increase=100*((lm.fit$coefficients[1]+lm.fit$coefficients[2]*nWeeks)/(lm.fit$coefficients[1]+lm.fit$coefficients[2]*1)-1)
}
