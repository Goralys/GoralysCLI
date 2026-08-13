/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool.
package utils

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/term"
)

const (
	// ColorReset resets terminal styling back to default.
	ColorReset = "\033[0m"

	// ColorRed is used for failure/error status output.
	ColorRed = "\033[31m"

	// ColorGreen is used for success status output.
	ColorGreen = "\033[32m"

	// ColorYellow is used for warning output. [unused; uncomment if needed]
	// ColorYellow = "\033[33m"

	// ColorBlue is currently unused. [unused; uncomment if needed]
	// ColorBlue = "\033[34m"

	// ColorGoralys is the brand color used for the [GoralysCLI] prefix.
	ColorGoralys = "\033[94m"

	// ColorCyan is used to color subprocess labels (e.g. [pnpm], [composer]).
	ColorCyan = "\033[36m"

	// PromptYes is the expected input for an affirmative [Y/n] prompt.
	PromptYes = "y"

	// PromptNo is the expected input for a negative [Y/n] prompt.
	PromptNo = "n"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// 0 = don't skip, 1 = skip [no], 2 = skip [yes]
var skipPrompt = 0

// InitSkipPrompt initializes the prompt skip mechanism. If the answer argument is true, the prompts will be skipped
// with answer 'yes', else the prompts will be skipped with answer 'no'
func InitSkipPrompt(answer bool) {
	// false = no, true = yes (only called if prompts should be skipped
	if answer {
		skipPrompt = 2
	} else {
		skipPrompt = 1
	}
}

func stripAnsi(_s string) string {
	return ansiRegex.ReplaceAllString(_s, "")
}

// PrefixWriter represents a custom terminal output. It contains a custom prefix that is printed before the content of
// the base output.
type PrefixWriter struct {
	label       string
	dest        io.Writer
	midProgress bool
}

func supportColor() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Colorize returns a string's colorized version with the CLI pre-defined colors (using ansi codes).
func Colorize(color, _s string) string {
	if !supportColor() {
		return _s
	}

	return color + _s + ColorReset
}

// GoralysText returns the colorized prefix for the CLI logs.
func GoralysText() string {
	return Colorize(ColorGoralys, "GoralysCLI")
}

// GoralysPrefixLen returns the length of the CLI's base prefix
func GoralysPrefixLen() int {
	return len(stripAnsi(GoralysText())) + 2
}

// Log outputs a message with the CLI's default prefix to the console.
func Log(msg string) {
	fmt.Printf("[%s] %s \n", GoralysText(), msg)
}

// LogNoPrefix outputs a message to the console without the CLI's default prefix (replaced by leading spaces)
func LogNoPrefix(msg string) {
	fmt.Printf("%s %s \n", strings.Repeat(" ", GoralysPrefixLen()), msg)
}

// Logf outputs a formatted message with the CLI's default prefix to the console.
func Logf(format string, a ...any) {
	Log(fmt.Sprintf(format, a...))
}

// LogfNoPrefix outputs a formatted message to the console without the CLI's default prefix (replaced by leading spaces)
func LogfNoPrefix(format string, a ...any) {
	LogNoPrefix(fmt.Sprintf(format, a...))
}

// NewPrefixWriter creates a new prefix writer with a custom label (prefix) for a given output.
func NewPrefixWriter(label string, dest io.Writer) *PrefixWriter {
	return &PrefixWriter{label: label, dest: dest}
}

// Write writes a given content (byte array) to a PrefixWriter's output.
func (w *PrefixWriter) Write(p []byte) (int, error) {
	text := string(p)

	// Split on \r or \n, keeping track of which delimiter followed each segment.
	i := 0
	for i < len(text) {
		j := strings.IndexAny(text[i:], "\r\n")
		if j == -1 {
			// no delimiter left in this chunk; nothing to flush yet
			break
		}
		j += i
		segment := text[i:j]
		delim := text[j]

		// \r\n counts as a single line ending
		next := j + 1
		if delim == '\r' && next < len(text) && text[next] == '\n' {
			next++
			delim = '\n'
		}

		if segment != "" || delim == '\n' {
			w.writeSegment(segment, delim == '\r')
		}

		i = next
	}

	return len(p), nil
}

func (w *PrefixWriter) writeSegment(line string, isOverwrite bool) {
	if isOverwrite {
		if _, err := fmt.Fprintf(w.dest, "\r[%s] => %s", Colorize(ColorCyan, w.label), line); err != nil {
			return
		}
		w.midProgress = true
		return
	}

	if line == "" {
		return
	}

	if w.midProgress {
		if _, err := fmt.Fprintln(w.dest); err != nil {
			return
		}
		w.midProgress = false
	}
	if _, err := fmt.Fprintf(w.dest, "[%s] => %s\n", Colorize(ColorCyan, w.label), line); err != nil {
		return
	}
}

// Prompt prompts the user with a given message (question) and writes the output the given string variable (dest).
func Prompt(dest *string, msg string) {
	if skipPrompt != 0 {
		*dest = []string{PromptNo, PromptYes}[skipPrompt-1]
		return
	}

	fmt.Printf("[%s] %s ", GoralysText(), msg)
	if _, err := fmt.Scan(dest); err != nil {
		return
	}
}

// Promptf prompts the user with a given formatted message (question) and writes the output the given string variable
// (dest).
func Promptf(dest *string, format string, a ...any) {
	Prompt(dest, fmt.Sprintf(format, a...))
}

// PromptBool prompts the user for a given yes or no question and writes the output to the given bool variable (dest).
func PromptBool(dest *bool, msg string) {
	var raw string
	Promptf(&raw, "%s [Y/n]: ", msg)

	*dest = strings.ToLower(strings.TrimSpace(raw)) == "y"
}

// PromptfBool prompts the user for a given yes or no question (formatted) and writes the output to the given bool
// variable (dest).
func PromptfBool(dest *bool, format string, a ...any) {
	var raw string
	Promptf(&raw, "%s [Y/n]: ", fmt.Sprintf(format, a...))

	*dest = strings.ToLower(strings.TrimSpace(raw)) == "y"
}
