package reception

import (
	"encoding/json"
	"net/http"

	"github.com/martengine/reception/service"
)

// Server controls the requests and forwards them to relevant service.
func Server() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if ResponseWriter supports Flusher.
		if flusher, ok := w.(http.Flusher); ok {
			streamResponse(w, r, flusher)
			return
		}
	}
}

func responses(r *http.Request) <-chan service.Response {
	responses := make(chan service.Response)
	go requestServices(r, responses)
	return responses
}

func requestServices(r *http.Request, responses chan<- service.Response) {
	// TODO
}

func streamResponse(w http.ResponseWriter, r *http.Request, flusher http.Flusher) {
	for response := range responses(r) {
		json.NewEncoder(w).Encode(response)
		flusher.Flush()
	}
}

// Info renders the list of public registered services with all available information.
func Info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := struct{ Services []service.Service }{Services: PublicServices()}
		body, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.Write(body)
	}
}
