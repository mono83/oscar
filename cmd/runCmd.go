package cmd

import (
	"errors"
	"github.com/mono83/oscar"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"regexp"
)

var verbose bool
var quiet bool
var environmentFile string
var filter string
var header string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs tests from lua file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("one lua file must be supplied")
		}

		values := map[string]string{}
		if len(environmentFile) > 0 {
			// Reading environment file
			bts, err := ioutil.ReadFile(environmentFile)
			if err != nil {
				return err
			}

			i, err := ini.Load(bts)
			if err != nil {
				return err
			}

			sec, err := i.GetSection("")
			if err != nil {
				return err
			}

			values = sec.KeysHash()
		}

		o := &oscar.Oscar{
			Debug: verbose,
			Vars:  values,
		}

		if !quiet {
			o.Output = os.Stdout
		}
		if len(header) > 0 {
			o.Include = []string{header}
		}

		if len(filter) > 0 {
			nameMatcher, err := regexp.Compile("(?i)" + filter)
			if err != nil {
				return err
			}
			o.CaseSelector = func(testCase *oscar.TestCase) bool {
				return nameMatcher.MatchString(testCase.Name)
			}
		}

		return o.StartFile(args[0])
	},
}

func init() {
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose (debug) mode")
	runCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress any output")
	runCmd.Flags().StringVarP(&environmentFile, "env", "e", "", "Root variables, passed to Oscar")
	runCmd.Flags().StringVarP(&filter, "filter", "f", "", "Test case name filter, regex")
	runCmd.Flags().StringVarP(&header, "lib", "l", "", "Add library lua file with helper functions")
}
