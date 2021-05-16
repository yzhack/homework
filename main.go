package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	_ "golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "pong")
	})
	signalChan := make(chan os.Signal, 1)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return newServer(":9090", nil).ListenAndServe()
	})
	g.Go(func() error {
		return newServer(":9091", nil).ListenAndServe()
	})

	//注册信号量
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)

	//老师 为啥我这样 听不掉server进程
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
			case <-signalChan:
				if err := newServer("9090", nil).Shutdown(ctx); err != nil {
					return err
				}
				if err := newServer("9091", nil).Shutdown(ctx); err != nil {
					return err
				}
			}
		}
	})
	err := g.Wait()
	fmt.Printf("err:%s", err.Error())
}

func newServer(addr string, handler http.Handler) (server *http.Server) {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
