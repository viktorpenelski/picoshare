package store

import (
	syscall "golang.org/x/sys/unix"
	"path"
	"runtime"
	"errors"
)

type storageSpaceBytes struct {
	All   uint64
	Free  uint64
}
	
type StorageSpaceGB struct {
    All   float64
    Free  float64
    Used  float64
	UsedPercentage float64
}

const (
    B  = 1
    KB = 1024 * B
    MB = 1024 * KB
    GB = 1024 * MB
)

func unixStorageSpaceFrom(path string) (*storageSpaceBytes, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}
	disk := &storageSpaceBytes {
		All:   fs.Blocks * uint64(fs.Bsize),
		Free:  fs.Bfree * uint64(fs.Bsize),
	}
	return disk, nil
}

func storageSpaceBytesToGb(bytes storageSpaceBytes) *StorageSpaceGB {
	return &StorageSpaceGB {
		All:   float64(bytes.All)/float64(GB),
		Free:  float64(bytes.Free)/float64(GB),
		Used:  float64(bytes.All - bytes.Free)/float64(GB),
		UsedPercentage: (float64(bytes.All - bytes.Free)/float64(bytes.All)) * 100,
	}
}

func FromDbPath(dbPath string) (*StorageSpaceGB, error) {
	// only tested on a linux distribution so far. 
	// Maybe the same implementation works on other unix OSs as well?
	if runtime.GOOS != "linux" {
		return nil, errors.New("storage space not supported on " + runtime.GOOS)
	}
	dir := path.Dir(dbPath)
	storageSpaceBytes, err := unixStorageSpaceFrom(dir)
	if err != nil {
		return nil, err
	}
	return storageSpaceBytesToGb(*storageSpaceBytes), nil
}
