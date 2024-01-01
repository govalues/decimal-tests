//go:build mutation

package decimal_test

import (
	"testing"

	"github.com/gtramontina/ooze"
)

func TestMutation(t *testing.T) {
	ooze.Release(t,
		ooze.WithRepositoryRoot("../decimal"),
		ooze.WithTestCommand("go test -timeout 10s"),
		ooze.Parallel(),
	)
}
