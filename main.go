package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/DmiTryAgain/secretSanta/pkg/app"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := newRootCmd()
	rootCmd.InitDefaultVersionFlag()
	rootCmd.InitDefaultHelpFlag()
	exitOnErr(rootCmd.ParseFlags(os.Args))

	// create app and run
	a := app.New(rootCmd)
	exitOnErr(a.Run(context.Background()))
}

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "secretSanta",
		Short:   "",
		Long:    ``,
		Version: appVersion(),
	}
}

func appVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return info.Main.Version
}
