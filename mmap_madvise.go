//go:build ignore

package mmap

import (
	"os"
	"syscall"
	"unsafe"
)

type AdviseFlag int

const (
	Normal     AdviseFlag = 0
	Sequential AdviseFlag = 1
	Random     AdviseFlag = 2
)

func Advise(data []byte, flag AdviseFlag) error {
	if len(data) == 0 {
		return os.ErrInvalid
	}
	var mflag int
	switch flag {
	case Normal:
		mflag = syscall.MADV_NORMAL
	case Sequential:
		mflag = syscall.MADV_SEQUENTIAL
	case Random:
		mflag = syscall.MADV_RANDOM
	default:
		return os.ErrInvalid
	}
	_, _, err := syscall.Syscall(
		uintptr(syscall.SYS_MADVISE),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(mflag),
		0,
	)
	if err != 0 {
		return err
	}
	return nil
}
