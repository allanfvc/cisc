package package_managers

type Execution interface {
  GetCommand() string
  GetTrace() string
  ParseTrace(trace interface{}) map[string][]Test
  GetExecutionTrace() []string
}

type LogParser interface {
  GetExecutions(log string) []Execution
  Parse(executionTrace []string) map[string][]Test
}