package build

import "github.com/allanfvc/cisc/internal/domain/miner"

type Runner interface {
	Run()
  Config(config RunConfig, detector miner.ICIDetector) Runner
  GetName() string
}

type RunConfig struct {
	Owner          string
	Project        string
	MainBranchName string
}
