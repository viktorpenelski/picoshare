package store

import (
	syscall "golang.org/x/sys/unix"
	"path"
	"runtime"
	"errors"
)

type StorageSpaceInBytes struct {
	All   uint64
	Free  uint64
}

func unixStorageSpaceFrom(path string) (*StorageSpaceInBytes, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}
	disk := &StorageSpaceInBytes {
		All:   fs.Blocks * uint64(fs.Bsize),
		Free:  fs.Bfree * uint64(fs.Bsize),
	}
	return disk, nil
}


func FromDbPath(dbPath string) (*StorageSpaceInBytes, error) {
	// only tested on a linux distribution so far. 
	// Maybe the same implementation works on other unix OSs as well?
	if runtime.GOOS != "linux" {
		return nil, errors.New("Storage space not supported on " + runtime.GOOS)
	}
	dir := path.Dir(dbPath)
	return unixStorageSpaceFrom(dir)
}
