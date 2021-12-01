package smell

import "github.com/allanfvc/cisc/internal/domain/miner"

type Smell interface {
  Detect(builds *miner.BuildHistory)
}
