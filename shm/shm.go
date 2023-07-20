package shm

import (
	"fmt"
	"os"
	"syscall"
)

func Alloc(filepath string, size int) ([]byte, error) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("filepath error: %s", err)
	}
	defer f.Close()

	b, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap(%s) size(%v) error: %s", filepath, size, err)
	}
	if b == nil {
		return nil, fmt.Errorf("error: mmap(%s) data is nil", filepath)
	}
	return b, nil
}
