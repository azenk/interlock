package semaphore

import (
	"os"
	"time"
	"testing"
	"io/ioutil"
)

func TestFlock(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	lockProc := func (f *os.File, ch chan bool) {
		err := flock(f)
		if err != nil {
			t.Error(err)
			ch <- false
		}
		d, _ := time.ParseDuration("10ms")
		time.Sleep(d)
		ch <- true
	}
	
	ch := make(chan bool)
	go lockProc(tmpfile, ch)
	d, _ := time.ParseDuration("5ms")
	time.Sleep(d)
	go lockProc(tmpfile, ch)

	a := <-ch
	b := <-ch

	switch {
	case a && b:
		t.Error("Locking file succeeded twice")
	case !(a || b):
		t.Error("Locking file failed twice")
	}

	if err := unflock(tmpfile); err != nil {
		t.Error("Can't unflock file")
	}

	if err := tmpfile.Close(); err != nil {
		t.Error(err)
	}

}
