/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package utils

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Version represents a version with format Major.Minor.Patch (or M.m.p for short)
type Version struct {
	Major int
	Minor int
	Patch int
}

// IsNil checks if the version is different from 0.0.0 (uninitialized var -> should be nil)
func (v Version) IsNil() bool {
	return (v.Major | v.Minor | v.Patch) == 0b0
}

// ToString converts a version into a formatted string with the Major.minor.patch format (or M.m.p)
func (v Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// EnvVar represents an environment variable with its key and value
type EnvVar struct {
	Key   string
	Value string
}

// EnvCategory represents an environment category with its name and its list of environment variable
type EnvCategory struct {
	Name string
	Vars []EnvVar
}

// EnvFile represents an environment file with its path and its list of environment categories
type EnvFile struct {
	File string

	Categories []EnvCategory
}

// RawTemplate represents a YAML template's raw content. It contains the path of both the final output file and the
// reference file (JSON) that contais the version of the related ressource. It also contains the list of the versions
// for the template.
type RawTemplate struct {
	File string `yaml:"file"`
	Ref  string `yaml:"ref"`

	Versions map[string]yaml.Node `yaml:",inline"`
}

// RawRefFile represents a reference file (JSON) for a template. It contains the version field of the file.
type RawRefFile struct {
	Version string `json:"version"`
}
