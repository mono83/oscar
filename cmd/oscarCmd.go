package cmd

import "github.com/spf13/cobra"

// OscarCmd is main command line command
var OscarCmd = &cobra.Command{
	Use: "oscar",
}

func init() {
	OscarCmd.AddCommand(
		runCmd,
	)
}
