/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool
package utils

import (
	"fmt"
	"path"
)

func mkDirsFromRoot(root string, first string, _d ...string) error {
	for _, d := range append([]string{first}, _d...) {
		if err := MkDirIfNotExist(path.Join(root, d)); err != nil {
			return fmt.Errorf("failed to create %s, %w", d, err)
		}
	}

	return nil
}

// CreateBackendDirs create all necessary directories for Goralys' backend
func CreateBackendDirs(root string) error {
	stop := StartSpinner("Creating Logs directory")
	if err := mkDirsFromRoot(root, "backend/Logs"); err != nil {
		stop(false)
		return err
	}
	stop(true)

	stop = StartSpinner("Creating RateLimiter directory")
	if err := mkDirsFromRoot(root, "backend/RateLimiter"); err != nil {
		stop(false)
		return err
	}
	stop(true)

	stop = StartSpinner("Creating Assets directory")
	if err := mkDirsFromRoot(root,
		"backend/Assets",
		"backend/Assets/Template",
		"backend/Assets/Mails",
		"backend/Assets/Template/Exports",
		"backend/Assets/StudentsDrafts",
	); err != nil {
		stop(false)
		return err
	}
	stop(true)

	return nil
}

// RemoveNonBackendDirs removes all directories that are unnecessary to the backend (e.g. frontend directories).
// It is only used when deploying the backend on a server.
func RemoveNonBackendDirs(root string) error {
	stop := StartSpinner("Removing non backend directories")
	if err := RmAllExcept(root,
		"LICENSE",
		"README.md",
		"CONTRIBUTING.md",
		"backend",
		".git",
		"scripts",
	); err != nil {
		stop(false)
		return err
	}

	stop(true)
	return nil
}
