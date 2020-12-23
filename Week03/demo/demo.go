package main

import (
	"context"
	"fmt"
	"net/http"
)

// func serveApp() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintln(resp, "Hello World!")
// 	})
// 	http.ListenAndServe("0.0.0.0:8080", mux) // app traffic
// }

// func serveDebug() {
// 	go http.ListenAndServe("127.0.0.1:8081", http.DefaultServeMux) // debug;
// }

// func main() {
// 	go serveDebug()
// 	serveApp()
// }

func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-stop
		s.Shutdown(context.Background())
	}()
	return s.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello World!")
	})
	return serve("0.0.0.0:8080", mux, stop)
}

func serveDebug(stop <-chan struct{}) error {
	return serve("127.0.0.1:8081", http.DefaultServeMux, stop)
}

func main() {
	done := make(chan error, 2)
	stop := make(chan struct{})

	go func() {
		done <- serveApp(stop)
	}()
	go func() {
		done <- serveDebug(stop)
	}()

	var stoped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if !stoped {
			stoped = true
			close(stop)
		}
	}
}
