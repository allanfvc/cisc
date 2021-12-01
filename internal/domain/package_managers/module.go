package package_managers

type Module struct {
	GoalRegex       string
	ExtractionRegex string
	Name            string
	Status          string
	Lines           []string
	Duration        string
	//Goals           []Goal
}

type Goal struct {
	Lines   []string
	Plugin  string
	Version string
	Name    string
	Status  string
}
