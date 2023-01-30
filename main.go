package main

import (
	"flag"
	"log"
	"net/http"

	"golang.org/x/net/webdav"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "the directory for WebDAV to serve")
	flag.Parse()

	// fileHandler := http.FileServer(http.Dir(*directory))
	webDAVHandler := &webdav.Handler{
		FileSystem: webdav.Dir(*directory),
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

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
