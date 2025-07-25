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

func Test_SemTimedWait_nowait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem", 0644, 1)
	if err := sem.TimedWait(1 * time.Second); err != nil {
		t.Fatalf("Should not have timedout: %v", err)
	}
	sem.Close()
}

func Test_SemTimedWait_wait(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem_wait", 0644, 1)
	sem.Wait()
	v, _ := sem.GetValue()
	fmt.Println(v)

	end := make(chan error, 1)
	go func() {
		var sem2 Semaphore
		sem2.Open("/testsem_wait", 0644, 1)
		end <- sem2.TimedWait(2 * time.Second)
		sem2.Close()
	}()

	time.Sleep(500 * time.Millisecond)
	sem.Post()
	err := <-end
	sem.Close()
	sem.Unlink()
	if err != nil {
		t.Fatalf("Should not have timedout: %v", err)
	}
}

// When doing a TimedWait any release of the semaphore should allow the
// waiting process to grab the semaphore.
func Test_SemTimedWait_wait_with_short_gap(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem_wait2", 0644, 1)
	sem.Wait()
	v, _ := sem.GetValue()
	fmt.Println(v)

	end := make(chan error, 1)
	semWait := make(chan time.Duration, 1)
	go func() {
		t1 := time.Now()
		var sem2 Semaphore
		sem2.Open("/testsem_wait2", 0644, 1)
		end <- sem2.TimedWait(2 * time.Second)
		semWait <- time.Since(t1)
		sem2.Post()
	}()

	time.Sleep(500 * time.Millisecond)
	sem.Post()
	// If you do an immediate sem.Wait call you wouldn't allow another process to take the semaphore
	// by doing a sleep even if it is a short sleep the semaphore can be taken by another process
	time.Sleep(1 * time.Millisecond)
	sem.Wait()
	time.Sleep(1600 * time.Millisecond)
	sem.Post()
	err := <-end
	semWaitDuration := <-semWait
	if semWaitDuration < 500*time.Millisecond {
		t.Fatalf("Impossible we still slept for 500 Milliseconds prior to releasing semaphore: sem wait was %s", semWaitDuration.String())
	}
	if semWaitDuration > 550*time.Millisecond {
		t.Fatalf("It took too long to acquire semaphore. We still slept for 500 Milliseconds prior to releasing semaphore: sem wait was %s", semWaitDuration.String())
	}

	sem.Close()
	sem.Unlink()
	if err != nil {
		t.Fatalf("Should not have timedout: %v", err)
	}
}

func Test_SemTimedWait_timeout(t *testing.T) {
	var sem Semaphore
	sem.Open("/testsem_wait", 0644, 1)
	sem.Wait()
	v, _ := sem.GetValue()
	fmt.Println(v)

	end := make(chan error, 1)
	go func() {
		var sem2 Semaphore
		sem2.Open("/testsem_wait", 0644, 1)
		end <- sem2.TimedWait(1 * time.Second)
		sem2.Close()
	}()

	time.Sleep(2 * time.Second)
	sem.Post()
	err := <-end
	sem.Close()
	sem.Unlink()
	if err == nil {
		t.Fatalf("Should have timedout: %v", err)
	}
}

func Test_SemDoubleClose(t *testing.T) {
	var sem Semaphore
	if err := sem.Open("/testsem", 0644, 1); err != nil {
		t.Fatalf("Failed to open: %v", err)
	}

	if err := sem.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
	if err := sem.Close(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
}

func Test_SemDoubleUnlink(t *testing.T) {
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
	if err := sem.Unlink(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
}

func Test_isSemaphoreInitialized(t *testing.T) {
	var sem Semaphore
	if err := sem.Close(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
	if _, err := sem.GetValue(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
	if err := sem.Post(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
	if err := sem.Wait(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
	if err := sem.TryWait(); err == nil {
		t.Fatalf("Should have received error: %v", err)
	}
}
