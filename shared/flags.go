/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package shared

// VersionFlag is a root level flag that triggers the "version command" (prints the tool's current version)
var VersionFlag bool

// MobileFlag is a persistent flag used to indicate that the tool is run for the GoralysCap repo
var MobileFlag bool

// BackendFlag is a persistent flag used to indicate that the tool is run for a "backend-only" environment
// (e.g.: server deployment in production)
var BackendFlag bool

// RepoRoot is a variable used to store the root directory of the repository
var RepoRoot string
