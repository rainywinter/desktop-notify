package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/golang/glog"
)

var (
	port       int
	defaultMsg = "某个未知的作业结束啦，去看看吧"
	msgKey     = "message"
	msgCh      chan string
	repeat     = 3
	interval   = time.Second * 10
)

func init() {
	flag.IntVar(&port, "port", 8080, "http listening port.")
	msgCh = make(chan string, 100)
}

type server struct {
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	msg := r.FormValue(msgKey)
	if msg == "" {
		msg = defaultMsg
		glog.Info("某个任务发送了一个空的提醒")
	} else {
		glog.Infof("某个任务发来一条提醒: %s", msg)
	}
	select {
	case msgCh <- msg:
	default:
		glog.Infof("msg blocked: %s", msg)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func display() {
	defer func() {
		if r := recover(); r != nil {
			glog.Error(string(debug.Stack()))

			// restart job
			go display()
		}
	}()
	for msg := range msgCh {
		for i := 0; i < repeat; i++ {
			notify(msg)
			time.Sleep(interval)
		}
	}
}

func notify(msg string) {
	err := beeep.Notify("啦啦啦", msg, "")
	if err != nil {
		glog.Warning(err)
	}
}

func main() {
	flag.Parse()
	// err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	// if err != nil {
	// 	panic(err)
	// }
	go display()
	glog.Infof("server start: %d", port)
	glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), &server{}))
}
