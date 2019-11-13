package main

import (
	"flag"
	"fmt"
	"log"
	"nftserver/handle"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	go InitLog()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var addr string
	flag.StringVar(&addr, "l", ":9989", "-l ip:port")
	flag.Parse()

	go handle.ListenAndServe(addr)
	//go handle.ListenAndServe(":9989")
	//daemon.Daemon()
	select {
	case <-c:
		fmt.Println("system terminate!")
		log.Println("system terminate!")
	}

}

func InitLog() {
	var fd *os.File = nil
	var err error

	cur := time.Now()
	for {

		file := filepath.Join(os.Getenv("HOME"), "log",
			fmt.Sprintf("nftproxy-%s.log", cur.Format("2006-01-02")))

		fd1 := fd
		fd, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
		if err != nil {
			fmt.Println("cannot create log file")
			os.Exit(1)

		}

		log.SetOutput(fd)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		if fd1 != nil {
			_ = fd1.Close()
		}

		// 计算下一个零点
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		cur = <-t.C
	}
}

func startTimer(t time.Duration, f func()) {
	go func() {
		ticker := time.NewTicker(t)
		for {
			f()
			<-ticker.C
		}
	}()
}
