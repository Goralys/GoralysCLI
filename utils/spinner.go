/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

// Package utils is the main package containing all the utilities functions for CLI tool.
package utils

import (
	"fmt"
	"strings"
	"time"
)

// StartSpinner starts a custom spinner (with the default CLI prefix) that prints out animated dots. It returns a
// function that takes one bool argument [ok]. When this function is called, the animation stops, and it prints out the
// status determined by the [ok] argument: if true, it prints [OK] in green; if false, it prints [FAIL] in red.
func StartSpinner(label string) func(ok bool) {
	return spinner(fmt.Sprintf("[%s]", GoralysText()), label)
}

// StartSpinnerNoPrefix starts a custom spinner (without the default CLI prefix) that prints out animated dots. It
// returns a function that takes one bool argument [ok]. When this function is called, the animation stops, and it
// prints out the status determined by the [ok] argument: if true, it prints [OK] in green; if false, it prints [FAIL]
// in red.
func StartSpinnerNoPrefix(label string) func(ok bool) {
	return spinner(strings.Repeat(" ", GoralysPrefixLen()), label)
}

func spinner(prefix string, label string) func(ok bool) {
	done := make(chan struct{})

	go func() {
		dots := []string{"", ".", "..", "..."}
		i := 0
		ticker := time.NewTicker(400 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fmt.Printf("\r%s %s%-3s", prefix, label, dots[i%len(dots)])
				i++
			}
		}
	}()

	return func(ok bool) {
		close(done)
		status := Colorize(ColorGreen, "[OK]")
		if !ok {
			status = Colorize(ColorRed, "[FAIL]")
		}
		fmt.Printf("\r%s %s %s%s\n", prefix, label, status, strings.Repeat(" ", 5))
	}
}
