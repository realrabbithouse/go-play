package filelock

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFileLock_RLock_success(t *testing.T) {
}

func TestFileLock_RLock_failed(t *testing.T) {
	tmpdir, err := os.MkdirTemp(".", "dir-")
	require.NoError(t, err)

	file, err := filepath.Abs(filepath.Join(tmpdir, "target"))
	require.NoError(t, err)

	cmd := exec.Command(os.Args[0], helperProcessArgs("wlock", file)...)
	cmd.Env = append(os.Environ(), "FILELOCK_HELPER_PROCESS=1")
	require.NoError(t, cmd.Start())
	go func() {
		require.NoError(t, cmd.Wait())
	}()

	time.Sleep(time.Second)
	l, err := New(file)
	require.NoError(t, err)

	require.NoError(t, l.RLock(WithTimeout(time.Second)))
}

func TestFileLock_WLock_success(t *testing.T) {
}
