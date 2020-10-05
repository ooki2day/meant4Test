package httprouterserver

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

// HTTPRouterServer is an http server.
type HTTPRouterServer struct {
	endpointsMap map[string]httprouter.Handle
	router       *httprouter.Router
	srv          *http.Server
}

// NewServer creates a new HTTPRouterServer object.
func NewServer() *HTTPRouterServer {
	return &HTTPRouterServer{
		make(map[string]httprouter.Handle),
		httprouter.New(),
		nil}
}

// Start starts up a new server on port. Returns error if failed.
func (server *HTTPRouterServer) Start(port int) error {
	if server.srv != nil {
		server.Stop()
	}
	err := server.startREST()
	if err != nil {
		return err
	}
	server.srv = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: server.router}
	return server.srv.ListenAndServe()
}

// AddEndpoint adds an endpoint to server. Returns error if failed.
func (server *HTTPRouterServer) AddEndpoint(path string, handler httprouter.Handle) error {
	server.endpointsMap[path] = httprouter.Handle(handler)
	return nil
}

// RemoveEndpoint removes an endpoint from server. Returns error if failed.
func (server *HTTPRouterServer) RemoveEndpoint(path string) error {
	delete(server.endpointsMap, path)
	return nil
}

// Stop is a function to stop the sesrver. Returns error if failed.
func (server *HTTPRouterServer) Stop() error {
	var err error
	if server.srv != nil {
		log.Println("Shutdown server")
		server.router = httprouter.New()
		err = server.srv.Shutdown(context.Background())
		server.srv = nil
	}
	return err
}

func (server *HTTPRouterServer) startREST() error {
	if len(server.endpointsMap) < 1 {
		return errors.New("Endpoints is empty")
	}

	for point, handle := range server.endpointsMap {
		server.router.PUT(point, handle)
		log.Printf("Added %s endpoint", point)
	}

	return nil
}
