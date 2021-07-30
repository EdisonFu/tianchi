package main

import (
	_ "go.uber.org/automaxprocs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tianchi/dao/cache"
	"tianchi/handler"
	_"tianchi/dao/mysql"
	_ "net/http/pprof"

	l4g "github.com/alecthomas/log4go"
)

func main() {
	l4g.LoadConfiguration("./log4go.xml")

	go func() {
		err := http.ListenAndServe(":6065", nil)
		if err != nil {
			l4g.Error("pprof err:%v", err)
			return
		}
		l4g.Info("pprof listen:6065")
	}()

	l4g.Info("server start！")
	//从数据库恢复数据
	cache.ReloadCacheFromDB()
	// 创建路由
	handler.InitRouter()

	//等待退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
}
