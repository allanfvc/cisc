package package_managers

import "time"

type Build struct {
	Log         string
	Modules     []Module
	StartedAt   time.Time
	Duration    time.Duration
	TotalMemory int
	Tests       map[string][]Test
}
