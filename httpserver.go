package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		return signalHandle(ctx)
	})
	group.Go(func() error {
		return httpServer(ctx)
	})

	if err := group.Wait(); err != nil {
		fmt.Println("shutdown:", err)
	}
}

//信号处理
func signalHandle(ctx context.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	select {
	case <-ctx.Done():
	    return ctx.Err()
	case sig := <-sigs:
	    return errors.Errorf("get os signal: %v", sig)
	}

}

//http server
func httpServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World")
	})
	server := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(timeoutCtx)
	}()
	return server.ListenAndServe()
}
