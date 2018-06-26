package cmd

import (
	"errors"
	"github.com/mono83/oscar/out"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var htmlCmd = &cobra.Command{
	Use:   "html file.json folder",
	Short: "Converts report from JSON format to HTML",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("both source report and destination folder must be provided")
		}

		bts, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		report, err := out.LoadJSON(bts)
		if err != nil {
			return err
		}

		return out.WriteHTMLFiles(args[1], report)
	},
}
