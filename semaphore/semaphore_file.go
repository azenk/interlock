package semaphore

import (
	"syscall"
	"os"
)

type SemaphoreFile struct {
	path string
}

func NewSemaphoreFile(path string) *SemaphoreFile {
	s := SemaphoreFile{path: path}
	return &s
}

func flock(f *os.File) error {
	fd := f.Fd()
	err := syscall.Flock(int(fd), syscall.LOCK_EX | syscall.LOCK_NB)
	return err 
}

func unflock(f *os.File) error {
	fd := f.Fd()
	err := syscall.Flock(int(fd), syscall.LOCK_UN | syscall.LOCK_NB)
	return err
}

// Open semaphore file for exclusive access, write
func (s *SemaphoreFile) Acquire(id string) bool {
	return false
}
	
