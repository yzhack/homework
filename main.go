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
	serversChan := make(chan *http.Server, 2)
	signalChan := make(chan os.Signal, 1)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return newServer(":9090", nil, serversChan).ListenAndServe()
	})
	g.Go(func() error {
		return newServer(":9091", nil, serversChan).ListenAndServe()
	})

	//注册信号量
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("done len:%d",len(serversChan))
			case <-signalChan:
				fmt.Printf("signalChan len:%d",len(serversChan))
				return stopServer(ctx, serversChan)
			}
		}
	})
	err := g.Wait()
	fmt.Printf("err:%s", err.Error())
}

func newServer(addr string, handler http.Handler, servers chan *http.Server) (server *http.Server) {
	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	servers <- s
	return s
}

func stopServer(ctx context.Context, servers chan *http.Server) (err error) {
	for v := range servers {
		if err := v.Shutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
