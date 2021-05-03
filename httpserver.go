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

var stop chan struct{}

func main() {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(httpServer)
	group.Go(signalHandle)

	if err := group.Wait(); err != nil {
		close(stop)
		fmt.Println("shutdown:", err)
	}
}

//信号处理
func signalHandle() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-stop
		close(sigs)
	}()
	if sig := <-sigs; sig != nil {
		return errors.New(fmt.Sprint("signel:", sig))
	} else {
		return nil
	}

}

//http server
func httpServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World")
	})
	server := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux,
	}
	go func() {
		<-stop
		_ = server.Shutdown(context.Background())
	}()
	return server.ListenAndServe()
}
