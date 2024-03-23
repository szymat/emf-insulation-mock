package insulationmock

import (
	"log"
	"net/http"
)

// MockServer represents the mock server.
type MockServer struct {
	Config    *config    // The configuration of the server
	luaRouter *LuaRouter // The router that handles Lua scripting and routing
}

// NewMockServer creates a new instance of MockServer.
func NewMockServer(config *config) *MockServer {
	log.Printf("Creating new mock server with port: %s and scripts path: %s", config.Port, config.ScriptsPath)

	if config.BeforeScript != "" {
		log.Printf("Before script: %s", config.BeforeScript)
	}

	if config.AfterScript != "" {
		log.Printf("After script: %s", config.AfterScript)
	}

	return &MockServer{
		Config:    config,
		luaRouter: NewLuaRouter(config),
	}
}

func (ms *MockServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Load and sort routes from Lua configuration
	if err := ms.luaRouter.LoadRoutes(); err != nil {
		log.Fatalf("Failed to load routes: %v", err)
	}

	log.Printf("Loaded %d routes", len(ms.luaRouter.routes))
	ms.luaRouter.SortRoutes()
	ms.handleRequest(w, r)
}

// handleRequest is the main HTTP handler function that delegates requests to the LuaRouter.
func (ms *MockServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	ms.luaRouter.FindAndExecuteRoute(w, r)
}
