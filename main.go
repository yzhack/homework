package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	done := make(chan error, 2)
	stop := make(chan struct{})
	signalChan := make(chan os.Signal)
	go func() {
		done <- serve("127.0.0.1:8080", r, stop, &gin.Context{}, signalChan)
	}()
	go func() {
		done <- serve("127.0.0.1:8001", r, stop, &gin.Context{}, signalChan)
	}()

	go func() {
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
	}()
	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Println("error:%v", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}

func serve(addr string, r *gin.Engine, stop chan struct{}, ctx *gin.Context, signal chan os.Signal) (err error) {
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		////老师 为什么我这样写 测试下来 只能关闭一个server
		//select {
		//case <-stop:
		//case <-signal:
		//	_ = server.Shutdown(ctx)
		//}

		// 这样写就能正常关闭 2个server
		select {
		case <-stop:
			_ = server.Shutdown(ctx)
		case <-signal:
			_ = server.Shutdown(ctx)
		}
	}()
	return server.ListenAndServe()
}
