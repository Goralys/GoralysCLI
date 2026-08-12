/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */
package utils

import "fmt"

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

func RunPhpCS() error {
	pnpm, err := ResolvePnpm("run", "phpcs")

	if err != nil {
		return fmt.Errorf("failed to run phpcs, %w", err)
	}

	if err = pnpm.Run(); err != nil {
		return fmt.Errorf("an error occurred while running phpcs, %w", err)
	}

	return nil
}

func RunPhpCBF() error {
	pnpm, err := ResolvePnpm("run", "phpcbf")

	if err != nil {
		return fmt.Errorf("failed to run phpcbf, %w", err)
	}

	if err = pnpm.Run(); err != nil {
		return fmt.Errorf("an error occurred while running phpcbf, %w", err)
	}

	return nil
}
