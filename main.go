package main

import (
	
	// this works the same using just "syscall", but https://pkg.go.dev/syscall state that it is deprecated
	syscall "golang.org/x/sys/unix"

	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"runtime"

	gorilla "github.com/mtlynch/gorilla-handlers"

	"github.com/mtlynch/picoshare/v2/garbagecollect"
	"github.com/mtlynch/picoshare/v2/handlers"
	"github.com/mtlynch/picoshare/v2/handlers/auth/shared_secret"
	"github.com/mtlynch/picoshare/v2/store/sqlite"
)

type DiskStatus struct {
    All   uint64 `json:"all"`
    Used  uint64 `json:"used"`
    Free  uint64 `json:"free"`
    Avail uint64 `json:"avail"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
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

func printCurrentlyRunningOs() {
	fmt.Printf("Running on %s %s\n", runtime.GOOS, runtime.GOARCH)
}

func main() {
	dbPath := flag.String("db", "data/store.db", "path to database")
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("failed to retrieve disk space: %v", err)
		return
	}
	diskSpace(wd)
	diskSpace(*dbPath)
	printCurrentlyRunningOs()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("Starting picoshare server")


	flag.Parse()

	authenticator, err := shared_secret.New(requireEnv("PS_SHARED_SECRET"))
	if err != nil {
		log.Fatalf("invalid shared secret: %v", err)
	}

	ensureDirExists(filepath.Dir(*dbPath))

	store := sqlite.New(*dbPath)

	gc := garbagecollect.NewScheduler(store, 7*time.Hour)
	gc.StartAsync()

	h := gorilla.LoggingHandler(os.Stdout, handlers.New(authenticator, store).Router())
	if os.Getenv("PS_BEHIND_PROXY") != "" {
		h = gorilla.ProxyIPHeadersHandler(h)
	}
	http.Handle("/", h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}
	log.Printf("Listening on %s", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("missing required environment variable: %s", key))
	}
	return val
}

func ensureDirExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
