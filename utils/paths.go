/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool.
package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Source - https://stackoverflow.com/a/10510783
// Posted by Mostafa, modified by community. See post 'Timeline' for change history
// Retrieved 2026-08-10, License - CC BY-SA 4.0

// DirExists checks whether a given directory exists.
func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// MkDirIfNotExist creates given directory if it does not already exist.
// If the directory already exists, it does not return an error.
func MkDirIfNotExist(path string) error {
	info, err := os.Stat(path)

	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path %s exists but is not a directory", path)
		}
		return nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	return os.MkdirAll(path, os.ModePerm)
}

// RmDir removes a given directory.
// If the directory does not exist, it returns no error
func RmDir(path string) error {
	info, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path %s exists but is not a directory, use RmFile instead", path)
	}

	return os.RemoveAll(path) // recursive
}

// RmFile removes a given file.
// If the file does not exist, it returns no error
func RmFile(path string) error {
	info, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("path %s is a directory, use RmDir instead", path)
	}

	return os.Remove(path)
}

// RmAllExcept removes all elements (files and directories) recursively from a given path, except for a specified keep
// list.
func RmAllExcept(origin string, keep ...string) error {
	entries, err := os.ReadDir(origin)
	if err != nil {
		return err
	}

	keepSet := make(map[string]bool)
	for _, name := range keep {
		keepSet[name] = true
	}

	for _, entry := range entries {
		if keepSet[entry.Name()] {
			continue
		}

		fullPath := filepath.Join(origin, entry.Name())

		if entry.IsDir() {
			if err := RmDir(fullPath); err != nil {
				return err
			}
		} else {
			if err := RmFile(fullPath); err != nil {
				return err
			}
		}
	}

	return nil
}
