// Project:         C1515
// Author:          Michael John
// Copyright:       2014-2016 Crown Equipment Corp.

// Package semaphore provides an interface to named userspace semaphores.
package semaphore

import (
	"errors"
	"syscall"
	"time"
	"unsafe"
)

// #cgo LDFLAGS: -pthread
// #include <stdlib.h>
// #include <fcntl.h>
// #include <sys/stat.h>
// #include <sys/types.h>
// #include <semaphore.h>
// #include <time.h>
// sem_t* Go_sem_open(const char *name, int oflag, mode_t mode, unsigned int value)
// {
//		return sem_open(name, oflag, mode, value);
// }
import "C"

type Semaphore struct {
	sem  *C.sem_t //semaphore returned by sem_open
	name string   //name of semaphore
}

func (s *Semaphore) isSemaphoreInitialized() (bool, error) {
	if s.sem == nil {
		return false, errors.New("Not a valid semaphore")
	}
	return true, nil
}

// Open creates a new POSIX semaphore or opens an existing semaphore.
// The semaphore is identified by name. The mode argument specifies the permissions
// to be placed on the new semaphore. The value argument specifies the initial
// value for the new semaphore. If the named semaphore already exist, mode and
// value are ignored.
// For details see sem_overview(7).
func (s *Semaphore) Open(name string, mode, value uint32) error {
	s.name = name
	n := C.CString(name)

	var err error
	s.sem, err = C.Go_sem_open(n, syscall.O_CREAT, C.mode_t(mode), C.uint(value))
	C.free(unsafe.Pointer(n))
	if s.sem == nil {
		return err
	}
	return nil
}

// Close closes the named semaphore, allowing any resources that the system has
// allocated to the calling process for this semaphore to be freed.
func (s *Semaphore) Close() error {
	if ok, err := s.isSemaphoreInitialized(); !ok {
		return err
	}

	ret, err := C.sem_close(s.sem)
	if ret != 0 {
		return err
	}
	s.sem = nil
	return nil
}

// GetValue returns the current value of the semaphore.
func (s *Semaphore) GetValue() (int, error) {
	if ok, err := s.isSemaphoreInitialized(); !ok {
		return 0, err
	}

	var val C.int
	ret, err := C.sem_getvalue(s.sem, &val)
	if ret != 0 {
		return int(ret), err
	}
	return int(val), nil
}

// Post increments the semaphore.
func (s *Semaphore) Post() error {
	if ok, err := s.isSemaphoreInitialized(); !ok {
		return err
	}

	ret, err := C.sem_post(s.sem)
	if ret != 0 {
		return err
	}
	return nil
}

// Wait decrements the semaphore. If the semaphore's value is greater than zero,
// then the decrement proceeds, and the function returns, immediately. If the
// semaphore currently has the value zero, then the call blocks until either
// it becomes possible to perform the decrement, or a signal interrupts the call.
func (s *Semaphore) Wait() error {
	if ok, err := s.isSemaphoreInitialized(); !ok {
		return err
	}

	ret, err := C.sem_wait(s.sem)
	if ret != 0 {
		return err
	}
	return nil
}

// TryWait is the same as Wait(), except that if the decrement cannot be immediately
// performed, then the call returns an error instead of blocking.
func (s *Semaphore) TryWait() error {
	if ok, err := s.isSemaphoreInitialized(); !ok {
		return err
	}

	ret, err := C.sem_trywait(s.sem)
	if ret != 0 {
		return err
	}
	return nil
}

// TimedWait is the same as Wait(), except that d specifies a limit on the
// amount of time that the call should block if the decrement cannot be
// immediately performed.
func (s *Semaphore) TimedWait(d time.Duration) error {
	if err := s.TryWait(); err == nil {
		// success
		return nil
	}
	time.Sleep(d)
	if err := s.TryWait(); err == nil {
		// success
		return nil
	}
	return errors.New("The call timed out before the semaphore could be locked")
}

// Unlink removes the named semaphore. The semaphore name is removed immediately.
// The semaphore is destroyed once all other processes that have the semaphore
// open close it.
func (s *Semaphore) Unlink() error {
	name := C.CString(s.name)
	ret, err := C.sem_unlink(name)
	C.free(unsafe.Pointer(name))
	if ret != 0 {
		return err
	}
	return nil
}
