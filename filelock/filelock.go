//go:build dragonfly || freebsd || linux || netbsd

package filelock

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

var (
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

	ErrTimeout         = errors.New("acquiring lock timeout")
	ErrNotAbsolutePath = errors.New("lock path is not absolute")
)

const (
	defaultLockTimeout = 30 * time.Second

	minWaitDuration = 200 * time.Millisecond
	maxWaitDuration = 600 * time.Millisecond
)

type config struct {
	// ctx is the context for the lock operation.
	ctx context.Context

	// timeout is the maximum amount of time to wait for a lock.
	// Defaults to defaultLockTimeout.
	timeout time.Duration

	// block is a flag that indicates whether to block until the lock is acquired.
	block bool

	// remove is a flag that indicates whether to remove the lock file when the lock is released.
	remove bool
}

// Option is a function type that can be used to customize the behavior of a FileLock.
type Option func(*config)

// WithContext returns an Option that sets the context for the lock operation.
// If no context is provided, the function will use the background context.
func WithContext(ctx context.Context) Option {
	return func(o *config) { o.ctx = ctx }
}

// WithTimeout returns an Option that sets the lock timeout duration.
// This function allows customization of the maximum amount of time to
// wait for a lock.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) { c.timeout = timeout }
}

// WithBlock returns an Option that sets the lock to block until acquired.
// This option configures the FileLock to wait indefinitely until the lock can
// be acquired, rather than returning immediately if the lock is not available.
func WithBlock() Option {
	return func(c *config) { c.block = true }
}

// WithRemove returns an Option that sets the lock to be removed when released.
// This option configures the FileLock to automatically remove the lock file
// when the lock is released, cleaning up any temporary files created during
// the locking process.
func WithRemove() Option {
	return func(c *config) { c.remove = true }
}

type FileLock struct {
	config *config    // config is the configuration for the lock operation
	path   string     // path is the target path which the FileLock protects
	file   *os.File   // file is the underlying file descriptor used for locking
	mu     sync.Mutex // guard against FileLock
}

// New creates and returns a new FileLock instance.
//
// It takes a path to the file that needs to be locked and creates a lock file
// with the same name plus a ".lock" extension in the same directory.
func New(path string) (*FileLock, error) {
	if !filepath.IsAbs(path) {
		return nil, ErrNotAbsolutePath
	}

	dir := filepath.Dir(path)
	name := filepath.Base(path)

	file, err := os.OpenFile(filepath.Join(dir, name+".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}

	return &FileLock{
		config: &config{
			ctx:     context.Background(),
			timeout: defaultLockTimeout,
			block:   false,
			remove:  false,
		},
		path: path,
		file: file,
	}, nil
}

// RLock acquires a shared lock on behalf of the current process on the file represented
// by the FileLock. If there is another process already held an exclusive lock on the file,
// RLock blocks until the lock is available.
//
// RLock optionally accepts a variable number of Option functions to customize the lock behavior.
func (l *FileLock) RLock(opts ...Option) error {
	for _, opt := range opts {
		opt(l.config)
	}

	// Create a Flock_t structure.
	lock := unix.Flock_t{
		Type:   unix.F_RDLCK, // read lock, shared lock
		Whence: 0,            // relative to the start of the file
		Start:  0,            // lock starts at byte 0
		Len:    0,            // lock extends to EOF
	}

	return l.acquireLock(lock)
}

// WLock acquires an exclusive lock on behalf of the current process on the file represented
// by the FileLock. If there is another process already holding a lock on the file,
// WLock blocks until the lock is available.
//
// If the WithBlock option is provided, WLock will use the fcntl F_SETLK syscall with the
// F_SETLKW operation, indicating that WLock wants to wait until the lock can be acquired,
// rather than returning immediately if the lock is not available.
//
// WLock optionally accepts a variable number of Option functions to customize the lock behavior.
func (l *FileLock) WLock(opts ...Option) error {
	for _, opt := range opts {
		opt(l.config)
	}

	if l.config.block {
		return l.acquireLockWait()
	}

	// Create a Flock_t structure.
	lock := unix.Flock_t{
		Type:   unix.F_WRLCK, // write lock, exclusive lock
		Whence: 0,            // relative to the start of the file
		Start:  0,            // lock starts at byte 0
		Len:    0,            // lock extends to EOF
	}

	return l.acquireLock(lock)
}

// Unlock releases the lock held by the FileLock.
//
// Unlock first releases the lock on the underlying file using the fcntl F_SETLK
// syscall with the F_UNLCK operation. After releasing the lock, it closes the
// file descriptor associated with the lock file.
func (l *FileLock) Unlock() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create a Flock_t structure.
	lock := unix.Flock_t{
		Type:   unix.F_UNLCK, // unlock
		Whence: 0,            // relative to the start of the file
		Start:  0,            // lock starts at byte 0
		Len:    0,            // lock extends to EOF
	}

	// Apply the lock using fcntl.
	if err := unix.FcntlFlock(l.file.Fd(), unix.F_SETLK, &lock); err != nil {
		return fmt.Errorf("releasing lock: %w", err)
	}

	err := l.file.Close()
	if l.config.remove {
		removeErr := os.Remove(l.path + ".lock")
		return errors.Join(err, removeErr)
	}

	return err
}

func (l *FileLock) acquireLock(lock unix.Flock_t) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Start a goroutine to enforce the timeout.
	timeoutC := time.After(l.config.timeout)

	for {
		select {
		case <-timeoutC:
			return ErrTimeout
		case <-l.config.ctx.Done():
			return fmt.Errorf("acquiring lock context canceled: %w", l.config.ctx.Err())
		default:
			// Acquire the lock.
			err := unix.FcntlFlock(l.file.Fd(), unix.F_SETLK, &lock)
			if err == nil {
				return nil
			}
			// Sleep for a while for the next retry.
			time.Sleep(randomDuration(minWaitDuration, maxWaitDuration))
		}
	}
}

func (l *FileLock) acquireLockWait() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create a Flock_t structure.
	lock := unix.Flock_t{
		Type:   unix.F_WRLCK, // write lock
		Whence: 0,            // relative to the start of the file
		Start:  0,            // lock starts at byte 0
		Len:    0,            // lock extends to EOF
	}

	// Start a goroutine to enforce the timeout.
	timeoutC := time.After(l.config.timeout)

	destroyC := make(chan struct{})
	errC := make(chan error, 1)
	go func() {
		// Wait until acquire the lock.
		errC <- unix.FcntlFlock(l.file.Fd(), unix.F_SETLKW, &lock)
		select {
		case <-destroyC:
			// Immediately release the lock after the lock been acquired.
			l.Unlock()
		default:
			close(destroyC)
		}
	}()

	select {
	case <-l.config.ctx.Done():
		close(destroyC)
		return fmt.Errorf("acquiring lock context canceled: %w", l.config.ctx.Err())
	case <-timeoutC:
		close(destroyC)
		return ErrTimeout
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("acquiring lock: %w", err)
		}
		return nil
	}
}

func randomDuration(minDuration, maxDuration time.Duration) time.Duration {
	return minDuration + time.Duration(randomGenerator.Int63n(int64(maxDuration-minDuration)))
}
