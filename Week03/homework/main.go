package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func ServeAPP(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// serve app
	g.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello World!")
		})
		server := http.Server{
			Addr:    ":8082",
			Handler: mux,
		}

		go func() {
			select {
			case <-ctx.Done():
				fmt.Println("http ctx done")
			}
			server.Shutdown(context.Background())
		}()
		return server.ListenAndServe()
	})

	// server signal
	g.Go(func() error {
		quitSigs := []os.Signal{syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM} // SIGTERM is POSIX specific
		quit := make(chan os.Signal, len(quitSigs))
		signal.Notify(quit, quitSigs...)

		for {
			select {
			case s := <-quit:
				return fmt.Errorf("Got termial signal: %v", s)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	return g.Wait()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	err := ServeAPP(ctx)
	fmt.Println("err: ", err)
}
