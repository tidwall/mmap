package mmap

import (
	"os"
	"runtime"
	"sync"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

var mmapMu sync.Mutex
var mmapFiles map[unsafe.Pointer]*os.File

// Open will mmap a file to a byte slice of data.
func Open(path string, writable bool) (data []byte, err error) {
	flag, prot := os.O_RDONLY, mmap.RDONLY
	if writable {
		flag, prot = os.O_RDWR, mmap.RDWR
	}
	f, err := os.OpenFile(path, flag, 0)
	if err != nil {
		return nil, err
	}
	m, err := mmap.Map(f, prot, 0)
	if err != nil {
		f.Close()
		return nil, err
	}
	if len(m) == 0 {
		// Empty file. Release map
		m.Unmap()
		f.Close()
		return nil, nil
	}
	if runtime.GOOS == "windows" {
		// Keep track of the file.
		mmapMu.Lock()
		if mmapFiles == nil {
			mmapFiles = make(map[unsafe.Pointer]*os.File)
		}
		mmapFiles[unsafe.Pointer(&m[0])] = f
		mmapMu.Unlock()
	} else {
		// Allowed to close the file.
		f.Close()
	}
	return []byte(m), nil
}

// Close releases the data.
func Close(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if runtime.GOOS == "windows" {
		// Close file first
		mmapMu.Lock()
		f, ok := mmapFiles[unsafe.Pointer(&data[0])]
		if ok {
			delete(mmapFiles, unsafe.Pointer(&data[0]))
		}
		mmapMu.Unlock()
		if f != nil {
			f.Close()
		}
	}
	m := mmap.MMap(data)
	return m.Unmap()
}
