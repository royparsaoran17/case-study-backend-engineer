package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"auth-se/cmd/genx"
	"auth-se/cmd/http"
	"auth-se/cmd/migration"
	"auth-se/pkg/logger"
)

// Start handler registering service command
func Start() {

	rootCmd := &cobra.Command{}
	logger.SetJSONFormatter()
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()

	migrateCmd := &cobra.Command{
		Use:   "db:migrate",
		Short: "database migration",
		Run: func(c *cobra.Command, args []string) {
			migration.MigrateDatabase()
		},
	}

	migrateCmd.Flags().BoolP("version", "", false, "print version")
	migrateCmd.Flags().StringP("dir", "", "database/migration/", "directory with migration files")
	migrateCmd.Flags().StringP("table", "", "db", "migrations table name")
	migrateCmd.Flags().BoolP("verbose", "", false, "enable verbose mode")
	migrateCmd.Flags().BoolP("guide", "", false, "print help")

	cmd := []*cobra.Command{
		{
			Use:   "http",
			Short: "Run HTTP Server",
			Run: func(cmd *cobra.Command, args []string) {
				http.Start(ctx)
			},
		},
		{
			Use:   "gen:all",
			Short: "Generator CRUD",
			Run: func(cmd *cobra.Command, args []string) {
				genx.Gen()
			},
		},
		{
			Use:   "gen:logic",
			Short: "Generator logic use case",
			Run: func(cmd *cobra.Command, args []string) {
				genx.GenLogic()
			},
		},
		{
			Use:   "gen:entity",
			Short: "Generator struct entity from table",
			Run: func(cmd *cobra.Command, args []string) {
				genx.GenEntity()
			},
		},
		{
			Use:   "gen:io",
			Short: "Generator struct presentation from table",
			Run: func(cmd *cobra.Command, args []string) {
				genx.GenPresentation()
			},
		},
		{
			Use:   "gen:all-table",
			Short: "Generator crud all table",
			Run: func(cmd *cobra.Command, args []string) {
				genx.GenerateAll()
			},
		},
		{
			Use:   "gen:all-entity",
			Short: "Generator all table to entity",
			Run: func(cmd *cobra.Command, args []string) {
				genx.GenerateAllEntity()
			},
		},
		{
			Use:   "gen:all-presentation",
			Short: "Generator all table to presentation",
			Run: func(cmd *cobra.Command, args []string) {
				genx.GenerateAllPresentation()
			},
		},
		migrateCmd,
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
