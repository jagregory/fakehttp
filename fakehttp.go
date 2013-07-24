package fakehttp

import (
	"fmt"
	"net"
	"net/http"
)

var servers = make(map[string]*http.Server)

// Listen on a port, and respond with canned responses.
// Response can either be a string or []byte for 200 responses,
// or an int representing the status code.
// Returns a net.Listener, always Close it before your test exits.
func Listen(port int, routes map[string]*Response) net.Listener {
	addr := fmt.Sprintf(":%d", port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err.Error())
	}

	srv := servers[addr]

	if srv == nil {
		srv = &http.Server{}
		servers[addr] = srv
	}

	srv.Handler = &fakeHandler{routes}

	go srv.Serve(l)

	return l
}

type Response struct {
	ContentOrStatusCode interface{}
}

type fakeHandler struct {
	routes map[string]*Response
}

func (h *fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h.routes[r.URL.Path]
	if response == nil {
		http.NotFound(w, r)
		return
	}

	switch val := response.ContentOrStatusCode.(type) {
	case string:
		fmt.Fprint(w, val)
	case []byte:
		w.Write(val)
	case int:
		w.WriteHeader(val)
	}
}
