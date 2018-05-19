package reception

import (
	"encoding/json"
	"net/http"

	"github.com/martengine/reception/service"
)

// Server controls the requests and forwards them to relevant service.
func Server() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
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
