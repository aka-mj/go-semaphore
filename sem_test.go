package semaphore

import (
	"fmt"
	"testing"
	"time"
)

func Test_Semaphore(t *testing.T) {
	var sem Semaphore
	if err := sem.Open("/testsem", 0644, 1); err != nil {
		t.Fatalf("Failed to open: %v", err)
	}

	if err := sem.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}

	if err := sem.Unlink(); err != nil {
		t.Fatalf("Failed to unlink: %v", err)
	}
}

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

func Test_SemWait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)

	if err := sem.Wait(); err != nil {
		t.Fatalf("Failed to wait: %v", err)
	}
	val, _ := sem.GetValue()
	if val != 0 {
		t.Fatal("Value incorrect")
	}

	if err := sem.Post(); err != nil {
		t.Fatalf("Failed to post: %v", err)
	}
	sem.Close()
}

func Test_SemTryWait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)

	if err := sem.TryWait(); err != nil {
		t.Fatalf("Failed to wait: %v", err)
	}

	sem.Post()
	sem.Close()
}

func Test_SemTryWaitFail(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)

	sem.Wait()
	err := sem.TryWait()
	if err == nil {
		t.Fatal("TryWait should have failed")
	}
	sem.Post()
	sem.Close()
}

func Test_SemPostFail(t *testing.T) {
	var sem Semaphore
	err := sem.Post()
	if err == nil {
		t.Fatal("Post should have failed")
	}
}

func Test_SemTimedWait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)
	sem.Wait()

	go func() {
		var sem2 Semaphore
		sem2.Open("/testsem", 0644, 1)
		i, _ := sem2.GetValue()
		fmt.Printf(">> %v\n", i)
		err := sem2.TimedWait(11 * time.Second)
		if err != nil {
			t.Fatalf("Failed: %v", err)
		} else {
			sem2.Post()
		}
		sem2.Close()
	}()

	time.Sleep(3 * time.Second)
	sem.Post()
	sem.Close()
}
