/*
 * Copyright (C) 2026 Sami Saubion
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */
package utils

import (
	"fmt"
	"strings"
	"time"
)

func StartSpinner(label string) func(ok bool) {
	return spinner(fmt.Sprintf("[%s]", GoralysText()), label)
}

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
