/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */
package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/term"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorGoralys = "\033[94m"
	ColorCyan    = "\033[36m"

	PromptYes = "y"
	PromptNo  = "n"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// 0 = don't skip, 1 = skip [no], 2 = skip [yes]
var skipPrompt int = 0

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

type PrefixWriter struct {
	label string
	dest  io.Writer
}

func supportColor() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func Colorize(color, _s string) string {
	if !supportColor() {
		return _s
	}

	return color + _s + ColorReset
}

func GoralysText() string {
	return Colorize(ColorGoralys, "GoralysCLI")
}

func GoralysPrefixLen() int {
	return len(stripAnsi(GoralysText())) + 2
}

func Log(msg string) {
	fmt.Printf("[%s] %s \n", GoralysText(), msg)
}

func LogNoPrefix(msg string) {
	fmt.Printf("%s %s \n", strings.Repeat(" ", GoralysPrefixLen()), msg)
}

func Logf(format string, a ...any) {
	Log(fmt.Sprintf(format, a))
}

func LogfNoPrefix(format string, a ...any) {
	LogNoPrefix(fmt.Sprintf(format, a))
}

func NewPrefixWriter(label string, dest io.Writer) *PrefixWriter {
	return &PrefixWriter{label: label, dest: dest}
}

func (w *PrefixWriter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(p)))
	for scanner.Scan() {
		if _, err := fmt.Fprintf(w.dest, "[%s] => %s\n", Colorize(ColorCyan, w.label), scanner.Text()); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

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

func Promptf(dest *string, format string, a ...any) {
	Prompt(dest, fmt.Sprintf(format, a...))
}
