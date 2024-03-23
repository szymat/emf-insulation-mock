package insulationmock_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/szymat/emf-insulation-mock/pkg/insulationmock"
)

func TestMockServerStart(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	scriptsPath := filepath.Join(dir, "../../tests/scripts/") + "/"

	config := insulationmock.NewConfig(scriptsPath, "routes.lua", "before.lua", "after.lua", "8989")

	r := mux.NewRouter()
	mockServer := insulationmock.NewMockServer(config)
	// r.HandleFunc("/*", mockServer.HandleRequest)
	// accept any path and method into HandleRequest
	r.PathPrefix("/").HandlerFunc(mockServer.HandleRequest)

	server := http.Server{
		Addr:              ":" + config.Port,
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server.ListenAndServe() returned an error: %v\n", err)
		}
		log.Print("server stopped")
	}()

	time.Sleep(time.Second)

	err := server.Shutdown(context.Background())
	if err != nil {
		t.Errorf("server.Shutdown() returned an error: %v", err)
	}

	log.Print("server stopped")
}

// func (ms *MockServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
// Load and sort routes from Lua configuration
// if err := ms.luaRouter.LoadRoutes(); err != nil {
// /        log.Fatalf("Failed to load routes: %v", err)
// }

//  log.Printf("Loaded %d routes", len(ms.luaRouter.routes))
//    ms.luaRouter.SortRoutes()
//      ms.handleRequest(w, r)
//}

func TestLoadRoutesFailed(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	scriptsPath := filepath.Join(dir, "../../tests/scripts/") + "/"

	config := insulationmock.NewConfig(scriptsPath, "routes.lua", "before.lua", "after.lua", "8989")

	mockServer := insulationmock.NewMockServer(config)

	// Create request and responseWriter for mock testing
	response := httptest.NewRecorder()

	body := []byte("test body")
	req := httptest.NewRequest("GET", "/", bytes.NewReader(body))

	mockServer.HandleRequest(response, req)
}

func TestMockServerHandleRequest(t *testing.T) {
	// Test only HandleRequest method
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	scriptsPath := filepath.Join(dir, "../../tests/scripts/") + "/"

	config := insulationmock.NewConfig(scriptsPath, "routes.lua", "before.lua", "after.lua", "8989")

	mockServer := insulationmock.NewMockServer(config)

	// Create request and responseWriter for mock testing
	response := httptest.NewRecorder()

	body := []byte("test body")
	req := httptest.NewRequest("GET", "/", bytes.NewReader(body))

	mockServer.HandleRequest(response, req)
}
