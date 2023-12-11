package run

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	m := New()
	err := m.Register(&MyThing{})
	if err != nil {
		require.NoError(t, err)
	}
	ctx := context.Background()
	err = m.Call(ctx, "mything:test123", []string{"1", "2"})
	require.NoError(t, err)
	// t.Fail()
}

type MyThing struct{}

func (m *MyThing) Test123(ctx context.Context, args []string) error {
	return nil
}

func (m *MyThing) private(xx string) {}

// testCmd := fmt.Sprintf(
// 	"set -o pipefail; go test -coverprofile=%s -race -test.v ./... ./pkg/... ./apis/... | tee %s",
// 	locations.UnitTestCoverageReport(), locations.UnitTestStdOut(),
// )

// // cgo needed to enable race detector -race
// testErr := sh.RunWithV(map[string]string{"CGO_ENABLED": "1"}, "bash", "-c", testCmd)
// must(sh.RunV("bash", "-c", "set -o pipefail; cat "+locations.UnitTestStdOut()+" | go tool test2json > "+locations.UnitTestExecReport()))
// must(testErr)
