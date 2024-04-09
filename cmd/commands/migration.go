package commands

import (
	"github.com/Pactus-Contrib/Indexer/config"
	"github.com/Pactus-Contrib/Indexer/db"
	"github.com/Pactus-Contrib/Indexer/logging"
	"github.com/Pactus-Contrib/Indexer/schema"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrationCmd)
}

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New(configPath)
		if err != nil {
			return err
		}

		if err := cfg.Validate(); err != nil {
			return err
		}

		logger, err := defaultLogging()
		if cfg.Logging != nil {
			logOpt := logging.Options{
				Development:  false,
				Debug:        false,
				EnableCaller: false,
				SkipCaller:   0,
			}

			logger, err = logging.New(cfg.Logging.Handler, logOpt)
		}
		if err != nil {
			return err
		}

		logger.InfoContext(cmd.Context(), false, "Migration started")
		p := db.NewPool(logger)
		logger.InfoContext(cmd.Context(), false, "Created new database pool")

		for _, d := range cfg.DBS {
			switch d.Engine {
			case schema.MARIADB, schema.MYSQL, schema.POSTGRESQL:
				sql, err := db.NewSQL(d)
				if err != nil {
					return err
				}
				p.RegisterEngine(sql)
			case schema.MONGODB:
				mgo, err := db.NewMongodb(cmd.Context(), d)
				if err != nil {
					return err
				}
				p.RegisterEngine(mgo)
			}
		}

		logger.InfoContext(cmd.Context(), false, "All database registered in pool")
		logger.InfoContext(cmd.Context(), false, "Migration in process...")
		if err := p.Migration(cmd.Context(), cfg.IndexerUuid, cfg.LastBlockHeight); err != nil {
			return err
		}
		logger.InfoContext(cmd.Context(), false, "Migration completed")

		return nil
	},
}
