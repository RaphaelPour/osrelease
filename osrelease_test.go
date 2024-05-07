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
	require.Greater(t, v.Major, 0)

	require.Equal(t, out, v.Original)
	require.Equal(
		t, out, fmt.Sprintf(
			"%d.%d.%d%s",
			v.Major,
			v.Minor,
			v.Patch,
			v.Suffix,
		),
	)
}
