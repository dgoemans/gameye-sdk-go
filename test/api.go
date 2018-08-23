package test

import (
	"fmt"
	"net/http"
)

/*
CreateAPITestServerMux creates the ServeMux for a api test server
*/
func CreateAPITestServerMux(
	state string,
	patchChannel chan string,
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
			_, err = fmt.Fprintln(w, state)
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

			for {
				select {
				case <-closeChannel:
					return

				case patch := <-patchChannel:
					_, err = fmt.Fprintln(w, patch)
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
	mux.HandleFunc("/action/start-match", handleAction)
	mux.HandleFunc("/action/stop-match", handleAction)

	mux.HandleFunc("/fetch/noop", handleFetch)
	mux.HandleFunc("/fetch/game", handleFetch)
	mux.HandleFunc("/fetch/match", handleFetch)
	mux.HandleFunc("/fetch/statistic", handleFetch)
	mux.HandleFunc("/fetch/template", handleFetch)

	return
}