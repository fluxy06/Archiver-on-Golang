package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Архиватор",
}

func HandleEror(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		HandleEror(err)
	}
}
