/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */
package utils

import "gopkg.in/yaml.v3"

type Version struct {
	Major int
	Minor int
	Patch int
}

type EnvVar struct {
	Key   string
	Value string
}

type EnvCategory struct {
	Name string
	Vars []EnvVar
}

type EnvFile struct {
	File string

	Categories []EnvCategory
}

type RawTemplate struct {
	File string `yaml:"file"`
	Ref  string `yaml:"ref"`

	Versions map[string]yaml.Node `yaml:",inline"`
}

type RawRefFile struct {
	Version string `json:"version"`
}
