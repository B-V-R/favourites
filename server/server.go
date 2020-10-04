/*
 * Copyright (c) 2020 BVR (Vighneswar Rao Bojja)
 * This file is subject to the terms and conditions defined in file 'LICENSE'.
 *
 */

package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Server interface for server
type Server interface {
	Start(address string)
	Register(path string, handlerFunc http.HandlerFunc)
	Shutdown(ctx context.Context) error
	ShutdownWithCallBack(fn func())
	Stop() error
}

type server struct {
	*http.Server
	*http.ServeMux
	log *log.Logger
}

// Start make the server listen and server at given address
func (s *server) Start(address string) {
	s.Addr = address
	s.Server.Handler = s.ServeMux
	s.log.Println(fmt.Sprintf("%s %s", "Server Started at", address))
	if err := s.ListenAndServe(); err != nil {
		s.log.Fatal("Failed to start server", err)
	}
}

// Register registers the path and handler function to server
func (s *server) Register(path string, handlerFunc http.HandlerFunc) {
	s.log.Println(fmt.Sprintf("server handler with route=%s is registered", path))
	s.HandleFunc(path, handlerFunc)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	s.ServeMux.ServeHTTP(w, r)
}

// ShutdownWithCallBack calls the given callback function before shutdowns the server
func (s *server) ShutdownWithCallBack(fn func()) {
	s.log.Println("Shutdown with Callback is registered")
	s.Server.RegisterOnShutdown(fn)
}

// Stop close the server
func (s *server) Stop() error {
	s.log.Println("Server is stopped")
	return s.Server.Close()
}

//New creates a new server
func New() Server {
	return &server{
		Server:   &http.Server{},
		ServeMux: http.NewServeMux(),
		log:      log.New(os.Stdout, "favourite-app-server:\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
