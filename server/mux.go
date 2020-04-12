package server

import (
	"io"
	"strings"
	"sync"
)

// NewMux instantiate a Mux.
func NewMux() *Mux {
	return &Mux{routes: make(map[string]Handler), rw: sync.RWMutex{}}
}

// Mux is a router.
type Mux struct {
	routes map[string]Handler
	rw     sync.RWMutex
}

// Serve method is a routing process.
func (m *Mux) Serve(r *Request, res ResponseWriter) {
	uri := m.parseURI(r.URI)

	m.rw.RLock()
	handler, ok := m.routes[uri]
	m.rw.RUnlock()

	if ok {
		handler.Serve(r, res)
	} else {

		m.notFound(r, res)
	}
}

func (m *Mux) parseURI(uri string) string {
	parts := strings.Split(uri, "/")
	return "/" + strings.Join(parts[1:], "/")
}

// HandleFunc is used to register a handler function on a path.
func (m *Mux) HandleFunc(pattern string, f func(*Request, ResponseWriter)) {
	m.rw.Lock()
	m.routes[pattern] = HandlerFunc(f)
	m.rw.Unlock()
}

func (m *Mux) notFound(r *Request, res ResponseWriter) {
	res.SetStatus("400")
	io.WriteString(res, "Handler Not Found!")
}

// HandlerFunc is adapter function between function and Handler interface.
type HandlerFunc func(*Request, ResponseWriter)

// Serve method is used to conform to Handler interface.
func (f HandlerFunc) Serve(r *Request, res ResponseWriter) {
	f(r, res)
}
