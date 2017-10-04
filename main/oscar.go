package main

import (
	"fmt"
	"github.com/mono83/oscar/cmd"
	"os"
)

func main() {
	command := cmd.OscarCmd
	command.SilenceUsage = true
	command.SilenceErrors = true
	if err := command.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
