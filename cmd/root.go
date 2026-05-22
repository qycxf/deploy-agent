package cmd

import (
    "github.com/spf13/cobra"
    "fmt"
)

var rootCmd = &cobra.Command{
    Use:   "deploy-agent",
    Short: "deploy-agent CLI",
    Long:  "Deploy Agent - HTTP server and tooling",
}

// Execute runs the root command.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
    }
}
