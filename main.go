package main

import (
	"github.com/allanfvc/cisc-action/internal/domain/build"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"os"
  "strings"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&nested.Formatter{
		HideKeys:    false,
		TimestampFormat: "01-02-2006 15:04:05.000",
	})
}

func main() {
	token := os.Getenv("GIHUB_ACCESS_TOKEN")
	url := os.Getenv("GITHUB_API_URL")
  fullRepo := os.Getenv("GITHUB_REPOSITORY")
  split := strings.SplitN(fullRepo, "/", 2)
	owner := split[0]
	repo := split[1]
	slowBuild := build.NewSlowBuildRunner(token, url)
	slowBuild.Run(owner, repo)
}
