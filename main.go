package main

import (
	"fmt"
	"os"

	"github.com/bcjordan/brian/cmd/brian"
)

func main() {
	if err := brian.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}