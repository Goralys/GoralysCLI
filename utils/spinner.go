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

// SpinnerStep represents a step of a multistep spinner.
type SpinnerStep struct {
	Callback func() error
	Name     string
}

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

func overwriteN(n int, lines []string) (string, error) {
	var result strings.Builder
	_, err := fmt.Fprintf(&result, "\x1b[%dA", n) // returns n lines up
	if err != nil {
		return "", err
	}

	for _, line := range lines {
		// returns to the start of the line, clears it and then write the new content
		_, err := fmt.Fprintf(&result, "\n\r\x1b[2K%s", line)
		if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

func spinnerMultiStep(prefix string, label string, steps []SpinnerStep) error {
	currentStep := 0
	stepsPrinted := 0 // keeps track of the steps that actually printed in the console
	doneGlobal := make(chan struct{})

	// lines
	firstLine := fmt.Sprintf("%s %s", prefix, label)
	var doneLines []string
	currentLine := ""

	output := func() {
		lines := append([]string{firstLine}, doneLines...)
		if currentLine != "" {
			lines = append(lines, currentLine)
		}

		out, err := overwriteN(stepsPrinted+1, lines)
		if err != nil {
			return
		}

		fmt.Print(out)
		fmt.Print("\n")
	}

	// global ticker
	go func() {
		dots := []string{"", ".", "..", "..."}
		i := 0
		ticker := time.NewTicker(400 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-doneGlobal:
				return
			case <-ticker.C:
				firstLine = fmt.Sprintf("%s %s%-3s", prefix, label, dots[i%len(dots)])
				i++
				output()
			}
		}
	}()

	// steps runner
	for idx, current := range steps {
		done := make(chan struct{})
		currentStep++
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
					currentLine = fmt.Sprintf("%s-> %s%-3s", strings.Repeat(" ", len(prefix)+3), current.Name, dots[i%len(dots)])
					if stepsPrinted < idx+1 {
						stepsPrinted++
					}
					i++
				}
			}
		}()

		result := current.Callback()
		close(done)
		currentLine = fmt.Sprintf(
			"%s-> %s %s",
			strings.Repeat(" ", len(prefix)+3),
			current.Name,
			Colorize(ColorGreen, "[OK]"))
		if result != nil {
			close(doneGlobal)
			currentLine = fmt.Sprintf(
				"%s-> %s %s",
				strings.Repeat(" ", len(prefix)+3),
				current.Name,
				Colorize(ColorRed, "[FAIL]"))
			firstLine = fmt.Sprintf("%s %s %s", prefix, label, Colorize(ColorRed, "[FAIL]"))
			output()
			return result
		}

		doneLines = append(doneLines, currentLine)
		currentLine = ""
	}

	// if loop ends, then everything succeeded
	close(doneGlobal)
	firstLine = fmt.Sprintf("%s %s %s", prefix, label, Colorize(ColorGreen, "[OK]"))
	output()
	return nil
}

// func SpinnerMultiStep(label string, steps []step) error {
// 	return spinnerMultiStep(GoralysText(), label, steps)
// }

// SpinnerMultiStepNoPrefix starts a new multistep spinner without any prefix. Multistep spinners run the given list
// of steps and prints their execution results in the console.
// A step is essentially a functions that returns an error, if the error is nil, the step is considered successful
// (green [OK] in the console); other whies the step is considered failed (red [FAIL] in the console).
func SpinnerMultiStepNoPrefix(label string, steps []SpinnerStep) error {
	return spinnerMultiStep(strings.Repeat(" ", GoralysPrefixLen()), label, steps)
}
