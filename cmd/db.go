package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chdorner/submarine/data"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func NewDBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database management",
	}
	cmd.AddCommand(NewDBMigrateCmd())
	cmd.AddCommand(NewDBRollbackCmd())
	return cmd
}

func NewDBMigrateCmd() *cobra.Command {
	var db *gorm.DB

	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database with newest changes",
		PreRun: func(cmd *cobra.Command, args []string) {
			ctx := createContext(cmd.Flags())
			configureLogging(ctx)

			db = initDBConn(cmd.Flags(), false)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("This will migrate your database and could lead to data-loss.")
			fmt.Println("Please create database backup first.")
			if confirm("Are you sure you want to continue?", 3) {
				err := data.NewMigrator(db).Migrate()
				if err != nil {
					log.Fatalf("failed to migrate database: %s", err)
					return err
				}
				fmt.Println("Successfully migrated database")
			}
			return nil
		},
	}
}

func NewDBRollbackCmd() *cobra.Command {
	var db *gorm.DB

	return &cobra.Command{
		Use:    "rollback",
		Short:  "Rollback latest migration",
		Hidden: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			ctx := createContext(cmd.Flags())
			configureLogging(ctx)

			db = initDBConn(cmd.Flags(), false)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("This will rollback the latest migration in your database and does lead to data-loss.")
			fmt.Println("Please create database backup first.")
			if confirm("Are you sure you want to continue?", 3) {
				err := data.NewMigrator(db).RollbackLast()
				if err == nil {
					fmt.Println("Successfully rolled back database")
					return nil
				}
				if err == gormigrate.ErrNoRunMigration {
					fmt.Println("Already on an empty database, nothing to do.")
					return nil
				}

				fmt.Printf("failed to rollback database: %s\n", err)
				return err
			}
			return nil
		},
	}
}

func confirm(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}
