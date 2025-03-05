package main

import (
	"fmt"
	"os"

	//nolint:depguard
	"github.com/spf13/cobra"
)

func main() {
	// Place your code here.
	if err := buildCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func buildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-envdir [env-dir] [command] [args...]",
		Short: "Run a command with environment variables from a directory",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			envDir := args[0]
			command := args[1:]

			env, err := ReadDir(envDir)
			if err != nil {
				return fmt.Errorf("reading environment directory: %w", err)
			}

			exitCode := RunCmd(command, env)
			os.Exit(exitCode)
			return nil
		},
	}
	return cmd
}
