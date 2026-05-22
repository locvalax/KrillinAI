package main

import (
	"fmt"
	"os"

	"github.com/krillinai/KrillinAI/internal/app"
	"github.com/krillinai/KrillinAI/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// Commit is the git commit hash set at build time
	Commit = "none"
	// BuildDate is the build date set at build time
	BuildDate = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "krillinai",
		Short: "KrillinAI - AI-powered video subtitle and translation tool",
		Long: `KrillinAI is an AI-powered tool for generating, translating,
and dubbing subtitles for videos using state-of-the-art AI models.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApp(cmd)
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level (debug, info, warn, error)")

	// Version subcommand
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("KrillinAI\n")
			fmt.Printf("  Version:    %s\n", Version)
			fmt.Printf("  Commit:     %s\n", Commit)
			fmt.Printf("  Build Date: %s\n", BuildDate)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runApp initializes and starts the KrillinAI application server.
func runApp(cmd *cobra.Command) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	logLevel, _ := cmd.Flags().GetString("log-level")

	// Load configuration
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override log level if specified via flag
	if logLevel != "" {
		cfg.LogLevel = logLevel
	}

	// Initialize and run the application
	application, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	return application.Run()
}
