package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	source   string
	Segments []int
}

var versionPattern = regexp.MustCompile(`^(\d+(\.\d+){0,3})(-[a-z][0-9a-z-]*)?$`)

func NewVersion(v string) (*Version, error) {
	matches := versionPattern.FindStringSubmatch(v)
	if matches == nil {
		return nil, fmt.Errorf("Malformed version: %s", v)
	}

	segmentsStr := strings.Split(matches[1], ".")
	segments := make([]int, len(segmentsStr))
	for i, str := range segmentsStr {
		val, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version: %s: %s", v, err)
		}

		segments[i] = int(val)
	}

	return &Version{
		source:   v,
		Segments: segments,
	}, nil
}

func (c *Version) String() string {
	return c.source
}

func (c *Version) Compare(v *Version) int {
	// A quick, efficient equality check
	if c == v || c.source == v.source {
		return 0
	}

	// Compare the segments
	for i := 0; i < len(c.Segments); i++ {
		lhs := c.Segments[i]
		rhs := v.Segments[i]

		if lhs == rhs {
			continue
		} else if lhs < rhs {
			return -1
		} else {
			return 1
		}
	}

	panic("You created a black hole!")
}

func (c *Version) EqualTo(v *Version) bool {
	return c.Compare(v) == 0
}

func (c *Version) GreaterThan(v *Version) bool {
	return c.Compare(v) == 1
}

func (c *Version) LessThan(v *Version) bool {
	return c.Compare(v) == -1
}

func (c *Version) GreaterOrEqualTo(v *Version) bool {
	return c.Compare(v) > -1
}

func (c *Version) LesserOrEqualTo(v *Version) bool {
	return c.Compare(v) < 1
}
