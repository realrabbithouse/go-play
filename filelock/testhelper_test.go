package filelock

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("FILELOCK_HELPER_PROCESS") != "1" {
		return
	}

	failed := os.Getenv("FILELOCK_TEST_FAILED")

	var args []string
	for i, arg := range os.Args {
		if arg == "--" {
			args = os.Args[i+1:]
		}
	}
	if len(args) < 2 {
		t.Fatal("Usage: go test -test.run=TestHelperProcess -- <action> <path> [options]")
	}

	action := args[0]
	path := args[1]
	opts, hold := parseOptions(t, args[2:])

	l, err := New(path)
	if err != nil {
		t.Fatal(err)
	}

	switch action {
	case "wlock":
		err := l.WLock(opts...)
		if failed != "" {
			if err == nil {
				t.Fatalf("%v expected WLock failed with %q, got nil", args, failed)
			}
			if !strings.Contains(err.Error(), failed) {
				t.Fatalf("%v expected WLock failed with %q, got %v", args, failed, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("%v expected WLock to succeed, got %v", args, err)
		}

		t.Logf("WLock acquired for %s, hold for %s", path, hold)
		time.Sleep(hold)
		if err := l.Unlock(); err != nil {
			t.Fatalf("%v expected Unlock to succeed, got %v", args, err)
		}
	case "rlock":
		err := l.RLock(opts...)
		if failed != "" {
			if err == nil {
				t.Fatalf("%v expected RLock failed with %q, got nil", args, failed)
			}
			if !strings.Contains(err.Error(), failed) {
				t.Fatalf("%v expected RLock failed with %q, got %v", args, failed, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("%v expected RLock to succeed, got %v", args, err)
		}

		t.Logf("RLock acquired for %s, hold for %s", path, hold)
		time.Sleep(hold)
		if err := l.Unlock(); err != nil {
			t.Fatalf("%v expected Unlock to succeed, got %v", args, err)
		}
	default:
		t.Fatalf("Unknown action: %s", action)
	}
}

func parseOptions(tb testing.TB, args []string) ([]Option, time.Duration) {
	hold := 3 * time.Second
	var opts []Option
	for _, arg := range args {
		switch {
		case arg == "--block":
			opts = append(opts, WithBlock())
		case arg == "--remove":
			opts = append(opts, WithRemove())
		case strings.HasPrefix(arg, "--timeout="):
			timeout, err := time.ParseDuration(strings.TrimPrefix(arg, "--timeout="))
			if err != nil {
				tb.Fatal("Invalid timeout: ", arg)
			}
			opts = append(opts, WithTimeout(timeout))
		case strings.HasPrefix(arg, "--hold="):
			var err error
			hold, err = time.ParseDuration(strings.TrimPrefix(arg, "--hold="))
			if err != nil {
				tb.Fatal("Invalid hold duration: ", arg)
			}
		default:
			tb.Fatal("Unknown option: ", arg)
		}
	}

	return opts, hold
}

func helperProcessArgs(args ...string) []string {
	return append([]string{"-test.paniconexit0", "-test.timeout=10m0s", "-test.v=true", "-test.run=TestHelperProcess", "--"}, args...)
}
