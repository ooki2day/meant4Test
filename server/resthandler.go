package httprouterserver

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"sync"
)

// factRequest is a type to parse a calculate request.
type factRequest struct {
	A int
	B int
}

type errorResponse struct {
	Error string
}

func getIncorrectInput() errorResponse {
	return errorResponse{"Incorrect input"}
}

// RESTHandler is a REST requests handler for HTTPRouterServer.
type RESTHandler struct {
	calculatedFactorials map[int]uint64
}

// NewRESTHandler creates a new handler.
func NewRESTHandler(calculatedFactorial map[int]uint64) *RESTHandler {
	if calculatedFactorial != nil {
		return &RESTHandler{calculatedFactorials: calculatedFactorial}
	}
	return &RESTHandler{make(map[int]uint64)}
}

// CalculateHandler is an endpoint /calculate handler.
func (handler *RESTHandler) CalculateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("method calculate")
	w.Header().Set("Content-Type", "application/json")

	req, err := parseRequest(r)
	if err != nil {
		log.Println(err)
		js, err := json.Marshal(getIncorrectInput())
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return
	}

	a, b, err := handler.calculateFactorial(req)
	if err != nil {
		log.Println(err)
		return
	}
	js, err := json.Marshal(map[string]uint64{"a": a, "b": b})
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(js)
}

func parseRequest(r *http.Request) (factRequest, error) {
	var req factRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if (err != nil) || (req.A < 1) || (req.B < 1) {
		return factRequest{}, errors.New(getIncorrectInput().Error)
	}
	return req, nil
}

func (handler *RESTHandler) calculateFactorial(req factRequest) (uint64, uint64, error) {
	var factA uint64
	var factB uint64
	var wg sync.WaitGroup
	wg.Add(2)
	var mut sync.Mutex

	go func() {
		defer wg.Done()

		if val, ok := handler.calculatedFactorials[req.A]; ok {
			factA = val
		} else {
			factA = factorial(req.A)
			mut.Lock()
			defer mut.Unlock()
			handler.calculatedFactorials[req.A] = factA
		}
	}()
	go func() {
		defer wg.Done()

		if val, ok := handler.calculatedFactorials[req.B]; ok {
			factB = val
		} else {
			factB = factorial(req.B)
			mut.Lock()
			defer mut.Unlock()
			handler.calculatedFactorials[req.B] = factB
		}
	}()

	wg.Wait()
	return factA, factB, nil
}

func factorial(n int) uint64 {
	if n == 1 {
		return 1
	}

	return uint64(n) * factorial(n-1)
}
