/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package cmd

import (
	_ "embed"
	"fmt"
	. "goralys-cli/utils"
	"os"

	"github.com/spf13/cobra"
)

//go:embed banner.txt
var banner string

func showBanner() {
	fmt.Println(Colorize(ColorGoralys, banner))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goralys-cli",
	Short: "The CLI tool for Goralys",
	Long: `GoralysCLI is the tool that powers Goralys' setup and backup	system.
    This CLI is separate from the monorepo at https://github.com/SAMSAM-55/Goralys.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		showBanner()

		if backendFlag && mobileFlag {
			return fmt.Errorf("cannot pass both --backend-only and --mobile")
		}

		if yesFlag && noFlag {
			return fmt.Errorf("cannot pass both --yes and --no")
		}

		if yesFlag || noFlag {
			InitSkipPrompt(yesFlag)
		}

		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var mobileFlag bool
var backendFlag bool
var yesFlag bool
var noFlag bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&mobileFlag,
		"mobile",
		"m",
		false,
		"Whether the cli tool is ran for the mobile app or for the main mono repo",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&backendFlag,
		"backend-only",
		"bo",
		false,
		"Whether the cli tool is ran for the backend only or for the main mono repo",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&yesFlag,
		"yes",
		"y",
		false,
		"The interactive prompts will be skipped and will be answered as 'yes'. If both this flag and the 'no'"+
			"flag are passed, the prompts will be answered as 'yes'",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&noFlag,
		"no",
		"n",
		false,
		"The interactive prompts will be skipped and will be answered as 'no'",
	)
}
