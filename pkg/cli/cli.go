package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Cli struct {
	RootCmd *cobra.Command
}

func (cli Cli) Execute() {
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func New() Cli {
	rootCmd := &cobra.Command{
		Use: "jgo",
		// Short: "JGO by JEOGA INC.",
		Long: `
JGO is a Fast and Flexible Go Boilerplate built with love @ JEOGA Inc.
Complete documentation is available at http://jgo.jeoga.io/`,
	}

	cli := Cli{
		RootCmd: rootCmd,
	}
	return cli
}
