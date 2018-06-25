package cmd

import (
	"errors"
	"github.com/mono83/oscar/out"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var reportCmd = &cobra.Command{
	Use:   "report file",
	Short: "Converts report from JSON format to other",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("report file name not provided")
		}

		bts, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		report, err := out.LoadJSON(bts)
		if err != nil {
			return err
		}

		out.PrintSummary(os.Stdout, report)
		return nil
	},
}
