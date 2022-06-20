package main

import (
	"fmt"

	"go.bbkane.com/warg/command"

	_ "embed"
)

//go:embed embedded/examples/repolines.sh
var examplesLinesText string

//go:embed embedded/examples/langlines.sh
var examplesLangLinesText string

func printString(s string) command.Action {
	return func(_ command.Context) error {
		fmt.Println(examplesLinesText)
		return nil
	}
}
