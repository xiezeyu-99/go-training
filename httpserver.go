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

	go func() {
		<-ctx.Done()
		close(sigs)
	}()
	if sig := <-sigs; sig != nil {
		return errors.New(fmt.Sprint("signel:", sig))
	} else {
		return nil
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
		_ = server.Shutdown(context.Background())
	}()
	return server.ListenAndServe()
}
