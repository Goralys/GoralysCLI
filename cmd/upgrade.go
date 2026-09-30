/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"goralys-cli/utils"
	templates "goralys-cli/utils/templates"

	"github.com/spf13/cobra"
)

var version string

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "This command upgrades the tool to the latest version.",
	Long: `This commands upgrades the tool to the latest version by fetching the latest release's file from the official
		   server. It first autodetects the current OS and architecture and then overwrites the current install.
			
		   You can also upgrade to a specific version if you wish to by passing the desired version to the version
		   flag. Please note that you can never downgrade to an older version of the tool using this command.
		   Only a manual reinstallation of the tool can accomplish that goal.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		var target templates.Version
		var current templates.Version

		current, err := templates.ParseVersionFromString(goralysCLIVer)
		if err != nil {
			return fmt.Errorf("failed to parse invalid build version from metadata %s, %w", goralysCLIVer, err)
		}

		if version != "none" {
			target, err = templates.ParseVersionFromString(version)

			if err != nil {
				return fmt.Errorf("could not parse version string %s, %w", version, err)
			}
		}

		os := strings.ToLower(runtime.GOOS)
		arch := strings.ToLower(runtime.GOARCH)
		ext := ""
		if os == "windows" {
			ext = ".exe"
		}

		var url string

		if target.Equal(current) && !target.IsNil() {
			utils.Logf("Target version %s is identical to the current version", target.ToString())
			utils.Log("Nothing to do, aborting")
			return nil
		}

		if current.AtLeast(target) && !target.IsNil() {
			utils.Logf("Target version %s is older than the current version", target.ToString())
			utils.Log("Nothing to do, aborting")
			return nil
		}

		if !target.IsNil() {
			url = fmt.Sprintf(
				"https://cli.goralys.fr/release/%s/goralys-cli-%s-%s%s",
				target.ToString(),
				os,
				arch,
				ext,
			)
			utils.Logf("Upgrading to version %s", target.ToString())
		} else {
			url = fmt.Sprintf("https://cli.goralys.fr/release/latest/goralys-cli-%s-%s%s", os, arch, ext)
			utils.Log("Upgrading to latest version")
		}

		err = utils.UpgradeSelf(url)
		if err != nil {
			return err
		}

		utils.Log("Upgrade done")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(
		&version,
		"version",
		"v",
		"none",
		"This flag is used to pass an optional target version to the command.",
	)
}
