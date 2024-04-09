package commands

import (
	"fmt"
	"github.com/Pactus-Contrib/Indexer/client"
	"github.com/Pactus-Contrib/Indexer/config"
	"github.com/Pactus-Contrib/Indexer/logging"
	"github.com/Pactus-Contrib/Indexer/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run indexer",
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
				Debug:        cfg.Logging.Debug,
				EnableCaller: cfg.Logging.EnableCaller,
				SkipCaller:   3,
			}

			if len(cfg.Logging.SentryDSN) != 0 {
				logOpt.Sentry = &logging.SentryConfig{
					DSN:              cfg.Logging.SentryDSN,
					AttachStacktrace: true,
					ServerName:       version.Application,
					Environment:      logging.PRODUCTION,
					Release:          version.Semantic(),
					Dist:             version.Full(),
					EnableTracing:    true,
					Debug:            true,
					TracesSampleRate: 1.0,
				}
			}

			logger, err = logging.New(cfg.Logging.Handler, logOpt)
		}
		if err != nil {
			return err
		}

		fmt.Println(logger)

		pactus, err := client.NewPactus(cmd.Context(), cfg.Pactus.RPC)
		if err != nil {
			return err
		}

		fmt.Println(pactus)

		return nil
	},
}

func defaultLogging() (logging.Logger, error) {
	return logging.New(logging.ConsoleHandler, logging.Options{
		Development:  false,
		Debug:        false,
		EnableCaller: true,
		SkipCaller:   3,
	})
}
