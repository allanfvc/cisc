package smell

import (
	"fmt"
	"github.com/allanfvc/cisc/internal/domain/miner"
	pkg "github.com/allanfvc/cisc/internal/domain/package_managers"
	"time"
)

type BuildFixPair struct {
	A *miner.BuildPoint
	B *miner.BuildPoint
}

type SkipFailedTests struct {
	NoTaggedBranches  bool
	Owner             string
	Project           string
	FixPairsPerBranch map[string]map[time.Time]BuildFixPair
	Detector          miner.ICIDetector
	Parser            pkg.LogParser
	smells            []string
}

func (s SkipFailedTests) Detect(builds *miner.BuildHistory) {
	history := builds.Builds
	if s.NoTaggedBranches {
		history = miner.RemoveTaggedBranches(history)
	}

	for branch, buildPoints := range history {
		pairs := s.computeBuildFixPairs(buildPoints)
    if s.FixPairsPerBranch == nil {
      s.FixPairsPerBranch = make(map[string]map[time.Time]BuildFixPair)
    }
		s.FixPairsPerBranch[branch] = pairs
	}
	s.calculate(s.Owner, s.Project)
}

func (s SkipFailedTests) computeBuildFixPairs(buildPoints map[time.Time]miner.BuildPoint) map[time.Time]BuildFixPair {
	buildFixPairs := make(map[time.Time]BuildFixPair)
	for date, buildPointA := range buildPoints {
		if buildPointA.BuildFeature.Status == "failure" {
			buildPointB := s.findNextBuildPair(buildPoints, date)
			fixPair := BuildFixPair{A: &buildPointA, B: buildPointB}
			buildFixPairs[date] = fixPair
		}
	}
	return buildFixPairs
}

func (s SkipFailedTests) findNextBuildPair(buildPoints map[time.Time]miner.BuildPoint, date time.Time) *miner.BuildPoint {

	for curDate, buildPointB := range buildPoints {
		if curDate == date {
			status := buildPointB.BuildFeature.Status
			if status == "passed" || status == "failure" {
				return &buildPointB
			}
		}
	}
	return nil
}

func (s SkipFailedTests) calculate(owner, project string) {
	for _, fixPairs := range s.FixPairsPerBranch {
		buildPointA, buildPointB := s.getBuildPoints(fixPairs)
		for _, jobFromA := range buildPointA.BuildFeature.Jobs {
			if jobFromA.State == "failure" {
				for _, jobFromB := range buildPointB.BuildFeature.Jobs {
					if jobFromB.Number == jobFromA.Number {
						logFromJobA := s.Detector.RetrieveLogPath(owner, project, jobFromA.ID)
						logFromJobB := s.Detector.RetrieveLogPath(owner, project, jobFromB.ID)
						s.compareJobs(buildPointB.ID, buildPointB.StartAt, jobFromB.ID, logFromJobA, logFromJobB)
					}
				}
			}
		}
	}
}

func (s SkipFailedTests) getBuildPoints(fixPairs map[time.Time]BuildFixPair) (*miner.BuildPoint, *miner.BuildPoint) {
	var buildPointA *miner.BuildPoint
	var buildPointB *miner.BuildPoint
	for _, buildFixPair := range fixPairs {
		buildPointA = buildFixPair.A
		buildPointB = buildFixPair.B
		if buildFixPair.B != nil {
			break
		}
	}
	return buildPointA, buildPointB
}

func (s SkipFailedTests) compareJobs(ID int, startDate time.Time, jobID int, logA string, logB string) {
	executionsFromA := s.extractExecutions(logA)
	executionsFromB := s.extractExecutions(logB)
	var testModuleA map[string][]pkg.Test
	var testModuleB map[string][]pkg.Test
	var executionFromB pkg.Execution
	for key, execution := range executionsFromA {
		executionFromB, hasKey := executionsFromB[key]
		testModuleA = s.getTestsByModule(execution)
		if hasKey {
			testModuleB = s.getTestsByModule(executionFromB)
			break
		}
	}
	if testModuleA == nil {
		return
	}
	for module, tests := range testModuleA {
		testsB := testModuleB[module]
		for _, test := range tests {
			for _, testB := range testsB {
				if test.Class == testB.Class {
					smell := s.computeAndStoreSmell(test, testB)
					if smell != "" {
						line := fmt.Sprintf("%d,%v,%d,%s,%s,%s", ID, startDate, jobID, executionFromB.GetCommand(), module, smell)
						s.smells = append(s.smells, line)
					}
					break
				}
			}
		}

	}

}

func (s SkipFailedTests) extractExecutions(logContent string) map[string]pkg.Execution {
	executionsMap := make(map[string]pkg.Execution)
	executions := s.Parser.GetExecutions(logContent)
	for _, execution := range executions {
		executionsMap[execution.GetCommand()] = execution
	}
	return executionsMap
}

func (s SkipFailedTests) getTestsByModule(execution pkg.Execution) map[string][]pkg.Test {
	executionTrace := execution.GetExecutionTrace()
	return s.Parser.Parse(executionTrace)
}

func (s SkipFailedTests) computeAndStoreSmell(testA, testB pkg.Test) string {
	className := testA.Class
	breaksTestA := testA.Errors + testA.Failures
	testCasesTestA := testA.TestRun
	ignoredTestA := testA.Skipped
	breaksTestB := testB.Errors + testB.Failures
	ignoredTestB := testB.Skipped
	testCasesTestB := testB.TestRun
	deltaBreaks := breaksTestB - breaksTestA
	deltaIgnored := ignoredTestB - ignoredTestA
	deltaTestCases := testCasesTestB - testCasesTestA
	warning := deltaBreaks < 0 && (deltaTestCases < 0 || deltaIgnored > 0)
	smell := fmt.Sprintf("%s,%d,%d,%d,%v", className, deltaBreaks, deltaIgnored, deltaTestCases, warning)
	if warning {
		return smell
	}
	return ""
}
