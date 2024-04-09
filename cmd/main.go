package main

import (
	"github.com/Pactus-Contrib/Indexer/cmd/commands"
	_ "go.uber.org/automaxprocs"
)

func main() {
	commands.Execute()
}
