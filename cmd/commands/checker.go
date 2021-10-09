package commands

import (
  "github.com/allanfvc/cisc/internal/domain/build"
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
  Run:     runUpdateDependent,
}

func init() {
  checker.Flags().StringP("ci-tool", "c", "github", "Type of CI/CD tool used: github")
}

func runUpdateDependent(cmd *cobra.Command, args []string) {
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
  fullRepo := os.Getenv("GITHUB_REPOSITORY")
  split := strings.SplitN(fullRepo, "/", 2)
  owner := split[0]
  repo := split[1]
  config := build.RunConfig{Owner: owner, Project: repo}
  log.Infof("Evaluating.... %s", fullRepo)
  slowBuild := build.NewSlowBuildRunner(token, url)
  log.Info("---- STEP1: slow builds")
  slowBuild.Run(config)
  log.Infof("---- Evaluation completed in %v", time.Now().Sub(start))
}



