package cmd

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gorm.io/gorm"

	"github.com/chdorner/submarine/data"
)

var rootCmd = &cobra.Command{
	Use:   "submarine",
	Short: "Submarine is a single-user bookmarks manager",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	rootCmd.PersistentFlags().StringP("db", "d", "submarine.db", "path to sqlite database")

	rootCmd.AddCommand(NewServeCmd())
	rootCmd.AddCommand(NewDBCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewVersionCommand())
}

func createContext(flags *pflag.FlagSet) context.Context {
	debug, _ := flags.GetBool("debug")
	return context.WithValue(context.Background(), "debug", debug) //nolint:staticcheck
}

func configureLogging(ctx context.Context) {
	if ctx.Value("debug").(bool) {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode turned on")
	}
}

func initDBConn(flags *pflag.FlagSet, migrate bool) *gorm.DB {
	path, _ := flags.GetString("db")
	db, err := data.Connect(path)
	if err != nil {
		panic(err)
	}

	if migrate {
		data.Migrate(db)
	}

	return db
}
