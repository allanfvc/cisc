package build

type Runner interface {
	Run(owner, project string)
}
