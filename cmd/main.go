package main

import "github.com/allanfvc/cisc/cmd/commands"

var version = "0.0.1"

func main() {
  commands.Execute(version)
}

