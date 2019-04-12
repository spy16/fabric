package server

import (
	"encoding/json"
	"net/http"

	"github.com/spy16/fabric"
)

// NewHTTP initializes the an http router with all the query routes and
// middlewares initialized.
func NewHTTP(fab *fabric.Fabric) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", queryHandler(fab))
	return mux
}

func queryHandler(fab *fabric.Fabric) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			wr.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(wr).Encode(map[string]string{
				"error": "method not allowed, use post",
			})
			return
		}

		var query fabric.Query
		if err := json.NewDecoder(req.Body).Decode(&query); err != nil {
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(map[string]string{
				"error": "not valid json body",
			})
			return
		}

		tri, err := fab.Query(req.Context(), query)
		if err != nil {
			wr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(wr).Encode(map[string]string{
				"error": err.Error(),
			})
		}

		json.NewEncoder(wr).Encode(tri)
	}
}
