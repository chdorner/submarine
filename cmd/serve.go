package cmd

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/chdorner/submarine/data"
	"github.com/chdorner/submarine/web"
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

			mux := web.NewRouter()
			logrus.WithField("addr", addr).Info("starting submarine")
			return http.ListenAndServe(addr, mux)
		},
	}

	fl := cmd.Flags()
	fl.StringP("addr", "a", "127.0.0.1:8080", "listen address")

	return cmd
}
