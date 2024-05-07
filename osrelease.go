package osrelease

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var (
	pattern = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)(.*)`)

	EmptyVersion = Version{}

	ErrReadOsRelease      = errors.New("error opening osrelease")
	ErrParseOsRelease     = errors.New("error parsing osrelease")
	ErrMatchCountMismatch = errors.New("match count mismatch")
	ErrParseVersionNumber = errors.New("error parsing version number")
)

type Version struct {
	Major, Minor, Patch int
	Suffix              string
	Original            string
}

func Parse() (Version, error) {
	raw, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return EmptyVersion, fmt.Errorf("%w: %w", ErrReadOsRelease, err)
	}

	return ParseString(string(raw))
}

func ParseString(s string) (Version, error) {
	matches := pattern.FindStringSubmatch(s)
	if len(matches) != 5 {
		return EmptyVersion, fmt.Errorf(
			"%w: %w: expected 5, got %d with %s",
			ErrParseOsRelease,
			ErrMatchCountMismatch,
			len(matches),
			s,
		)
	}

	var version Version
	var err error
	version.Major, err = strconv.Atoi(matches[1])
	if err != nil {
		return EmptyVersion, fmt.Errorf(
			"%w: %w: major version should be int, is '%s'",
			ErrParseVersionNumber,
			err,
			matches[1],
		)
	}

	version.Minor, err = strconv.Atoi(matches[2])
	if err != nil {
		return EmptyVersion, fmt.Errorf(
			"%w: %w: minor version should be int, is '%s'",
			ErrParseVersionNumber,
			err,
			matches[2],
		)
	}

	version.Patch, err = strconv.Atoi(matches[3])
	if err != nil {
		return EmptyVersion, fmt.Errorf(
			"%w: %w: patch version should be int, is '%s'",
			ErrParseVersionNumber,
			err,
			matches[3],
		)
	}

	version.Original = matches[0]
	version.Suffix = matches[4]

	return version, nil
}

func (v Version) NewerThan(other Version) bool {
	if v.Major > other.Major {
		return true
	} else if v.Major == other.Major {
		if v.Minor > other.Minor {
			return true
		} else if v.Minor == other.Minor {
			return v.Patch > other.Patch
		}
	}

	return false
}

func (v Version) NewerThanOrEqual(other Version) bool {
	if v.Major > other.Major {
		return true
	} else if v.Major == other.Major {
		if v.Minor > other.Minor {
			return true
		} else if v.Minor == other.Minor {
			return v.Patch >= other.Patch
		}
	}

	return false
}

func (v Version) OlderThan(other Version) bool {
	if v.Major < other.Major {
		return true
	} else if v.Major == other.Major {
		if v.Minor < other.Minor {
			return true
		} else if v.Minor == other.Minor {
			return v.Patch < other.Patch
		}
	}

	return false
}

func (v Version) OlderThanOrEqual(other Version) bool {
	if v.Major < other.Major {
		return true
	} else if v.Major == other.Major {
		if v.Minor < other.Minor {
			return true
		} else if v.Minor == other.Minor {
			return v.Patch <= other.Patch
		}
	}

	return false
}

func (v Version) Equal(other Version) bool {
	return v == other
}

func (v Version) NotEqual(other Version) bool {
	return v != other
}
