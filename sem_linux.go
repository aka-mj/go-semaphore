package semaphore


// #cgo LDFLAGS: -pthread
// #include <stdlib.h>
// #include <fcntl.h>
// #include <sys/stat.h>
// #include <sys/types.h>
// #include <semaphore.h>
// #include <time.h>
import "C"

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
