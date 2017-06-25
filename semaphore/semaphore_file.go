package semaphore

import (
	"syscall"
	"os"
	"io"
)

type SemaphoreFile struct {
	path string
	max int
}

func NewSemaphoreFile(path string, max int) *SemaphoreFile {
	s := SemaphoreFile{path: path, max: max}
	return &s
}

func flock(f *os.File) error {
	fd := f.Fd()
	lock := new(syscall.Flock_t)
	lock.Type = syscall.F_WRLCK
	err := syscall.FcntlFlock(fd, syscall.F_SETLK, lock)
	return err
}

func unflock(f *os.File) error {
	fd := f.Fd()
	lock := new(syscall.Flock_t)
	lock.Type = syscall.F_UNLCK
	err := syscall.FcntlFlock(fd, syscall.F_SETLK, lock)
	return err
}

func (s *SemaphoreFile) operation(op func (*SemaphoreData,string) (bool,error), id string) (bool,error) {
	f,err := os.OpenFile(s.path, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}

	defer f.Close()
	if err:= flock(f); err != nil {
		return false, err
	}
	defer unflock(f)

	sem_data,err := Load(f)
	if err == io.EOF {
		sem_data = New(s.max)
	} else if err != nil {
		return false, err
	}

	result,err := op(sem_data, id)
	if !result {
		return false, err
	}

	if _, err := f.Seek(0,0); err != nil {
		return false, err
	}

	if err := f.Truncate(0); err != nil {
		return false, err
	}

	json_repr, err := sem_data.ToJSON()
	if err != nil {
		return false, err
	}

	if _, err := f.WriteString(json_repr); err != nil {
		return false, err
	}

	return true, nil
}

// Open semaphore file for exclusive access, write
func (s *SemaphoreFile) Acquire(id string) (bool,error) {
	return s.operation((*SemaphoreData).Acquire, id)
}

func (s *SemaphoreFile) Release(id string) (bool,error) {
	return s.operation((*SemaphoreData).Release, id)
}

func (s *SemaphoreFile) Holds(id string) (bool,error) {
	return s.operation((*SemaphoreData).Holds, id)
}
