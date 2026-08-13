/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const phpMinMajor int = 8
const phpMinMinor int = 5

var errComposerNotFound = errors.New("composer not found in path")
var errPhpNotFound = fmt.Errorf("php >= %d.%d not found in path", phpMinMajor, phpMinMinor)

func isPhpBinAtLeast(phpBin string) (bool, error) {
	out, err := exec.Command(phpBin, "-r", "echo PHP_MAJOR_VERSION . '|' . PHP_MINOR_VERSION;").Output()
	if err != nil {
		return false, err
	}

	parts := strings.SplitN(strings.TrimSpace(string(out)), "|", 2)

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, err
	}

	if major != phpMinMajor {
		return major > phpMinMajor, nil
	}

	return minor >= phpMinMinor, nil
}

func phpCandidates() []string {
	if runtime.GOOS == "windows" {
		return []string{
			"php85",
			`C:\php85\php.exe`,
			`C:\xampp\php85\php.exe`,
			"php",
		}
	}

	return []string{
		"php85",
		"/opt/alt/php85/usr/bin/php",
		"php",
	}
}

// ResolvePhp is used to locate the php executable.
func ResolvePhp() (string, error) {
	for _, candidate := range phpCandidates() {
		bin, err := exec.LookPath(candidate)
		if err != nil {
			continue
		}

		if ok, _ := isPhpBinAtLeast(bin); ok {
			return bin, nil
		}
	}

	return "", errPhpNotFound
}

// ResolveComposer is used to retrieve the composer executable and return it as an
// executable command
func ResolveComposer(phpBin string, args ...string) (*exec.Cmd, error) {
	composerBin, err := exec.LookPath("composer")
	if err != nil {
		return nil, errComposerNotFound
	}

	if runtime.GOOS == "windows" {
		return exec.Command(composerBin, args...), nil
	}

	cmdArgs := append([]string{composerBin}, args...)
	return exec.Command(phpBin, cmdArgs...), nil
}

// ResolvePnpm locates the pnpm executable and returns it as an executable command
func ResolvePnpm(args ...string) (*exec.Cmd, error) {
	pnpmBin, err := exec.LookPath("pnpm")
	if err != nil {
		return nil, err
	}

	return exec.Command(pnpmBin, args...), nil
}

// ResolveNpx locates the npx executable and returns it as an executable command
func ResolveNpx(args ...string) (*exec.Cmd, error) {
	npxBin, err := exec.LookPath("npx")
	if err != nil {
		return nil, err
	}

	return exec.Command(npxBin, args...), nil
}
