package main

import (
	"os"

	"github.com/NoobAuthor/pomodoro/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
