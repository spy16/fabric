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
	handleQuery := queryHandler(fab)
	handleInsert := insertHandler(fab)
	handleReWeight := reweightHandler(fab)
	handleDelete := deleteHandler(fab)

	mux := http.NewServeMux()
	mux.HandleFunc("/triples", func(wr http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			handleQuery(wr, req)

		case http.MethodPost:
			handleInsert(wr, req)

		case http.MethodPatch:
			handleReWeight(wr, req)

		case http.MethodDelete:
			handleDelete(wr, req)

		default:
			writeResponse(wr, req, http.StatusMethodNotAllowed, map[string]string{
				"error": "method not allowed",
			})
		}
	})
	return withLogs(mux)
}

func queryHandler(fab *fabric.Fabric) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		query, err := readQuery(req.URL.Query())
		if err != nil {
			writeResponse(wr, req, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		tri, err := fab.Query(req.Context(), *query)
		if err != nil {
			writeResponse(wr, req, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeTriples(wr, req, http.StatusOK, tri)
	}
}

func insertHandler(fab *fabric.Fabric) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		var tri fabric.Triple
		if err := json.NewDecoder(req.Body).Decode(&tri); err != nil {
			writeResponse(wr, req, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if err := tri.Validate(); err != nil {
			writeResponse(wr, req, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if err := fab.Insert(req.Context(), tri); err != nil {
			writeResponse(wr, req, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		writeResponse(wr, req, http.StatusCreated, nil)
	}
}

func reweightHandler(fab *fabric.Fabric) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		var payload struct {
			fabric.Query `json:",inline"`

			Delta   float64 `json:"delta"`
			Replace bool    `json:"replace"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			writeResponse(wr, req, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		updates, err := fab.ReWeight(req.Context(), payload.Query, payload.Delta, payload.Replace)
		if err != nil {
			writeResponse(wr, req, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeResponse(wr, req, http.StatusOK, map[string]interface{}{
			"updated": updates,
		})
	}
}

func deleteHandler(fab *fabric.Fabric) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		query, err := readQuery(req.URL.Query())
		if err != nil {
			writeResponse(wr, req, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		deleted, err := fab.Delete(req.Context(), *query)
		if err != nil {
			writeResponse(wr, req, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeResponse(wr, req, http.StatusOK, map[string]interface{}{
			"deleted": deleted,
		})
	}
}

func writeTriples(wr http.ResponseWriter, req *http.Request, status int, triples []fabric.Triple) {
	switch outputFormat(req) {
	case "dot":
		wr.Write([]byte(fabric.ExportDOT("fabric", triples)))

	case "plot":
		plotTemplate.Execute(wr, map[string]interface{}{
			"graphVizStr": "`" + fabric.ExportDOT("fabric", triples) + "`",
		})

	default:
		writeResponse(wr, req, status, triples)
	}
}

func writeResponse(wr http.ResponseWriter, req *http.Request, status int, body interface{}) {
	wr.Header().Set("Content-Type", "application/json; charset=utf-8")
	wr.WriteHeader(status)
	if body == nil || status == http.StatusNoContent {
		return
	}

	switch outputFormat(req) {
	default:
		json.NewEncoder(wr).Encode(body)
	}
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

func outputFormat(req *http.Request) string {
	f := strings.TrimSpace(req.URL.Query().Get("format"))
	if f != "" {
		return f
	}

	return "json"
}
