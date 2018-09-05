package lib

import (
	"fmt"
	"log"
	"sync"
)

type HttpServerMgr struct {
	servers map[int]*HttpServer
	size    int
	mutex   *sync.Mutex
}

var HttpServerMgrIns *HttpServerMgr
var once sync.Once

func GetHttpMgrInstance() *HttpServerMgr {
	once.Do(func() {
		HttpServerMgrIns = &HttpServerMgr{map[int]*HttpServer{}, 0, &sync.Mutex{}}
	})
	return HttpServerMgrIns
}

func (self *HttpServerMgr) GetHttpServer(port int) *HttpServer {
	if value, ok := self.servers[port]; ok {
		return value
	} else {
		log.Fatal("No Such HttpServer Linstend!")
		return nil
	}
}

func (self *HttpServerMgr) StopHttpServer(port int) error {
	return self.GetHttpServer(port).Stop()
}

func (self *HttpServerMgr) DeleteHttpServer(port int) error {
	err := self.GetHttpServer(port).Stop()
	if nil == err {
		self.mutex.Lock()
		delete(self.servers, port)
		self.size -= 1
		self.mutex.Unlock()
		return nil
	} else {
		return err
	}
}

func (self *HttpServerMgr) Start() {
	for _, s := range self.servers {
		if s.status != SERVING {
			err := s.Start()
			if nil != err {
				log.Fatal(fmt.Sprintf("Port[%v] start Failed!", s.port))
			}
		}
	}
}

func (self *HttpServerMgr) Stop() {
	for _, s := range self.servers {
		if s.status == SERVING {
			err := s.Stop()
			if nil != err {
				log.Fatal(fmt.Sprintf("Port[%v] stop Failed!", s.port))
			}
		}
	}
}
