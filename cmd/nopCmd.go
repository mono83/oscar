package cmd

import "github.com/spf13/cobra"

var nopCmd = &cobra.Command{Hidden: true, Use: "nop", Run: func(*cobra.Command, []string) {}}
