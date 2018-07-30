/*
  Duration sorting.

  Copyright 2018 Nicolas S. Dade
*/
package durationsort

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Less returns true if x < y in a duration-aware comparison.
// It is suitable for use with Go's standard sort.Interface.
func less(a, b string) bool {
	// the idea is to scan along a and b rune-by-rune (might as well the UTF-8 ready),
	// until a numeric [0-9] or '-' rune is reached in both strings. Then I attempt to decode
	// the strings as if they were durations. If that works, the durations are compared.
	// Otherwise they are compared as text.
	digits := "-0123456789"

	for {
		i := strings.IndexAny(a, digits)
		j := strings.IndexAny(b, digits)
		if i < 0 || j < 0 {
			// no durations to compare. finish up with a straight string comparison
			return a < b
		}
		if i != j || a[:i] != b[:j] {
			// text differs. finish by comparing the text
			return a[:i] < b[:j]
		}
		// a and b match up to i (which equals j at this point); advance over the matching texts
		a = a[i:]
		b = b[j:]

		// and see if the next bytes can be decoded as a duration
		x, ar, okx := extractDuration(a)
		y, br, oky := extractDuration(b)
		if okx && oky {
			if x != y {
				return x < y
			}
			a, b = ar, br
		} else {
			// one or both of the strings doesn't start with a duration
			// compare them as text
			if len(a) > 0 && len(b) > 0 && a[0] != b[0] {
				return a[0] < b[0]
			}
			if len(a) == 0 {
				return true
			}
			if len(b) == 0 {
				return false
			}
			// advance forward over the digit
			a, b = a[1:], b[1:]
		}
	}
}
func Less(a, b string) bool {
	r := less(a, b)
	return r
}

// extractDuration extracts the duration prefix of a.
// It returns the decoded duration and the remaining non-duration part of a.
// If a doesn't start with a duration then (0, a, false) is returned.
func extractDuration(a string) (time.Duration, string, bool) {
	var i int
	var r rune
loop:
	for i, r = range a {
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'y', 'w', 'd', 'h', 'm', 's', 'u', 'µ', 'n', '.':
			// all possible chars of a duration
		case '-': // - can only be a leading char. if it appears later then it terminates the duration
			if i != 0 {
				i--
				break loop
			}
		default:
			i--
			break loop
		}
	}
	i++
	dur, err := parseDuration(a[:i])
	if err == nil {
		return dur, a[i:], true
	}
	return 0, a, false
}

// parseDuration parses a text-encoded duration
// I've extended the text decoding to understand 'y','w' and 'd' as rough measures matching what humans expect
func parseDuration(str string) (time.Duration, error) {
	orig_str := str // save for error msgs

	// remove any leading minus sign and save it for later
	negative := false
	if len(str) > 0 && str[0] == '-' {
		negative = true
		str = str[1:]
	}

	// first process any y,w or d sections
	var days, weeks, years int64
	units := "ywd"
loop:
	for {
		idx := strings.IndexAny(str, units)
		if idx == -1 {
			break loop
		}
		// see if what comes before the letter is a decimal number
		x, err := strconv.ParseInt(str[:idx], 10, 64)
		if err != nil {
			// for the moment we don't suport non-integer d/w/y values
			return 0, fmt.Errorf("Can't parse %q as a duration: %q is not a number: %v", orig_str, str[:idx], err)
		}
		// there shouldn't be any embedded negative signs
		if x < 0 {
			return 0, fmt.Errorf("Can't parse %q as a duration: %q is not a positive number", orig_str, str[:idx])
		}
		unit := str[idx]
		str = str[idx+1:]
		switch unit {
		case 'y':
			years = x
			units = "wd"
		case 'w':
			weeks = x
			units = "d"
		case 'd':
			days = x
			break loop
		}
	}

	// if there's anything left, parse it with the stdlib's parser
	var dur time.Duration
	if str != "" {
		var err error
		dur, err = time.ParseDuration(str)
		if err != nil {
			return 0, err
		}
	} else {
		// don't allow "" as the full argument
		if orig_str == "" {
			return 0, fmt.Errorf("Can't parse %q as a duration", orig_str)
		}
	}

	dur += time.Duration(days) * time.Hour * 24
	dur += time.Duration(weeks) * time.Hour * 24 * 7
	dur += time.Duration(years) * time.Hour * 24 * (4*365 + 1) / 4 // use Julian years like humans do

	if negative {
		dur = -dur
	}

	return dur, nil
}

// StringSlice is a slice of strings sortable in duration-aware order
// It implements sort.Interface. The name matches the equivalent type in the standard sort package.
type StringSlice []string

func (s StringSlice) Len() int           { return len(s) }
func (s StringSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s StringSlice) Less(i, j int) bool { return Less(s[i], s[j]) }

// Strings is a utility function to sort a slice of strings using duration-aware sort order.
// The name matched the name of the equivalent function in the standard sort package.
func Strings(strs []string) {
	sort.Sort(StringSlice(strs))
}
