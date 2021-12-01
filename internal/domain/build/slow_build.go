package build

import (
  "github.com/allanfvc/cisc/internal/domain/miner"
  "github.com/allanfvc/cisc/utils"
  "github.com/pterm/pterm"
  log "github.com/sirupsen/logrus"
  "sort"
  "time"
)

type SlowBuildRunner struct {
	detector miner.ICIDetector
  config RunConfig
}

func (r SlowBuildRunner) Config(config RunConfig, detector miner.ICIDetector) Runner {
  r.config = config
  r.detector = detector
  return r
}

func (r SlowBuildRunner) GetName() string {
  return "slow builds"
}

func (r SlowBuildRunner) Run() {
	builds, err := r.detector.RetrieveBuildHistory(r.config.Owner, r.config.Project)
	if err != nil {
		log.Fatal(err)
	}
	r.visitHistory(builds)
}

func (r SlowBuildRunner) visitHistory(builds *miner.BuildHistory) {
	history, _ := r.detector.LinearizeBuildHistory(builds)
	r.calculate(history)
}

func (r SlowBuildRunner) calculate(history map[time.Time]miner.BuildPoint) {
	durationByWeek := make(map[string][]int)
	for _, bp := range history {
		week := utils.WeekStartDate(bp.StartAt).Format("2006-01-02")
		value := durationByWeek[week]
		value = append(value, bp.DurationInSeconds())
		durationByWeek[week] = value
	}
	bars := pterm.Bars{}
	for week, durations := range durationByWeek {
		avg := utils.Average(durations...)
		bars = append(bars, pterm.Bar{
			Label: week,
			Value: int(avg),
		})
	}
  sort.Slice(bars, func(i, j int) bool {
    return bars[i].Label < bars[j].Label
  })
	_ = pterm.DefaultBarChart.WithBars(bars).WithShowValue(true).Render()
}
