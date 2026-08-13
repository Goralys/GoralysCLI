/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool
package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"goralys-cli/utils"
)

// LoadStaticTemplate loads a static template by copying its content into the given path. It also creates the directory
// of the destination file if it does not exist.
func LoadStaticTemplate(template string, path string) error {
	info, err := os.Stat(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("cannot write template to path %s because it is a directory", path)
	}

	err = utils.MkDirIfNotExist(filepath.Dir(path))
	if err != nil {
		return err
	}
	err = os.WriteFile(path, []byte(template), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
