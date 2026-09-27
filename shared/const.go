/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package shared is the package that holds all shared information such as flags and constants.
package shared

import (
	"goralys-cli/utils"
)

var PHPCS_TEST = utils.TestRunner{
	Name: "phpcs",
	Callback: func() error {
		stop := utils.StartSpinner("Running phpcs")
		err := utils.RunPhpCS(BackendFlag)
		if err != nil {
			stop(false)

			var reRun bool
			utils.Logf("Phpcs error: %s", err)
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

var ESLINT_TEST = utils.TestRunner{
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
