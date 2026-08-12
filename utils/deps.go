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

const PhpMinMajor int = 8
const PhpMinMinor int = 5

var ErrComposerNotFound = errors.New("composer not found in path")
var ErrPhpNotFound = errors.New(fmt.Sprintf("php >= %d.%d not found in path", PhpMinMajor, PhpMinMinor))

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

	if major != PhpMinMajor {
		return major > PhpMinMajor, nil
	}

	return minor >= PhpMinMinor, nil
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

	return "", ErrPhpNotFound
}

func ResolveComposer(phpBin string, args ...string) (*exec.Cmd, error) {
	composerBin, err := exec.LookPath("composer")
	if err != nil {
		return nil, ErrComposerNotFound
	}

	if runtime.GOOS == "windows" {
		return exec.Command(composerBin, args...), nil
	}

	cmdArgs := append([]string{composerBin}, args...)
	return exec.Command(phpBin, cmdArgs...), nil
}

func ResolvePnpm(args ...string) (*exec.Cmd, error) {
	pnpmBin, err := exec.LookPath("pnpm")
	if err != nil {
		return nil, err
	}

	return exec.Command(pnpmBin, args...), nil
}
