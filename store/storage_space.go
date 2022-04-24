package store

import (
	syscall "golang.org/x/sys/unix"
	"fmt"
)

type StorageSpace struct {
    All   uint64 `json:"all"`
    Used  uint64 `json:"used"`
    Free  uint64 `json:"free"`
    Avail uint64 `json:"avail"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk StorageSpace) {
    fs := syscall.Statfs_t{}
    err := syscall.Statfs(path, &fs)
    if err != nil {
        return
    }
    disk.All = fs.Blocks * uint64(fs.Bsize)
    disk.Avail = fs.Bavail * uint64(fs.Bsize)
    disk.Free = fs.Bfree * uint64(fs.Bsize)
    disk.Used = disk.All - disk.Free
    return
}

const (
    B  = 1
    KB = 1024 * B
    MB = 1024 * KB
    GB = 1024 * MB
)

func diskSpace2(path string) {
    disk := DiskUsage(path)
    fmt.Println("")
    fmt.Println(path, ":")
    fmt.Printf("Storage: %.2f GB\n", float64(disk.All)/float64(GB))
    fmt.Printf("Avail: %.2f GB\n", float64(disk.Avail)/float64(GB))
    fmt.Printf("Used: %.2f GB\n", float64(disk.Used)/float64(GB))
}

func diskSpace(path string) {
	var stat syscall.Statfs_t
	log.Printf(path)
	diskSpace2(path)
	syscall.Statfs(path, &stat)
	log.Printf("disk space: %d bytes free", stat.Bavail*uint64(stat.Bsize))
}