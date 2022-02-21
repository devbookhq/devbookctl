package main

import (
	"log"

	"github.com/devbookhq/devbookctl/cmd"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	cmd.Execute()
}
