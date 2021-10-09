package build

type Runner interface {
	Run(config RunConfig)
}

type RunConfig struct {
	Owner          string
	Project        string
	MainBranchName string
}
