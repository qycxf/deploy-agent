package cmd

import (
    "github.com/spf13/cobra"
    "github.com/qycxf/deploy-agent/internal/server"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Run HTTP server",
    RunE: func(cmd *cobra.Command, args []string) error {
        port, _ := cmd.Flags().GetString("port")
        return server.Start(port)
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)
    serverCmd.Flags().StringP("port", "p", "", "port to listen on (overrides env)")
}
