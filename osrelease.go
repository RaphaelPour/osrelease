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
	major, minor, patch int
	suffix              string
}

type Option func(v *Version)

func WithSuffix(suffix string) Option {
	return func(v *Version) {
		v.suffix = suffix
	}
}

func New(major, minor, patch int, options ...Option) Version {
	v := Version{
		major: major,
		minor: minor,
		patch: patch,
	}

	for _, option := range options {
		option(&v)
	}

	return v
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
	version.major, err = strconv.Atoi(matches[1])
	if err != nil {
		return EmptyVersion, fmt.Errorf(
			"%w: %w: major version should be int, is '%s'",
			ErrParseVersionNumber,
			err,
			matches[1],
		)
	}

	version.minor, err = strconv.Atoi(matches[2])
	if err != nil {
		return EmptyVersion, fmt.Errorf(
			"%w: %w: minor version should be int, is '%s'",
			ErrParseVersionNumber,
			err,
			matches[2],
		)
	}

	version.patch, err = strconv.Atoi(matches[3])
	if err != nil {
		return EmptyVersion, fmt.Errorf(
			"%w: %w: patch version should be int, is '%s'",
			ErrParseVersionNumber,
			err,
			matches[3],
		)
	}

	version.suffix = matches[4]

	return version, nil
}

func (v Version) Original() string {
	return fmt.Sprintf("%d.%d.%d%s", v.major, v.minor, v.patch, v.suffix)
}

func (v Version) String() string {
	return v.Original()
}

func (v Version) Major() int {
	return v.major
}

func (v Version) Minor() int {
	return v.minor
}

func (v Version) Patch() int {
	return v.patch
}

func (v Version) Suffix() string {
	return v.suffix
}

func (v Version) NewerThan(other Version) bool {
	if v.major > other.major {
		return true
	} else if v.major == other.major {
		if v.minor > other.minor {
			return true
		} else if v.minor == other.minor {
			return v.patch > other.patch
		}
	}

	return false
}

func (v Version) NewerThanOrEqual(other Version) bool {
	if v.major > other.major {
		return true
	} else if v.major == other.major {
		if v.minor > other.minor {
			return true
		} else if v.minor == other.minor {
			return v.patch >= other.patch
		}
	}

	return false
}

func (v Version) OlderThan(other Version) bool {
	if v.major < other.major {
		return true
	} else if v.major == other.major {
		if v.minor < other.minor {
			return true
		} else if v.minor == other.minor {
			return v.patch < other.patch
		}
	}

	return false
}

func (v Version) OlderThanOrEqual(other Version) bool {
	if v.major < other.major {
		return true
	} else if v.major == other.major {
		if v.minor < other.minor {
			return true
		} else if v.minor == other.minor {
			return v.patch <= other.patch
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
