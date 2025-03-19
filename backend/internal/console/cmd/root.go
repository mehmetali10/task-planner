package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task-planner",
	Short: "Task Planner is a CLI tool for managing tasks",
	Long:  `Task Planner is a CLI tool that fetches tasks from providers and processes them using a worker pool.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
