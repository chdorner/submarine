package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/chdorner/submarine/data"
)

func NewInitCmd() *cobra.Command {
	var db *gorm.DB
	var password string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize or update submarine settings",
		PreRun: func(cmd *cobra.Command, args []string) {
			ctx := createContext(cmd.Flags())
			configureLogging(ctx)

			db = initDBConn(cmd.Flags(), true)
			password, _ = cmd.Flags().GetString("password")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := data.NewSettingsRepository(db)
			settings := data.SettingsUpsert{
				Password: password,
			}
			err := repo.Upsert(settings)
			if err != nil {
				return err
			}

			logrus.Info("Successfully initialized submarine")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("password", "p", "", "set password")

	return cmd
}
