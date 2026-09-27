/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package shared

import (
	"goralys-cli/utils"
)

// PhpCsTest is a simple test runner which checks for phpcs code violations. If it finds violations, it can run phpcbf
// to fix what it can.
var PhpCsTest = utils.TestRunner{
	Name: "phpcs",
	Callback: func() error {
		stop := utils.StartSpinner("Running phpcs")
		err := utils.RunPhpCS(BackendFlag)
		if err != nil {
			stop(false)

			var reRun bool
			utils.PromptBool(&reRun, "phpcs violations were found, do you want setup to try to fix them ?")

			if reRun {
				stop = utils.StartSpinner("Running phpcbf")
				if err = utils.RunPhpCBF(BackendFlag); err != nil {
					stop(false)
					return err
				}
				stop(true)

				stop = utils.StartSpinner("Re-running phpcs after fixes")
				if err = utils.RunPhpCS(BackendFlag); err != nil {
					stop(false)
					return err
				}
				stop(true)
				utils.Log("phpcs clean after fixes")
			}
		} else {
			stop(true)
		}
		return nil
	},
}

// EslintTest is a simple test runner which runs eslint.
var EslintTest = utils.TestRunner{
	Name: "eslint",
	Callback: func() error {
		stop := utils.StartSpinner("Running eslint")
		if err := utils.RunEslint(); err != nil {
			stop(false)
			return err
		}
		stop(true)
		return nil
	},
}
