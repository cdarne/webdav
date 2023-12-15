package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/webdav"
)

func main() {
	port := flag.Uint("p", 8100, "port to serve on")
	directory := flag.String("d", ".", "the directory for WebDAV to serve")
	flag.Parse()

	server := setupServer(*port, *directory)
	go serve(server)
	waitForSignal()
	shutdown(server)
	log.Println("server exited gracefully")
}

func setupServer(port uint, rootDirectory string) *http.Server {
	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	// fileHandler := http.FileServer(http.Dir(rootDirectory))
	webDAVHandler := &webdav.Handler{
		FileSystem: webdav.Dir(rootDirectory),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("webDAV error: %v\n", err)
			}
		},
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			// Some basic auth
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			username, password, authOK := r.BasicAuth()

			if authOK == false {
				http.Error(w, "Not authorized", 401)
				return
			}

			if username != *g_username || password != *g_password {
				http.Error(w, "Not authorized", 401)
				return
			}
		*/

		log.Printf("%s - %s %s %s %s\n", r.RemoteAddr, r.UserAgent(), r.Method, r.URL, r.Proto)

		// if r.Method == "GET" || r.Method == "HEAD" {
		// fileHandler.ServeHTTP(w, r)
		// } else {
		webDAVHandler.ServeHTTP(w, r)
		// }
	}))

	log.Printf("Serving %s on HTTP port: %d\n", rootDirectory, port)

	return server
}

func serve(server *http.Server) {
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen failed: %v\n", err)
	}
	log.Println("server stopped")
}

// Wait for interrupt signal to gracefully shutdown the server with
func waitForSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
}

func shutdown(server *http.Server) {
	ctx, stop := context.WithTimeout(context.Background(), time.Second)
	defer stop()

	log.Println("shutting down server")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v\n", err)
	}
}
