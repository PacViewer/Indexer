package commands

import (
	"fmt"
	"github.com/Pactus-Contrib/Indexer/version"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "indexer information",
	RunE: func(cmd *cobra.Command, args []string) error {
		build := new(strings.Builder)
		if _, err := build.WriteString(fmt.Sprintf("Application: %s\n", version.Application)); err != nil {
			return err
		}

		if _, err := build.WriteString(fmt.Sprintf("Description: %s\n", version.Description)); err != nil {
			return err
		}

		if _, err := build.WriteString(fmt.Sprintf("Version: %s\n", version.Semantic())); err != nil {
			return err
		}

		if _, err := build.WriteString(fmt.Sprintf("Commit ID: %s\n", version.CommitID)); err != nil {
			return err
		}

		if _, err := build.WriteString(fmt.Sprintf("Build Time: %s\n", version.BuildTime)); err != nil {
			return err
		}

		fmt.Println(build.String())

		return nil
	},
}
