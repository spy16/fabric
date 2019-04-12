package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/spy16/fabric"
)

// NewHTTP initializes the an http router with all the query routes and
// middlewares initialized.
func NewHTTP(fab *fabric.Fabric) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", queryHandler(fab))
	mux.HandleFunc("/query.js", queryHandler(fab))
	return mux
}

func queryHandler(fab *fabric.Fabric) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		query, err := readQuery(req.URL.Query())
		if err != nil {
			writeResponse(wr, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		tri, err := fab.Query(req.Context(), *query)
		if err != nil {
			writeResponse(wr, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeResponse(wr, http.StatusOK, tri)
	}
}

func writeResponse(wr http.ResponseWriter, status int, body interface{}) {
	wr.Header().Set("Content-Type", "application/json; charset=utf-8")
	wr.WriteHeader(http.StatusOK)
	json.NewEncoder(wr).Encode(body)
}

func readQuery(vals url.Values) (*fabric.Query, error) {
	var q fabric.Query
	if err := readInto(vals, "source", &q.Source); err != nil {
		return nil, err
	}

	if err := readInto(vals, "predicate", &q.Predicate); err != nil {
		return nil, err
	}

	if err := readInto(vals, "target", &q.Target); err != nil {
		return nil, err
	}

	return &q, nil
}

func readInto(vals url.Values, name string, cl *fabric.Clause) error {
	parts := strings.Fields(vals.Get(name))
	if len(parts) == 0 {
		return nil
	}

	if len(parts) != 2 {
		return fmt.Errorf("invalid %s clause", name)
	}

	cl.Type = parts[0]
	cl.Value = parts[1]
	return nil
}
