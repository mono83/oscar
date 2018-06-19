package cmd

import (
	"errors"
	"github.com/fatih/color"
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/out"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
)

var verbose bool
var veryVerbose bool
var quiet bool
var noAnsi bool
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
		d := &out.Dispatcher{}
		if !quiet && (verbose || veryVerbose) {
			d.List = append(d.List, out.GetTracer(os.Stdout))
		}
		if !quiet {
			d.List = append(
				d.List,
				out.GetAftermath(os.Stdout),
				out.GetTestCasePrinter(os.Stdout, veryVerbose),
			)
		}
		context.OnEvent = d.OnEvent
		color.NoColor = noAnsi

		// Building Oscar runner
		o := &oscar.Oscar{
			TestContext: context,
		}
		for _, luaFile := range args {
			if err := o.AddTestSuiteFile(luaFile, header, filter); err != nil {
				return err
			}
		}

		// Starting
		return o.Start()
	},
}

func init() {
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose (debug) mode")
	runCmd.Flags().BoolVar(&veryVerbose, "vv", false, "Even more verbose mode")
	runCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress any output")
	runCmd.Flags().BoolVar(&noAnsi, "no-ansi", false, "Disable ANSI color output")
	runCmd.Flags().StringVarP(&environmentFile, "env", "e", "", "Root variables, passed to TestSuite")
	runCmd.Flags().StringVarP(&filter, "filter", "f", "", "Test case name filter, regex")
	runCmd.Flags().StringVarP(&header, "lib", "l", "", "Add library lua file with helper functions")
}
