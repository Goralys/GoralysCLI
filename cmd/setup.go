/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package cmd

import (
	_ "embed"
	"fmt"
	. "goralys-cli/utils"
	utils "goralys-cli/utils/templates"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed templates/dynamic/env.yaml
var envTemplate string

//go:embed templates/dynamic/env.local.yaml
var envLocalTemplate string

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "This command is used to setup Goralys.",
	Long: `This command is used to create the necessary files and install the
    dependencies (pnpm and composer) for the project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if mobileFlag {
			name = "GoralysCap"
		} else {
			name = "Goralys"
		}

		Logf("Setting up %s...", name)
		Log("Finding repo root...")

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get wd, %s", err)
		}

		root, err := FindRepoRoot(cwd, mobileFlag)
		if err != nil {
			return fmt.Errorf("setup failed, %s", err)
		}

		Logf("Found, setting up for repo at %s", root)

		stop := StartSpinner("Checking for pnpm")
		pnpm, err := ResolvePnpm("install")
		if err != nil {
			stop(false)
			return fmt.Errorf("failed to find pnpm, %s", err)
		}
		stop(true)

		pnpm.Stdout = NewPrefixWriter("pnpm", os.Stdout)
		pnpm.Stderr = NewPrefixWriter("pnpm", os.Stderr)
		if err = pnpm.Run(); err != nil {
			return fmt.Errorf("failed to run pnpm, %s", err)
		}

		Log("pnpm dependencies installed")

		if !mobileFlag {
			stop := StartSpinner("Checking for composer ")
			php, err := ResolvePhp()
			if err != nil {
				stop(false)
				return fmt.Errorf("failed to find php, %s", err)
			}
			stop(true)

			composer, err := ResolveComposer(php, "install", "--working-dir=backend", "--ansi")
			if err != nil {
				return fmt.Errorf("failed to find composer, %s", err)
			}

			composer.Stdout = NewPrefixWriter("composer", os.Stdout)
			composer.Stderr = NewPrefixWriter("composer", os.Stderr)
			if err := composer.Run(); err != nil {
				return fmt.Errorf("failed to run composer, %s", err)
			}

			Log("Composer dependencies installed.")

			err = CreateBackendDirs(root)
			if err != nil {
				return err
			}

			err = RemoveNonBackendDirs(root)
			if err != nil {
				return err
			}
		}

		Log("Configuring environments")

		if !mobileFlag {
			Log("(1/2) Creating .env")
			err = utils.MakeEnvFileFromTemplate(root, envTemplate)
			if err != nil {
				return err
			}

			Log("(2/2) Creating .env.local")
			err = utils.MakeEnvFileFromTemplate(root, envLocalTemplate)
			if err != nil {
				return err
			}
		} else {
			Log("Creating .env.local")
			err = utils.MakeEnvFileFromTemplate(root, envLocalTemplate)
			if err != nil {
				return err
			}
		}

		var testsStr = "eslint + phpcs"
		if mobileFlag {
			testsStr = "eslint"
		}

		var tests string
		Promptf(&tests, "Do you want the setup to run checks (%s) ? [Y/n]: ", testsStr)

		if strings.ToLower(tests) == "y" {
			stop = StartSpinner("Running phpcs")
			err = RunPhpCS()
			if err != nil {
				stop(false)

				var reRun string
				Prompt(&reRun, "phpcs violations were found, do you want setup to try to fix them ? [Y/n]")

				if strings.ToLower(reRun) == "y" {
					stop = StartSpinner("Running phpcbf")
					if err = RunPhpCBF(); err != nil {
						stop(false)
						return err
					}
					stop(true)

					stop = StartSpinner("Re-running phpcs after fixes")
					if err = RunPhpCS(); err != nil {
						stop(false)
						return err
					}
					stop(true)
					Log("phpcs clean after fixes")
				}
			} else {
				stop(true)
			}

			stop = StartSpinner("Running eslint")
			if err = RunEslint(); err != nil {
				stop(false)
				return err
			}
			stop(true)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
