/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool
package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"goralys-cli/utils"

	"gopkg.in/yaml.v3"
)

func formatEnvValue(val any) string {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("\"%v\"", v) // already a string, e.g. from "localhost"
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("\"%v\"", v)
	}
}

func parseVersionFromString(_s string) (Version, error) {
	parts := strings.Split(_s, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("wrong version format, expected M.m.p got %s", _s)
	}

	M, err := strconv.Atoi(parts[0]) // Major
	if err != nil {
		return Version{}, err
	}
	m, err := strconv.Atoi(parts[1]) // minor
	if err != nil {
		return Version{}, err
	}
	p, err := strconv.Atoi(parts[2]) // patch
	if err != nil {
		return Version{}, err
	}

	return Version{M, m, p}, nil
}

// AtLeast compares a version against a target, it returns true if the version is superior or equal to the target
// and returns false otherwise
func (_v Version) AtLeast(target Version) bool {
	if _v.Major != target.Major {
		return _v.Major > target.Major
	}
	if _v.Minor != target.Minor {
		return _v.Minor > target.Minor
	}

	return _v.Patch >= target.Patch
}

// Equal compares a version against a target and returns true if the 2 are equal
func (_v Version) Equal(target Version) bool {
	return _v.Major == target.Major && _v.Minor == target.Minor && _v.Patch == target.Patch
}

func parseYamlTemplate(root string, data []byte) (EnvFile, error) {
	var raw RawTemplate
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return EnvFile{}, err
	}

	file, err := os.ReadFile(filepath.Join(root, raw.Ref))
	if err != nil {
		return EnvFile{}, err
	}

	var ref RawRefFile
	err = json.Unmarshal(file, &ref)
	if err != nil {
		return EnvFile{}, err
	}

	target, err := parseVersionFromString(ref.Version)
	if err != nil {
		return EnvFile{}, err
	}

	var result EnvFile
	result.File = raw.File

	for vName, vNode := range raw.Versions {
		v, err := parseVersionFromString(vName)
		if err != nil {
			return EnvFile{}, err
		}

		if !target.AtLeast(v) && !target.Equal(v) {
			continue
		}

		var dynamicVer map[string]map[string]any
		if err := vNode.Decode(&dynamicVer); err != nil {
			return EnvFile{}, err
		}

		// iterate through categories
		for catName, vars := range dynamicVer {
			var cat EnvCategory
			cat.Name = catName

			for key, val := range vars {
				strVal := formatEnvValue(val)
				cat.Vars = append(cat.Vars, EnvVar{Key: key, Value: strVal})
			}

			// append category
			result.Categories = append(result.Categories, cat)
		}
	}

	return result, nil
}

func buildEnvFileContents(file EnvFile) string {
	var contents strings.Builder

	for _, cat := range file.Categories {
		contents.WriteString("# " + cat.Name + "\n")
		for _, envVar := range cat.Vars {
			contents.WriteString(envVar.Key + "=" + envVar.Value + "\n")
		}
		contents.WriteString("\n")
	}

	return contents.String()
}

func mergeEnvFile(root string, original EnvFile) (EnvFile, error) {
	existingKeys := make(map[string]bool)
	existingVals := make(map[string]string)

	file, err := os.Open(filepath.Join(root, original.File))
	if err != nil {
		if os.IsNotExist(err) {
			goto filterStep
		}

		return EnvFile{}, err
	}
	{
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// skip empty and comments
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) >= 2 {
				key := strings.TrimSpace(parts[0])
				if key != "" {
					existingKeys[key] = true
					existingVals[key] = strings.TrimSpace(parts[1])
				}
			}
		}
		if err = file.Close(); err != nil {
			return EnvFile{}, err
		}
	}
filterStep:
	var filteredCategories []EnvCategory

	for _, category := range original.Categories {
		var filteredVars []EnvVar

		for _, v := range category.Vars {
			if !existingKeys[v.Key] {
				filteredVars = append(filteredVars, v)
			} else {
				filteredVars = append(filteredVars, EnvVar{Key: v.Key, Value: existingVals[v.Key]})
			}
		}
		if len(filteredVars) > 0 {
			newCategory := EnvCategory{
				Name: category.Name,
				Vars: filteredVars,
			}
			filteredCategories = append(filteredCategories, newCategory)
		}
	}

	return EnvFile{
		File:       original.File,
		Categories: filteredCategories,
	}, nil
}

// MakeEnvFileFromTemplate loads a YAML template file and parses it. Then it retrieves the version of the current copy
// of the Goralys project, it then outputs the generated content to the target file specified in the template.
func MakeEnvFileFromTemplate(root string, template string) error {

	stop := utils.StartSpinnerNoPrefix("-> Parsing template")
	env, err := parseYamlTemplate(root, []byte(template))
	if err != nil {
		stop(false)
		return err
	}
	stop(true)

	stop = utils.StartSpinnerNoPrefix("-> Merging files")
	finalEnv, err := mergeEnvFile(root, env)
	if err != nil {
		stop(false)
		return err
	}
	stop(true)

	stop = utils.StartSpinnerNoPrefix("-> Writing final content")
	content := buildEnvFileContents(finalEnv)
	err = os.WriteFile(filepath.Join(root, finalEnv.File), []byte(content), os.ModePerm)
	if err != nil {
		stop(false)
		return err
	}
	stop(true)

	return nil
}
