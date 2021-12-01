package package_managers

import "time"

type Test struct {
	Class    string
	Module   string
	TestRun  int
	Failures int
	Errors   int
	Skipped  int
  Duration time.Duration
}
