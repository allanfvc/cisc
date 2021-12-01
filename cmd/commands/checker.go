package commands

import (
  "github.com/allanfvc/cisc/internal/domain/build"
  "github.com/allanfvc/cisc/internal/domain/miner"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "os"
  "strings"
  "time"
)

var checker = &cobra.Command{
  Use:     "checker",
  Aliases: []string{},
  Short:   "Run a smell check on the CI/CD workflow",
  Run:     run,
}

func init() {
  checker.Flags().StringP("ci-tool", "c", "github", "Type of CI/CD tool used: github")
}

func run(cmd *cobra.Command, args []string) {
  manager, _ := cmd.Flags().GetString("ci-tool")
  switch manager {
  case "github":
    checkGithubAction()
  default:
    log.Fatalf("unknow CI/CD tool %s", manager)
  }
}

func checkGithubAction() {
  start := time.Now()
  token := os.Getenv("GITHUB_ACCESS_TOKEN")
  url := os.Getenv("GITHUB_API_URL")
  fullRepo := "speedment/speedment"//os.Getenv("GITHUB_REPOSITORY")
  split := strings.SplitN(fullRepo, "/", 2)
  owner := split[0]
  repo := split[1]
  config := build.RunConfig{Owner: owner, Project: repo}
  detector := miner.NewGitHubCIDetector(token, url)
  runners := getRunners()
  log.Infof("Evaluating ---- %s", fullRepo)
  for i, runner := range runners {
    log.Infof("---- STEP %d: %s ----", i+1, runner.GetName())
    runner.Config(config, detector).Run()
  }
  
  
  //slowBuild := build.NewSlowBuildRunner(detector)
  //log.Infof("---- STEP %d: %s ----")
  //slowBuild.Run(config)
  //log.Info("---- STEP 02: skiped tests ----")
  //slowBuild.Run(config)
  log.Infof("---- Evaluation completed in %v", time.Now().Sub(start))
}

func getRunners() []build.Runner {
  return []build.Runner{
    //build.SlowBuildRunner{},
    build.SkipTestsRunner{},
  }
}

func initSkipTestsChecker(detector miner.ICIDetector) build.Runner {
  return build.NewSkipTestsRunner(detector)
}



