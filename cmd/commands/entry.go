package commands

import (
  "fmt"
  nested "github.com/antonfisher/nested-logrus-formatter"
  log "github.com/sirupsen/logrus"
  "os"

  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "cisc",
  Short: "continuous integration smell checker tools to automation.",
  Run:   defaultRun,
}

var commands = []*cobra.Command{checker}

func init() {
  log.SetOutput(os.Stdout)
  log.SetFormatter(&nested.Formatter{
    HideKeys:    false,
    TimestampFormat: "01-02-2006 15:04:05.000",
  })
  for _, cmd := range commands {
    rootCmd.AddCommand(cmd)
  }
}

func Execute(version string) {
  rootCmd.Version = version
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func defaultRun(cmd *cobra.Command, args []string) {
  cmd.Help()
}

