package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mono83/oscar"
	"github.com/mono83/oscar/lua"
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
var outJSONFile string

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

		// Vars
		color.NoColor = noAnsi

		// Building ROOT testing context
		context := oscar.NewContext()
		context.SetInitial(values)
		context.Set("lua.engine", "Oscar ][")

		// Adding event dispatcher
		d := &out.Dispatcher{}
		context.OnEvent = d.OnEvent

		// Registering event listeners (logging and etc)
		reporter := &out.Report{}
		d.List = append(d.List, reporter.OnEvent)

		if !quiet {
			defer func() {
				fmt.Fprintln(os.Stdout, "")
				fmt.Fprintln(os.Stdout, "")
				out.PrintTestCaseErrorsSummary(os.Stdout, reporter)
				fmt.Fprintln(os.Stdout, "")
				fmt.Fprintln(os.Stdout, "")
				out.PrintSummary(os.Stdout, reporter)
				fmt.Fprintln(os.Stdout, "")
			}()
		}

		if len(outJSONFile) > 0 {
			defer func() {
				ioutil.WriteFile(outJSONFile, []byte(reporter.JSON()), 0644)
			}()
		}

		if !quiet {
			if verbose || veryVerbose {
				d.List = append(d.List, out.FullRealTimePrinter(os.Stdout, veryVerbose, veryVerbose))
			} else {
				d.List = append(d.List, out.DotRealTimePrinter(os.Stdout))
			}
		}

		// Loading LUA files
		var suites []oscar.Suite
		for _, file := range args {
			var suite oscar.Suite
			var err error
			if len(header) > 0 {
				suite, err = lua.SuiteFromFiles(header, file)
			} else {
				suite, err = lua.SuiteFromFiles(file)
			}

			if err != nil {
				return err
			}
			suites = append(suites, suite)
		}

		// Running
		return oscar.RunSequential(context, suites)
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
	runCmd.Flags().StringVarP(&outJSONFile, "json-report", "j", "", "JSON report filename")
}
