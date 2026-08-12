/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var ErrMonoRepoRootNotFound = errors.New("goralys mono repo root not found")

const GoralysPkgName = "goralys"
const GoralysCapPkgName = "goralys-cap"
const GoralysComposerName = "goralys/goralys"

type PkgJson struct {
	Name string `json:"name"`
}

type ComposerJson struct {
	Name string `json:"name"`
}

func readJsonField(path string, extract func([]byte) (string, error)) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return extract(data)
}

func isGoralysRoot(dir string) bool {
	pkgName, err := readJsonField(path.Join(dir, "package.json"), func(b []byte) (string, error) {
		var p PkgJson
		if err := json.Unmarshal(b, &p); err != nil {
			return "", err
		}

		return p.Name, nil
	})
	if err != nil || pkgName != GoralysPkgName {
		return false
	}

	composerName, err := readJsonField(path.Join(dir, "backend", "composer.json"), func(b []byte) (string, error) {
		var c ComposerJson
		if err := json.Unmarshal(b, &c); err != nil {
			return "", err
		}

		return c.Name, nil
	})
	if err != nil || composerName != GoralysComposerName {
		return false
	}

	return true
}

func isGoralysCapRoot(dir string) bool {
	pkgName, err := readJsonField(path.Join(dir, "package.json"), func(b []byte) (string, error) {
		var p PkgJson
		if err := json.Unmarshal(b, &p); err != nil {
			return "", err
		}

		return p.Name, nil
	})
	if err != nil || pkgName != GoralysCapPkgName {
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
			return "", ErrMonoRepoRootNotFound
		}

		dir = parent
		i++
	}
}

func FindRepoRoot(start string, isMobile bool) (string, error) {
	if isMobile {
		return findRepoRoot(start, isGoralysCapRoot)
	}

	return findRepoRoot(start, isGoralysRoot)
}
