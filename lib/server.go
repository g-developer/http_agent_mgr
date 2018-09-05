package lib

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type HttpServer struct {
	port         int
	readTimeout  int
	writeTimeout int
	ins          *http.Server
	status       int
	mutex        *sync.Mutex
	IsRunning    chan int
}

type HttpHandler func(http.ResponseWriter, *http.Request)

func NewHttpServer(port int) *HttpServer {
	tmp := &HttpServer{
		port,
		10,
		10,
		&http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler:        http.NewServeMux(),
			ReadTimeout:    600 * time.Second,
			WriteTimeout:   600 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		INIT,
		&sync.Mutex{},
		make(chan int),
	}
	return tmp
}

func (self *HttpServer) stop(w http.ResponseWriter, r *http.Request) {
	self.Stop()
	w.WriteHeader(200)
	str := `{"code": 0, "msg": "success", "data": null}\n`
	w.Write([]byte(str))
}

func (self *HttpServer) AddDefaultHandler() {
	self.AddHandler("stop", self.stop)
}

func (self *HttpServer) AddFileServer(pattern string, dir string) error {
	//self.ins.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(dir))))
	if mux, ok := self.ins.Handler.(*http.ServeMux); ok {
		fmt.Println("AddHandler---", pattern)
		mux.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(dir))))
		return nil
	} else {
		return errors.New("Handler Is Not Type *http.ServeMux")
	}
}

func (self *HttpServer) AddHandler(pattern string, handler HttpHandler) error {
	if INIT == self.status || AVAILABLE == self.status || SERVING == self.status {
		if 0 != strings.Index(pattern, "/") {
			pattern = "/" + pattern
		}
		if mux, ok := self.ins.Handler.(*http.ServeMux); ok {
			fmt.Println("AddHandler---", pattern)
			self.mutex.Lock()
			mux.HandleFunc(pattern, handler)
			self.status = AVAILABLE
			self.mutex.Unlock()

			return nil
		} else {
			return errors.New("Handler Is Not Type *http.ServeMux")
		}
	} else {
		return errors.New(fmt.Sprintf("Port:%v Status Error! Status=%v", self.port, self.status))
	}

}

func (self *HttpServer) Start() error {
	fmt.Println("Start HttpServer Succ!")
	if AVAILABLE == self.status || STOP == self.status {

		var e error = nil
		self.mutex.Lock()
		go func(e error) {
			err := self.ins.ListenAndServe()
			if nil != err {
				self.status = AVAILABLE
				e = fmt.Errorf("Start Server Failed! Port=%v; Err=%v", self.port, err)
			} else {
				self.status = SERVING
				e = nil
			}
		}(e)
		self.mutex.Unlock()
		if nil != e {
			panic(e)
		}
		return e
	} else {
		e := errors.New(fmt.Sprintf("Port:%v Status Error! Status=%v", self.port, self.status))
		panic(e)
		return e
	}

}

func (self *HttpServer) Restart() error {
	err := self.Stop()
	if nil != err {
		return err
	} else {
		return self.Start()
	}
}

func (self *HttpServer) Stop() error {
	if SERVING == self.status {
		self.mutex.Lock()
		if SERVING == self.status {
			err := self.ins.Close()
			if nil == err {
				self.status = STOP
				self.IsRunning <- QUIT
			}
			self.mutex.Unlock()
			return err
		} else {
			self.IsRunning <- QUIT
			self.mutex.Unlock()
			return nil
		}
	} else {
		self.mutex.Lock()
		self.status = STOP
		self.IsRunning <- QUIT
		self.mutex.Unlock()
		return nil
	}
}
