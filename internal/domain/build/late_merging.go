package build

import (
	"fmt"
	"github.com/allanfvc/cisc/internal/domain/miner"
	log "github.com/sirupsen/logrus"
	"time"
)

type LateMergingRunner struct {
	detector miner.ICIDetector
}

func (l LateMergingRunner) Run(config RunConfig) {
	history, err := l.detector.RetrieveBuildHistory(config.Owner, config.Project)
	if err != nil {
		log.Fatal(err)
	}
	l.visitHistory(history, config)
}

func (l LateMergingRunner) visitHistory(buildHistory *miner.BuildHistory, config RunConfig) {
	builds := buildHistory.Builds
	//linearHistory, _ := l.detector.LinearizeBuildHistory(buildHistory)
	buildEvents := make(map[time.Time]miner.BuildEvent)
	for branch, buildPoints := range builds {
		if branch == config.MainBranchName {
			log.Infoln("main")
			//l.recordBuildEventsForMainBranch(buildPoints, buildEvents, builds)
		} else {
			l.recordBuildEventsForBranch(branch, buildPoints, buildEvents, builds[config.MainBranchName], config)
		}
	}
}

func (l LateMergingRunner) recordBuildEventsForBranch(
	branch string, buildPoints map[time.Time]miner.BuildPoint,
	buildEvents map[time.Time]miner.BuildEvent,
	masterHistory map[time.Time]miner.BuildPoint,
	config RunConfig) {
	for key, _ := range buildPoints {
		buildEvents[key] = miner.BuildEvent{Date: key, Branch: branch, EventType: miner.BuildEventFork}
	}

	for key, buildPoint := range buildPoints {
    eventType := l.getEventType(buildPoint, config.MainBranchName, masterHistory)
		buildEvents[key] = miner.BuildEvent{Date: key, Branch: branch, EventType: eventType}
	}
}

func (l LateMergingRunner) getEventType(buildPoint miner.BuildPoint, mainBranchName string, masterHistory map[time.Time]miner.BuildPoint) string {
	if fmt.Sprintf("Merge branch %s", mainBranchName) == buildPoint.VCSfeature.Message {
		return miner.BuildEventSync
	}
	isMasterRebasedToBranch := l.isMasterRebasedToFeatureBranch(buildPoint, masterHistory)
	if isMasterRebasedToBranch {
		return miner.BuildEventSync
	}
	return miner.BuildEventCommit
}

func (l LateMergingRunner) isMasterRebasedToFeatureBranch(buildPoint miner.BuildPoint, masterHistory map[time.Time]miner.BuildPoint) bool {
	for _, masterBP := range masterHistory {
		if masterBP.VCSfeature.Message == buildPoint.VCSfeature.Message {
			isSameCommitter := buildPoint.VCSfeature.CommitterName == masterBP.VCSfeature.CommitterName
			isSameCommitterDate := buildPoint.VCSfeature.CommitterDate == masterBP.VCSfeature.CommitterDate
			if isSameCommitter && isSameCommitterDate {
				return true
			}
		}
	}
	return false
}
