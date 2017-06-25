package semaphore

import (
	"os"
	"testing"
	"io/ioutil"
)

func TestFileAcquire(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}

	defer os.Remove(tmpfile.Name())
	if err := tmpfile.Close(); err != nil {
		t.Error(err)
	}

	sf := NewSemaphoreFile(tmpfile.Name(), 1)
	ok, err := sf.Acquire("test")
	if !ok {
		t.Error("Failed to acquire semaphore")
		t.Error(err)
	}
	ok, err = sf.Acquire("test2")
	if ok {
		t.Error("Able to acquire semaphore more than once")
		t.Error(err)
	}
}
