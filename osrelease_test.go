package osrelease

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	rawOut, err := exec.Command("uname", "-r").Output()
	out := strings.TrimSpace(string(rawOut))
	require.NoError(t, err)

	v, err := Parse()
	require.NoError(t, err)
	require.Greater(t, v.Major(), 0)

	require.Equal(t, out, v.Original())
	require.Equal(
		t, out, fmt.Sprintf(
			"%d.%d.%d%s",
			v.Major(),
			v.Minor(),
			v.Patch(),
			v.Suffix(),
		),
	)
}

func TestNewVersion(t *testing.T) {
	v := New(1, 2, 3)
	require.Equal(t, 1, v.Major())
	require.Equal(t, 2, v.Minor())
	require.Equal(t, 3, v.Patch())
	require.Empty(t, v.Suffix())
	require.Equal(t, "1.2.3", v.Original())
}

func TestNewVersionWithSuffix(t *testing.T) {
	v := New(1, 2, 3, WithSuffix("-gabc+3"))
	require.Equal(t, 1, v.Major())
	require.Equal(t, 2, v.Minor())
	require.Equal(t, 3, v.Patch())
	require.Equal(t, "-gabc+3", v.Suffix())
	require.Equal(t, "1.2.3-gabc+3", v.Original())
}

func TestExclusiveCompare(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		v1, v2    string
		expection bool
	}{
		{name: "6.6.2 > 6.6.1", v1: "6.6.2", v2: "6.6.1", expection: true},
		{name: "6.7.1 > 6.6.1", v1: "6.7.2", v2: "6.6.1", expection: true},
		{name: "7.6.1 > 6.6.1", v1: "7.6.2", v2: "6.6.1", expection: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			v1, err := ParseString(testCase.v1)
			require.NoError(t, err)

			v2, err := ParseString(testCase.v2)
			require.NoError(t, err)

			require.Equal(t, testCase.expection, v1.NewerThan(v2))
			require.Equal(t, !testCase.expection, v2.NewerThan(v1))

			require.Equal(t, !testCase.expection, v1.OlderThan(v2))
			require.Equal(t, testCase.expection, v2.OlderThan(v1))
		})
	}
}

func TestInclusiveCompare(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		v1, v2    string
		expection bool
	}{
		{name: "6.6.2 >= 6.6.1", v1: "6.6.2", v2: "6.6.1", expection: true},
		{name: "6.7.1 >= 6.6.1", v1: "6.7.2", v2: "6.6.1", expection: true},
		{name: "7.6.1 >= 6.6.1", v1: "7.6.2", v2: "6.6.1", expection: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			v1, err := ParseString(testCase.v1)
			require.NoError(t, err)

			v2, err := ParseString(testCase.v2)
			require.NoError(t, err)

			require.True(t, v1.NotEqual(v2))

			require.Equal(t, testCase.expection, v1.NewerThanOrEqual(v2))
			require.Equal(t, !testCase.expection, v2.NewerThanOrEqual(v1))

			require.Equal(t, !testCase.expection, v1.OlderThanOrEqual(v2))
			require.Equal(t, testCase.expection, v2.OlderThanOrEqual(v1))
		})
	}
}

func TestInclusiveComparisonWithEqualVersions(t *testing.T) {
	v1, err := ParseString("1.2.3")
	require.NoError(t, err)

	v2, err := ParseString("1.2.3")
	require.NoError(t, err)

	require.True(t, v1.Equal(v2))
	require.True(t, v1.OlderThanOrEqual(v2))
	require.True(t, v2.OlderThanOrEqual(v1))
	require.True(t, v2.NewerThanOrEqual(v1))
	require.True(t, v1.NewerThanOrEqual(v2))
}

func TestErrorCases(t *testing.T) {
	for _, testCase := range []struct {
		name string
		v    string
		err  error
	}{
		{"empty version", "", ErrMatchCountMismatch},
		{"bad major", "99999999999999999999999999999999999999999999999999999999999999999.0.0", ErrParseVersionNumber},
		{"bad minor", "0.99999999999999999999999999999999999999999999999999999999999999999.0", ErrParseVersionNumber},
		{"bad patch", "0.0.99999999999999999999999999999999999999999999999999999999999999999", ErrParseVersionNumber},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := ParseString(testCase.v)
			require.NotNil(t, err)
			require.ErrorIs(t, err, testCase.err)
		})
	}
}
