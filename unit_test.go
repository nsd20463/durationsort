package durationsort

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// test cases which should work
	cases := []struct {
		s string
		d time.Duration
	}{
		{"1h", time.Hour},
		{"1m", time.Minute},
		{"-1m", -time.Minute},
		{"1h2m3.5s", time.Hour + 2*time.Minute + 3*time.Second + time.Second/2},
		{"1d", 24 * time.Hour},
		{"1w", 7 * 24 * time.Hour},
		{"1y", 24 * time.Hour * (4*365 + 1) / 4},
		{"1y2w3d4h5m6s", (24 * time.Hour * (4*365 + 1) / 4) + (2 * 7 * 24 * time.Hour) + (3 * 24 * time.Hour) + 4*time.Hour + 5*time.Minute + 6*time.Second},
		{"5ns", 5 * time.Nanosecond},
		{"6us", 6 * time.Microsecond},
		{"7µs", 7 * time.Microsecond},
		{"8ms", 8 * time.Millisecond},
		{"9.876ms", 9876 * time.Microsecond},
	}

	for _, c := range cases {
		d, e := parseDuration(c.s)
		var v interface{}
		if e == nil {
			v = d
		} else {
			v = e
		}
		t.Logf("parseDuration(%q) = %v", c.s, v)

		if e != nil {
			t.Errorf("parseDuration(%q) returned error %v", c.s, e)
		}

		if d != c.d {
			t.Errorf("parseDuration(%q) returned unexpected result %v != expected  %v", c.s, d, c.d)
		}
	}
}

func TestParseFailures(t *testing.T) {
	// test cases which should fail
	cases := []string{
		"",
		"1h1d",
		"1.5d", // in theory we could allow non-integer days if we wanted it. stdlib does. "1.5h" is parsed as 1h30m
		"h",
		"1hm",
		"1dy", // "dy" is not an abbreviation for day
	}

	for _, s := range cases {
		d, e := parseDuration(s)
		var v interface{}
		if e == nil {
			v = d
		} else {
			v = e
		}
		t.Logf("parseDuration(%q) = %v", s, v)

		if e == nil {
			t.Errorf("parseDuration(%q) should have but did not return an error. It returned %v", s, d)
		}
	}
}

func TestExtract(t *testing.T) {
	// test cases which should work
	cases := []struct {
		s    string
		d    time.Duration
		left string
	}{
		{"0d", 0, ""},
		{"1h-", time.Hour, "-"},
		{"1m.", time.Minute, "."},
		{"-1m\"", -time.Minute, "\""},
		{"1d,2h", 24 * time.Hour, ",2h"},
	}

	for _, c := range cases {
		d, s, ok := extractDuration(c.s)
		t.Logf("extractDuration(%q) = %v, %q, %v", c.s, d, s, ok)

		if !ok {
			t.Errorf("extractDuration(%q) failed", c.s)
		} else if d != c.d || s != c.left {
			t.Errorf("extractDuration(%q) returned unexpected result %v,%q != expected  %v,%q", c.s, d, s, c.d, c.left)
		}
	}
}

func TestEqual(t *testing.T) {
	// test cases for equality (Less(x,x) should be false)
	cases := []string{
		"",
		"a",
		"0",
		"0h",
		"-1m",
	}

	for _, x := range cases {
		if Less(x, x) {
			t.Errorf("Less(%q, %q) != false", x, x)
		}
	}
}

func TestSort(t *testing.T) {
	result := []string{"1y", "1d", "1us", "1ms", "0h1m0s", "1s", "2.5ms"}
	Strings(result)
	t.Logf("%q", result)
	prev := time.Duration(0)
	for _, s := range result {
		n, err := parseDuration(s)
		if err != nil {
			t.Errorf("unexpected error %s", err)
			break
		}
		if prev >= n {
			t.Errorf("out of order. %v came before %v", prev, n)
		}
		prev = n
	}
}
