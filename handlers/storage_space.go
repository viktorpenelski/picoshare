package handlers

import (
	"net/http"
	"github.com/mtlynch/picoshare/v2/store"
	"encoding/json"
)

func (s Server) getStorageSpace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbPath := s.store.GetDbPath()
		storageSpace, err := store.FromDbPath(dbPath)
		if err != nil {
			// TODO(vik): should we log this? It might polute the logs on non-linux OS where it is expected to fail
			// log.Printf("Could not get the system storage space", err) 
			http.Error(w, "Storage space not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(storageSpace); err != nil {
			panic(err)
		}
	}
}
