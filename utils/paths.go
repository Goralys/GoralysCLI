/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool.
package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// FromHomeDir returns the user's home path followed by an additional and optional path.
func FromHomeDir(_p ...string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	paths := append([]string{homeDir}, _p...)
	return filepath.Join(paths...), nil
}

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

// Cp copies a source directory into a destination directory recursively.
func Cp(src string, dest string) error {
	src = filepath.Clean(src)
	dest = filepath.Clean(dest)

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, relPath)

		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(target, info.Mode().Perm())
		}

		return copyFile(path, target)
	})
}

// CopyFile copies a given file into a destination file
func copyFile(src string, dest string) error {
	err := MkDirIfNotExist(filepath.Dir(dest))
	if err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			return
		}
	}(in)

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Sync()
}

// CopyFile copies a given file into from its original directory into a target directory
func CopyFile(origin string, target string, name string) error {
	return copyFile(filepath.Join(origin, name), filepath.Join(target, name))
}
