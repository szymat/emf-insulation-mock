package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/szymat/emf-insulation-mock/pkg/insulationmock"
)

func main() {
	luaScriptsPath := flag.String("scripts", "scripts/", "Path to the directory containing the Lua scripts")

	routesFile := flag.String("routes", "routes.lua", "Name of the Lua file containing route definitions")
	beforeScript := flag.String("before", "before.lua", "Name of the Lua file containing the before script")
	afterScript := flag.String("after", "after.lua", "Name of the Lua file containing the after script")

	port := flag.String("port", "8081", "Port on which the server will listen")

	flag.Parse()

	ctx := context.Background()
	// Create config
	config := insulationmock.NewConfig(*luaScriptsPath, *routesFile, *beforeScript, *afterScript, *port)

	log.Printf("Starting mock server with scripts path: %s", *luaScriptsPath)

	// check if directory exists
	if _, err := os.Stat(*luaScriptsPath); os.IsNotExist(err) {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}
}
