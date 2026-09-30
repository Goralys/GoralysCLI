/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"goralys-cli/shared"
)

// RunEslint runs the pnpm lint command.
func RunEslint() error {
	pnpm, err := ResolvePnpm("run", "lint")

	if err != nil {
		return fmt.Errorf("failed to run eslint, %w", err)
	}

	if err = pnpm.Run(); err != nil {
		return fmt.Errorf("an error occurred while running eslint, %w", err)
	}

	return nil
}

// RunPhpCS runs the PHP Code Sniffer and Beautifier tool. The actual executable invoked (phpcs/phpcbf) depends on the
// provided command.
func RunPhpCS(backendFlag bool, command string) error {
	var cmd *exec.Cmd
	var err error
	var backendDir = filepath.Join(shared.RepoRoot, "backend")

	if backendFlag {
		php, errPhp := ResolvePhp()

		if errPhp == nil {
			cmd = RunPhp(
				php,
				filepath.Join(backendDir, "vendor", "bin", command),
				".",
				"--standard="+filepath.Join(backendDir, "phpcs.xml"),
			)
			cmd.Dir = backendDir
		} else {
			return fmt.Errorf("failed to resolve php, %w", errPhp)
		}
	} else {
		cmd, err = ResolvePnpm("run", command)
		if err != nil {
			return fmt.Errorf("failed to run %s, %w", command, err)
		}
	}

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("an error occurred while running %s, %w", command, err)
	}

	return nil
}

// phpCsTest is a simple test runner that checks for phpcs code violations. If it finds violations, it can run phpcbf
// to fix what it can.
var phpCsTest = shared.TestRunner{
	Name: "phpcs",
	Callback: func() error {
		stop := StartSpinner("Running phpcs")
		err := RunPhpCS(shared.BackendFlag, "phpcs")
		if err != nil {
			Logf("(phpcs error): %s", err)
			stop(false)

			var reRun bool
			PromptBool(&reRun, "phpcs violations were found, do you want setup to try to fix them ?")

			if reRun {
				stop = StartSpinner("Running phpcbf")
				if err = RunPhpCS(shared.BackendFlag, "phpcbf"); err != nil {
					stop(false)
					return err
				}
				stop(true)

				stop = StartSpinner("Re-running phpcs after fixes")
				if err = RunPhpCS(shared.BackendFlag, "phpcs"); err != nil {
					stop(false)
					return err
				}
				stop(true)
				Log("phpcs clean after fixes")
			}
		} else {
			stop(true)
		}
		return nil
	},
}

// eslintTest is a simple test runner that runs eslint.
var eslintTest = shared.TestRunner{
	Name: "eslint",
	Callback: func() error {
		stop := StartSpinner("Running eslint")
		if err := RunEslint(); err != nil {
			stop(false)
			return err
		}
		stop(true)
		return nil
	},
}

// Tests is a list of all the tests that can be run.
var Tests = []shared.TestRunner{phpCsTest, eslintTest}
