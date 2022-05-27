package main

import (
	"context"
	"errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Engine struct{}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second * 20)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func main() {
	engine := new(Engine)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      engine,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	gg := func(ctx context.Context) ([]Result, errot) {
		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			log.Println("HTTP服务器启动", "http://localhost"+server.Addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Println(err)
				return err
			}
			log.Println("HTTP服务关闭请求")
		})
		g.Go(func() error {
			log.Println("HTTP服务器启动", "http://localhost"+server.Addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Println(err)
				return err
			}
			log.Println("HTTP服务关闭请求")
		})

		if err := g.Wait(); err != nil {
			return nil, err
		}
	}

	// 监听信号 优雅退出http服务
	Watch(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})
}

func Watch(fns ...func() error) {
	// 程序无法捕获信号 SIGKILL 和 SIGSTOP （终止和暂停进程），因此 os/signal 包对这两个信号无效。
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	// 阻塞
	s := <-ch
	close(ch)
	log.Println("收到信号", s.String(), "执行关闭函数")
	for i := range fns {
		if err := fns[i](); err != nil {
			log.Println(err)
		}
	}
	log.Println("关闭函数执行完成")
}
