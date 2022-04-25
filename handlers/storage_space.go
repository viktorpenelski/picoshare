package handlers

import (
	"net/http"
	"github.com/mtlynch/picoshare/v2/store"
	"encoding/json"
	"math"
	//"fmt"
)

	
type StorageSpace struct {
    All   float64 `json:"all"`
    Free  float64 `json:"free"`
    Used  float64 `json:"used"`
	UsedPercentage float64 `json:"usedPercentage"`
}

const (
    B  = 1
    KB = 1024 * B
    MB = 1024 * KB
    GB = 1024 * MB
)


func (s Server) getStorageSpace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbPath := s.store.GetDbPath()
		storageSpaceInBytes, err := store.FromDbPath(dbPath)
		if err != nil {
			http.Error(w, "Could not fetch storage space info, likely unsupported OS.", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		storageSpaceInGB := storageSpaceBytesToGb(*storageSpaceInBytes)
		if err := json.NewEncoder(w).Encode(storageSpaceInGB); err != nil {
			panic(err)
		}
	}
}

func storageSpaceBytesToGb(storageSpaceInBytes store.StorageSpaceInBytes) *StorageSpace {
	// return with double precision
	return &StorageSpace{
		All: doublePrecision(float64(storageSpaceInBytes.All)/float64(GB)),
		Free: doublePrecision(float64(storageSpaceInBytes.Free)/float64(GB)),
		Used: doublePrecision(float64(storageSpaceInBytes.All - storageSpaceInBytes.Free)/float64(GB)),
		UsedPercentage: doublePrecision(float64(storageSpaceInBytes.All - storageSpaceInBytes.Free)/float64(storageSpaceInBytes.All) * 100),
	}
}

func doublePrecision(number float64) float64 {
	// TODO(vik) should we just leave this to the FE?
	return math.Round(number * 100) / 100
}
