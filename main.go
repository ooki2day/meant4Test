package main

import (
	"bufio"
	"github.com/ooki2day/meant4Test/server"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	port = 8989
)

func main() {
	srv := httprouterserver.NewServer()
	restHander := httprouterserver.NewRESTHandler(nil)
	err := srv.AddEndpoint("/calculate", restHander.CalculateHandler)
	if err != nil {
		log.Fatal(err)
	}

	var serverStop sync.WaitGroup
	serverStop.Add(1)
	go func() {
		defer serverStop.Done()

		if err := srv.Start(port); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		if (len(text) > 3) && (strings.Compare(strings.ToLower(text[0:4]), "stop") == 0) {
			srv.Stop()
			break
		}
	}

	serverStop.Wait()
}
