package main

import (
	"github.com/g-developer/http_agent_mgr/lib"
	"net/http"
	"fmt"
	"time"
	"strconv"
	"os"
	"sync"
	"encoding/json"
)

type NgxApi interface {
	Serve(w http.ResponseWriter, r *http.Request)
	Strategy () string
}

type Api struct {
	Name string
	Balance string
	Per int
	Sleep int64
	Degree string
}

var mutex = &sync.Mutex{}
var allIps = map[string]*Api{}

func NewApi (balance string, name string, per int, sleep int64, degree string) *Api {
	return &Api{name, balance, per, sleep, degree}
}

func (self *Api) Strategy () string {
	res := ""
	switch self.Balance {
	case "iphash" : res = ipHash(self.Per)
	//case "fair" : return fair(self.Count)
	//case "" : return weight(self.Count)
	}
	return res
}

func (self *Api) Serve (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	strBody := self.Strategy()
	if 0 < self.Sleep {
		switch self.Degree {
		case "second" :
			time.Sleep(time.Second * time.Duration(self.Sleep))
		case "micro" :
			time.Sleep(time.Microsecond * time.Duration(self.Sleep))
		}
	}
	w.Write([]byte(strBody))
}


func ipHash (count int) string {
	str := "ip_hash;\n"
	for i:=0; i<count; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += "up;\n"
		}
	}
	return str
}


func getApis (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	resp, _ := json.Marshal(allIps)
	w.Write([]byte(resp))
}

func apis (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	r.ParseForm()
	prefix := "zyc"
	var sleep int64 = 0
	balance := ""
	degree := "second"
	count := 1
	per := 10
	if 0 < len(r.Form["prefix"]) {
		prefix = r.Form["prefix"][0]
	}
	if 0 < len(r.Form["balance"]) {
		balance = r.Form["balance"][0]
	}
	if 0 < len(r.Form["per"]) {
		per, _ = strconv.Atoi(r.Form["per"][0])
	}
	if 0 < len(r.Form["count"]) {
		count, _ = strconv.Atoi(r.Form["count"][0])
	}
	if 0 < len(r.Form["sleep"]) {
		sleep, _ = strconv.ParseInt(r.Form["sleep"][0], 10, 64)
	}
	if 0 < len(r.Form["degree"]) {
		degree = r.Form["degree"][0]
	}
	mutex.Lock()
	for i:=0; i<count; i++ {
		name := prefix + "." + strconv.Itoa(i)
		apiTmp := NewApi(balance, name, per, sleep, degree)
		go func(Name string, tmp *Api) {
			server.AddHandler(Name, tmp.Serve)
			allIps[name] = apiTmp
		}(name, apiTmp)
	}
	mutex.Unlock()
	resp := ""
	resp = `{"code": 0, "msg": "add "` + prefix + `" success"}\n`
	w.Write([]byte(resp))
}

var server *lib.HttpServer = nil

func main () {
	server = lib.NewHttpServer(8081)
	server.AddDefaultHandler()
	server.AddHandler("apis", getApis)
	err := server.AddHandler("add", apis)
	if nil != err {
		fmt.Println("Add Handler Error!", err)
	}
	go server.Start()

	for {
		select {
			case isQuit := <- server.IsRunning :
				if isQuit != lib.QUIT {
					fmt.Println("Server ---------- Running")
				} else {
					fmt.Println("Server ----------- Quit")
					os.Exit(0)
				}
			default:
				time.Sleep(1 * time.Second)
		}
	}

}
