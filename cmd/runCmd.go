package cmd

import (
	"errors"
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/out"
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
		if len(args) < 1 {
			return errors.New("at least one lua file must be supplied")
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

		// Building testing context
		context := &oscar.TestContext{Vars: values}

		// Building and configuring Oscars
		oscars := []*oscar.TestSuite{}
		for range args {
			d := &out.Dispatcher{}
			o := &oscar.TestSuite{
				TestContext: &oscar.TestContext{
					Parent: context,
				},
			}
			o.OnEvent = d.OnEmit

			if !quiet && verbose {
				d.List = append(d.List, out.GetTracer(os.Stdout))
			}
			if !quiet {
				d.List = append(
					d.List,
					out.GetAftermath(os.Stdout),
					out.GetTestCasePrinter(os.Stdout),
				)
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

			oscars = append(oscars, o)
		}

		// Running all oscars
		for i, luaFile := range args {
			oscars[i].StartFile(luaFile)
		}

		return nil
	},
}

func init() {
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose (debug) mode")
	runCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress any output")
	runCmd.Flags().StringVarP(&environmentFile, "env", "e", "", "Root variables, passed to TestSuite")
	runCmd.Flags().StringVarP(&filter, "filter", "f", "", "Test case name filter, regex")
	runCmd.Flags().StringVarP(&header, "lib", "l", "", "Add library lua file with helper functions")
}
