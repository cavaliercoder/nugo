package main

import (
	"testing"
)

func TestVersions(t *testing.T) {
	// version strings to test
	tests := []string{
		"1.0.0.0",         // 0
		"0.1.0.0",         // 1
		"0.0.1.0",         // 2
		"0.0.0.1",         // 3
		"0.999.999.999",   // 4
		"0.0.999.999",     // 5
		"0.0.0.999",       // 6
		"999.999.999.999", // 7
	}

	// parse all version strings
	versions := make([]*Version, len(tests))
	for i, test := range tests {
		v, err := NewVersion(test)
		if err != nil {
			t.Fatalf(err.Error())
		}

		versions[i] = v
	}

	// test greater than
	greater := map[*Version][]*Version{
		versions[0]: []*Version{versions[1], versions[2], versions[3], versions[4], versions[5], versions[6]},
		versions[1]: []*Version{versions[2], versions[3], versions[5], versions[6]},
		versions[2]: []*Version{versions[3], versions[6]},
	}

	for h, ls := range greater {
		for _, l := range ls {
			if !h.GreaterThan(l) {
				t.Fatalf("Version test %s > %s failed", h, l)
			}

			if l.GreaterThan(h) {
				t.Fatalf("Version test %s !> %s failed", h, l)
			}
		}
	}

	// test less than
	lesser := map[*Version][]*Version{
		versions[0]: []*Version{versions[7]},
		versions[1]: []*Version{versions[0], versions[4], versions[7]},
		versions[2]: []*Version{versions[0], versions[1], versions[4], versions[5], versions[7]},
	}

	for l, hs := range lesser {
		for _, h := range hs {
			if !l.LessThan(h) {
				t.Fatalf("Version test %s < %s failed", l, h)
			}

			if h.LessThan(l) {
				t.Fatalf("Version test %s !< %s failed", l, h)
			}
		}
	}
}
