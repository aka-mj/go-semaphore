package semaphore

import (
	"testing"
)

func Test_SemGetValue(t *testing.T) {
	var sem Semaphore
	if err := sem.Open("/testsem", 0644, 1); err != nil {
		t.Fatalf("Failed to open: %v", err)
	}

	val, err := sem.GetValue()
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	if val != 1 {
		t.Error("Value incorrect")
	}

	if err := sem.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
}

func Test_longName(t *testing.T) {
	var sem Semaphore
	name := make([]byte, 256)
	name[0] = '/'
	for i := 1; i < 256; i++ {
		name[i] = 'X'
	}
	if err := sem.Open(string(name), 0644, 1); err == nil {
		t.Fatalf("Failed to open: %v", err)
	}
}
