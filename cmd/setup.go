/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"goralys-cli/shared"

	"goralys-cli/utils"
	templates "goralys-cli/utils/templates"

	"github.com/spf13/cobra"
)

//go:embed templates/dynamic/env.yaml
var envTemplate string

//go:embed templates/dynamic/env.next.yaml
var envNextTemplate string

//go:embed templates/dynamic/env.cap.yaml
var envCapTemplate string

//go:embed templates/static/.htaccess
var htAccessTemplate string

//go:embed templates/static/MainActivity.java
var mainActivityTemplate string

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "This command is used to setup Goralys.",
	Long: `This command is used to create the necessary files and install the
    dependencies (pnpm and composer) for the project.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		var name = "Goralys"
		var suffix = "frontend"
		if shared.MobileFlag {
			name = "GoralysCap"
		}
		if shared.BackendFlag {
			suffix = "backend"
			name = "Goralys [backend]"
		}

		utils.Logf("Setting up %s", name)

		backupPath, err := utils.FromHomeDir(GoralysBackupDir...)
		if err != nil {
			return fmt.Errorf("failed to construct the backup path, %w", err)
		}

		backupPath = filepath.Join(backupPath, strings.ToLower(name), strings.ToLower(suffix))
		stop := utils.StartSpinner("Finding repo root")

		cwd, err := os.Getwd()
		if err != nil {
			stop(false)
			return fmt.Errorf("failed to get wd, %s", err)
		}

		shared.RepoRoot, err = utils.FindRepoRoot(cwd, shared.MobileFlag)
		if err != nil {
			stop(false)
			return fmt.Errorf("setup failed, %s", err)
		}

		stop(true)
		utils.Logf("Found, setting up for repo at %s", shared.RepoRoot)

		// backup restore
		stop = utils.StartSpinner("Checking for existing backup at " + backupPath)
		exists, err := utils.DirExists(backupPath)
		if err != nil {
			stop(false)
			return fmt.Errorf("failed to check for existing backup, %s", err)
		}
		stop(true)

		var restore bool
		if exists {
			utils.PromptBool(&restore, "An existing backup was found, do you want to restore it ?")
		}

		if restore {
			utils.Logf("Restoring backup from %s", backupPath)

			stop = utils.StartSpinnerNoPrefix("-> Copying env files")
			if !shared.BackendFlag {
				err = utils.CopyFile(backupPath, shared.RepoRoot, ".env.local")
				if err != nil {
					stop(false)
					return fmt.Errorf("failed to copy .env.local file, %s", err)
				}
			}
			if !shared.MobileFlag {
				err = utils.CopyFile(backupPath, shared.RepoRoot, filepath.Join("backend", ".env"))
				if err != nil {
					stop(false)
					return fmt.Errorf("failed to copy backend/.env file, %s", err)
				}
			}
			stop(true)

			if !shared.MobileFlag {
				stop = utils.StartSpinnerNoPrefix("-> Copying backend/Assets")
				err = utils.Cp(filepath.Join(backupPath, "backend", "Assets"), filepath.Join(shared.RepoRoot, "backend", "Assets"))
				if err != nil {
					stop(false)
					return fmt.Errorf("failed to copy backend/Assets directory, %s", err)
				}
				stop(true)
			}

			utils.Log("Backup successfully restored")
		} else {
			utils.Log("No backup found")
		}

		if !shared.BackendFlag {
			stop := utils.StartSpinner("Checking for pnpm")

			pnpmApprove, errApprove := utils.ResolvePnpm("approve-builds", "--all", "--color")
			pnpm, err := utils.ResolvePnpm("install", "--color")
			if err != nil || errApprove != nil {
				stop(false)
				return fmt.Errorf("failed to find pnpm, %s", err)
			}
			stop(true)

			var approveBuilds bool
			utils.PromptBool(&approveBuilds, "Do you want to approve builds before installing dependencies")
			if approveBuilds {
				pnpmApprove.Stdout = utils.NewPrefixWriter("pnpm", os.Stdout)
				pnpmApprove.Stderr = utils.NewPrefixWriter("pnpm", os.Stderr)
				if err = pnpmApprove.Run(); err != nil {
					return fmt.Errorf("failed to approve builds before installing dependencies, %s", err)
				}
			}

			pnpm.Stdout = utils.NewPrefixWriter("pnpm", os.Stdout)
			pnpm.Stderr = utils.NewPrefixWriter("pnpm", os.Stderr)
			if err = pnpm.Run(); err != nil {
				return fmt.Errorf("failed to run pnpm, %s", err)
			}

			utils.Log("pnpm dependencies installed")
		}

		if !shared.MobileFlag {
			stop := utils.StartSpinner("Checking for composer")
			php, err := utils.ResolvePhp()
			if err != nil {
				stop(false)
				return fmt.Errorf("failed to find php, %s", err)
			}
			stop(true)

			composer, err := utils.ResolveComposer(php, "install", "--working-dir=backend", "--ansi")
			if err != nil {
				return fmt.Errorf("failed to find composer, %s", err)
			}

			composer.Stdout = utils.NewPrefixWriter("composer", os.Stdout)
			composer.Stderr = utils.NewPrefixWriter("composer", os.Stderr)
			if err := composer.Run(); err != nil {
				return fmt.Errorf("failed to run composer, %s", err)
			}

			utils.Log("Composer dependencies installed.")

			err = utils.CreateBackendDirs(shared.RepoRoot)
			if err != nil {
				return err
			}

			err = utils.ElevateToExecutable(filepath.Join(shared.RepoRoot, "backend", "vendor", "bin", "phpcs"))
			if err != nil {
				return err
			}
			err = utils.ElevateToExecutable(filepath.Join(shared.RepoRoot, "backend", "vendor", "bin", "phpcbf"))
			if err != nil {
				return err
			}
		}

		utils.Log("Configuring environments")

		if !shared.MobileFlag {
			err = templates.MakeEnvFileFromTemplate(shared.RepoRoot, envTemplate, "(1/2) Creating .env")
			if err != nil {
				return err
			}

			err = templates.MakeEnvFileFromTemplate(shared.RepoRoot, envNextTemplate, "(2/2) Creating .env.local")
			if err != nil {
				return err
			}
		} else {
			err = templates.MakeEnvFileFromTemplate(shared.RepoRoot, envCapTemplate, "Creating .env.local")
			if err != nil {
				return err
			}
		}

		if shared.BackendFlag || shared.MobileFlag {
			utils.Log("Finalizing your configuration, you are almost there")
			if shared.BackendFlag {
				stop := utils.StartSpinnerNoPrefix("-> Creating .htaccess")
				err := templates.LoadStaticTemplate(htAccessTemplate, filepath.Join(shared.RepoRoot, ".htaccess"))
				if err != nil {
					stop(false)
					return err
				}
				stop(true)
			}

			if shared.MobileFlag {
				stop := utils.StartSpinnerNoPrefix("-> Generating assets")
				pnpmGenerate, err := utils.ResolvePnpm("run", "assets:generate")
				if err != nil {
					stop(false)
					return err
				}

				if err = pnpmGenerate.Run(); err != nil {
					stop(false)
					return err
				}
				stop(true)

				stop = utils.StartSpinnerNoPrefix("-> Building project")
				pnpmBuild, err := utils.ResolvePnpm("run", "build")
				if err != nil {
					stop(false)
					return err
				}

				if err = pnpmBuild.Run(); err != nil {
					stop(false)
					return err
				}
				stop(true)

				androidExists, err := utils.DirExists(filepath.Join(shared.RepoRoot, "android"))
				if err != nil {
					return err
				}

				capCmdTxt := "add"
				if androidExists {
					capCmdTxt = "sync"
				}

				stop = utils.StartSpinnerNoPrefix("-> Setting up capacitor")
				capCmd, err := utils.ResolveNpx("cap", capCmdTxt, "android")
				if err != nil {
					stop(false)
					return err
				}

				if err = capCmd.Run(); err != nil {
					stop(false)
					return err
				}
				stop(true)

				stop = utils.StartSpinnerNoPrefix("-> Creating MainActivity.java")
				err = templates.LoadStaticTemplate(
					mainActivityTemplate,
					filepath.Join(
						shared.RepoRoot,
						"android",
						"app",
						"src",
						"main",
						"java",
						"fr",
						"goralys",
						"app",
						"MainActivity.java",
					),
				)
				if err != nil {
					stop(false)
					return err
				}
				stop(true)
			}
		}

		if shared.BackendFlag {
			utils.Log("Backend only setup detected, removing non backend dir")
			err = utils.RemoveNonBackendDirs(shared.RepoRoot)
			if err != nil {
				return err
			}
		}

		var tests []string
		if shared.MobileFlag || !shared.BackendFlag {
			tests = append(tests, "eslint")
		}
		if shared.BackendFlag || !shared.MobileFlag {
			tests = append(tests, "phpcs")
		}

		var runTests bool
		utils.PromptfBool(&runTests, "Do you want the setup to run checks (%s) ?", strings.Join(tests, " + "))

		if runTests {
			for _, t := range utils.Tests {
				if slices.Contains(tests, t.Name) {
					err = t.Callback()
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
