package mmap

import (
	"os"
	"runtime"
	"sync"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

type mapContext struct {
	f      *os.File
	opened bool
}

var mmapMu sync.Mutex
var mmapFiles map[unsafe.Pointer]mapContext

// MapFile maps an opened file to a byte slice of data.
func MapFile(f *os.File, writable bool) (data []byte, err error) {
	prot := mmap.RDONLY
	if writable {
		prot = mmap.RDWR
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		return nil, nil
	}
	m, err := mmap.Map(f, prot, 0)
	if err != nil {
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
			mmapFiles = make(map[unsafe.Pointer]mapContext)
		}
		mmapFiles[unsafe.Pointer(&m[0])] = mapContext{f, false}
		mmapMu.Unlock()
	}
	return []byte(m), nil
}

// Open will mmap a file to a byte slice of data.
func Open(path string, writable bool) (data []byte, err error) {
	flag := os.O_RDONLY
	if writable {
		flag = os.O_RDWR
	}
	f, err := os.OpenFile(path, flag, 0)
	if err != nil {
		return nil, err
	}
	m, err := MapFile(f, writable)
	if err != nil {
		f.Close()
		return nil, err
	}
	if len(m) == 0 {
		// Empty file. Release map
		Close(m)
		f.Close()
		return nil, nil
	}
	if runtime.GOOS == "windows" {
		mmapMu.Lock()
		ctx := mmapFiles[unsafe.Pointer(&m[0])]
		ctx.opened = true
		mmapFiles[unsafe.Pointer(&m[0])] = ctx
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
		var ctx mapContext
		var ok bool
		mmapMu.Lock()
		ctx, ok = mmapFiles[unsafe.Pointer(&data[0])]
		if ok {
			delete(mmapFiles, unsafe.Pointer(&data[0]))
		}
		mmapMu.Unlock()
		if ok && ctx.opened {
			ctx.f.Close()
		}
	}
	m := mmap.MMap(data)
	return m.Unmap()
}

// Create a new mmap file with the provided size
func Create(path string, size int) ([]byte, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if _, err := f.WriteAt([]byte{0}, int64(size)-1); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return Open(path, true)
}
