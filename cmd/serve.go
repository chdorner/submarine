package cmd

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/router"
)

func NewServeCmd() *cobra.Command {
	var db *gorm.DB
	var addr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start webserver",
		PreRun: func(cmd *cobra.Command, args []string) {
			ctx := createContext(cmd.Flags())
			configureLogging(ctx)

			db = initDBConn(cmd.Flags())
			addr, _ = cmd.Flags().GetString("addr")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := data.NewSettingsRepository(db)
			if !repo.IsInitialized() {
				return errors.New("submarine is not initialized yet, run `submarine init` first")
			}

			e := router.New(db)
			logrus.WithField("addr", addr).Info("starting submarine")
			return e.Start(addr)
		},
	}

	fl := cmd.Flags()
	fl.StringP("addr", "a", "127.0.0.1:9876", "listen address")

	return cmd
}
