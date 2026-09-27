/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package cmd

import (
	"fmt"
	"goralys-cli/shared"
	"os"
	"slices"
	"strings"

	"goralys-cli/utils"

	"github.com/spf13/cobra"
)

var eslintFlag bool
var phpcsFlag bool

// testsCmd represents the tests command
var testsCmd = &cobra.Command{
	Use:   "tests",
	Short: "Runs the tests",
	Long: `This command runs the tests for the local repository.
		   By default, it runs both eslint and phpcs for a full Goralys repository (mono-repo).
		   For GoralysCap, it only runs eslint and for a backend only setup, it only runs phpcs.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		var name = "Goralys"
		if shared.MobileFlag {
			name = "GoralysCap"
		}
		if shared.BackendFlag {
			name = "Goralys [backend]"
		}

		utils.Logf("Running test for %s", name)

		stop := utils.StartSpinner("Finding repo root")

		cwd, err := os.Getwd()
		if err != nil {
			stop(false)
			return fmt.Errorf("failed to get wd, %s", err)
		}

		root, err := utils.FindRepoRoot(cwd, shared.MobileFlag)
		if err != nil {
			stop(false)
			return fmt.Errorf("setup failed, %s", err)
		}

		var tests []string
		if shared.MobileFlag || !shared.BackendFlag {
			tests = append(tests, "eslint")
		}
		if shared.BackendFlag || !shared.MobileFlag {
			tests = append(tests, "phpcs")
		}

		stop(true)
		utils.Logf("Found, running tests (%s) for repo at %s", strings.Join(tests, " + "), root)

		for _, t := range []utils.TestRunner{shared.PHPCS_TEST, shared.ESLINT_TEST} {
			if slices.Contains(tests, t.Name) {
				err = t.Callback()
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testsCmd)

	testsCmd.Flags().BoolVar(&eslintFlag, "eslint", false, "Explicit flag for running eslint")
	testsCmd.Flags().BoolVar(&phpcsFlag, "phpcs", false, "Explicit flag for running phpcs")
}
