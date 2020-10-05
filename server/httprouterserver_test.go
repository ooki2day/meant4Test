package httprouterserver

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	port = 8989
)

func TestHTTPRouterServerPositive(t *testing.T) {
	srv := NewServer()
	restHander := NewRESTHandler(nil)
	err := srv.AddEndpoint("/calculate", restHander.CalculateHandler)
	if err != nil {
		t.Errorf("Failed to add an endpoint: %s", err)
		t.FailNow()
	}

	var serverStop sync.WaitGroup
	serverStop.Add(1)
	go func() {
		defer serverStop.Done()

		if err := srv.Start(port); err != http.ErrServerClosed {
			t.Errorf("Server runtime error: %s", err)
		}
	}()

	time.Sleep(2 * time.Second)
	srv.Stop()

	serverStop.Wait()
}

func TestHTTPRouterServerWOEndpoint(t *testing.T) {
	srv := NewServer()
	restHander := NewRESTHandler(nil)
	err := srv.AddEndpoint("/calculatetest", restHander.CalculateHandler)
	if err != nil {
		t.Errorf("Failed to add an endpoint: %s", err)
		t.FailNow()
	}

	err = srv.RemoveEndpoint("/calculatetest")
	if err != nil {
		t.Errorf("Failed to delete an endpoint: %s", err)
		t.FailNow()
	}

	var serverStop sync.WaitGroup
	serverStop.Add(1)
	go func() {
		defer serverStop.Done()

		if err := srv.Start(port); err != http.ErrServerClosed {
			if err.Error() != "Endpoints is empty" {
				t.Errorf("Server runtime error: %s", err)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	srv.Stop()

	serverStop.Wait()
}

func TestHTTPRouterServerPositiveStartTwice(t *testing.T) {
	srv := NewServer()
	restHander := NewRESTHandler(nil)
	err := srv.AddEndpoint("/calculateTest", restHander.CalculateHandler)
	if err != nil {
		t.Errorf("Failed to add an endpoint: %s", err)
		t.FailNow()
	}

	var serverStop sync.WaitGroup
	serverStop.Add(2)
	go func() {
		defer serverStop.Done()

		if err := srv.Start(port); err != http.ErrServerClosed {
			t.Errorf("Server runtime error: %s", err)
		}
	}()
	go func() {
		defer serverStop.Done()

		if err := srv.Start(port); err != http.ErrServerClosed {
			t.Errorf("Server runtime error: %s", err)
		}
	}()

	time.Sleep(2 * time.Second)
	srv.Stop()

	serverStop.Wait()
}
