package tcp

import (
	"context"
	"gopro/interface/tcp"
	"gopro/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config,
	handler tcp.Handler) error {

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGILL, syscall.SIGQUIT, syscall.SIGTERM)
	// 异步监听系统信号，监听到关闭信号，就是用closeChan 发送一个关闭信息，让另一个异步线程感知到并做资源释放
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGILL, syscall.SIGQUIT, syscall.SIGTERM:
			logger.Warn(sig)
			closeChan <- struct{}{}
		}
	}()
	// 创建tcp连接
	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	ListenerAndServe(listen, handler, closeChan)
	return nil
}

func ListenerAndServe(listener net.Listener,
	handler tcp.Handler,
	closeChan <-chan struct{}) {

	// 监听closeChan通道， 比如调用Kill -9 或者系统杀死程序时 关闭资源
	go func() {
		<-closeChan
		logger.Info("shutting down 关闭资源")
		_ = listener.Close()
		_ = handler.Close()
	}()

	// 关闭资源
	defer func() {
		logger.Info("关闭资源")
		_ = listener.Close()
		_ = handler.Close()
	}()

	wg := sync.WaitGroup{}
	ctx := context.Background()
	for {

		// 循环监听多个客户端的存在.支持同时处理多个client的连接
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accept link")
		wg.Add(1)
		go func() {
			// 为了防止在Handler中发生panic， 所以在defer中进行wait的减少
			defer func() {
				wg.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}

	// 因为Server可能会服务多个client的连接，所以某个连接挂了的时候需要等待一下，等其他几个连接处理完
	wg.Wait()
}
