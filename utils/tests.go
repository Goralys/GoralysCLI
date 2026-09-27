/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package utils

import (
	"fmt"
	"os/exec"
)

// TestRunner represents a simple test that has a name and can be run via a callback.
type TestRunner struct {
	Name     string
	Callback func() error
}

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

// RunPhpCS runs the composer phpcs command.
func RunPhpCS(backendFlag bool) error {
	var bin *exec.Cmd
	var err error

	if backendFlag {
		php, errPhp := ResolvePhp()

		if errPhp == nil {
			bin, err = ResolveComposer(php, "--working-dir=backend", "phpcs")
		} else {
			return errPhp
		}
	} else {
		bin, err = ResolvePnpm("run", "phpcs")
	}

	if err != nil {
		return fmt.Errorf("failed to run phpcs, %w", err)
	}

	if err = bin.Run(); err != nil {
		return fmt.Errorf("an error occurred while running phpcs, %w", err)
	}

	return nil
}

// RunPhpCBF runs the composer phpcbf command.
func RunPhpCBF(backendFlag bool) error {
	var bin *exec.Cmd
	var err error

	if backendFlag {
		php, errPhp := ResolvePhp()

		if errPhp == nil {
			bin, err = ResolveComposer(php, "--working-dir=backend", "phpcbf")
		} else {
			return errPhp
		}
	} else {
		bin, err = ResolvePnpm("run", "phpcbf")
	}

	if err != nil {
		return fmt.Errorf("failed to run phpcbf, %w", err)
	}

	if err = bin.Run(); err != nil {
		return fmt.Errorf("an error occurred while running phpcbf, %w", err)
	}

	return nil
}
