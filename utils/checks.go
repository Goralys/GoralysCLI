/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var errMonoRepoRootNotFound = errors.New("goralys mono repo root not found")

const goralysPkgName = "goralys"
const goralysCapPkgName = "goralys-cap"
const goralysComposerName = "goralys/goralys"

type pkgJSON struct {
	Name string `json:"name"`
}

type composerJSON struct {
	Name string `json:"name"`
}

func readJSONField(path string, extract func([]byte) (string, error)) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return extract(data)
}

func isGoralysRoot(dir string) bool {
	pkgName, err := readJSONField(path.Join(dir, "package.json"), func(b []byte) (string, error) {
		var p pkgJSON
		if err := json.Unmarshal(b, &p); err != nil {
			return "", err
		}

		return p.Name, nil
	})
	if err != nil || pkgName != goralysPkgName {
		return false
	}

	composerName, err := readJSONField(path.Join(dir, "backend", "composer.json"), func(b []byte) (string, error) {
		var c composerJSON
		if err = json.Unmarshal(b, &c); err != nil {
			return "", err
		}

		return c.Name, nil
	})
	if err != nil || composerName != goralysComposerName {
		return false
	}

	return true
}

func isGoralysCapRoot(dir string) bool {
	pkgName, err := readJSONField(path.Join(dir, "package.json"), func(b []byte) (string, error) {
		var p pkgJSON
		if err := json.Unmarshal(b, &p); err != nil {
			return "", err
		}

		return p.Name, nil
	})
	if err != nil || pkgName != goralysCapPkgName {
		return false
	}

	return true
}

func findRepoRoot(start string, checker func(string) bool) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	var i = 0
	for {
		hasGit, gitErr := DirExists(path.Join(dir, ".git"))
		if checker(dir) && hasGit && gitErr == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)

		if parent == dir || i > 5 {
			return "", errMonoRepoRootNotFound
		}

		dir = parent
		i++
	}
}

// FindRepoRoot retrieves the path of the repository root for the project in which the CLI tool is ran.
// It recursively tries to go up in the path tree with a maximum depth of 5.
func FindRepoRoot(start string, isMobile bool) (string, error) {
	if isMobile {
		return findRepoRoot(start, isGoralysCapRoot)
	}

	return findRepoRoot(start, isGoralysRoot)
}
