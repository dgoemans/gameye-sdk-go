package test

import (
	"encoding/json"
	"net/http"
)

type QueryPatch struct {
	Path  []string    `json:"path"`
	Value interface{} `json:"value"`
}

func CreateApiTestServerMux(
	state map[string]interface{},
	patchChannel chan QueryPatch,
) (
	mux *http.ServeMux,
) {
	handleNoop := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	handleAction := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}

	handleFetch := func(w http.ResponseWriter, r *http.Request) {
		var err error
		accept := r.Header.Get("Accept")
		switch accept {
		case "application/json":
			encoder := json.NewEncoder(w)
			err = encoder.Encode(state)
			if err != nil {
				panic(err)
			}

		case "application/x-ndjson":
			flusher := w.(http.Flusher)
			closeNotifier := w.(http.CloseNotifier)
			closeChannel := closeNotifier.CloseNotify()
			w.Header().Add("Transfer-Encoding", "chunked")
			w.WriteHeader(http.StatusOK)
			flusher.Flush()

			encoder := json.NewEncoder(w)
			for {
				select {
				case <-closeChannel:
					break

				case patch := <-patchChannel:
					err = encoder.Encode([]QueryPatch{patch})
					if err != nil {
						panic(err)
					}
					flusher.Flush()
				}
			}
		}
	}

	mux = http.NewServeMux()
	mux.HandleFunc("/noop", handleNoop)
	mux.HandleFunc("/action/noop", handleAction)
	mux.HandleFunc("/fetch/noop", handleFetch)

	return
}
