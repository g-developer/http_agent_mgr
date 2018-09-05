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
		"math/rand"
)

type NgxApi interface {
	Serve(w http.ResponseWriter, r *http.Request)
	Strategy () string
}

type Api struct {
	Name string
	Per int
	Sleep int64
	Degree string
}

var mutex = &sync.Mutex{}
var allIps = map[string]*Api{}

func NewApi (name string, per int, sleep int64, degree string) *Api {
	return &Api{name, per, sleep, degree}
}

func (self *Api) Strategy () string {
	res := normal(self.Per)
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


func normal (count int) string {
	str := ""
	for i:=0; i<count; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += ";\n"
		}
	}
	return str
}

func ipHash (count int) string {
	str := "ip_hash;\n"
	for i:=0; i<count; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += ";\n"
		}
	}
	return str
}


func getApis (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	resp, _ := json.Marshal(allIps)
	w.Write([]byte(resp))
}


//数据前几行换行
func firstEnter (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := "\n\n\n\n"
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += ";\n"
		}
	}
	w.Write([]byte(str))
}

//数据前几行换行
func firstBlank (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := `             `
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += ";\n"
		}
	}
	w.Write([]byte(str))
}

//数据中间换行
func middleEnter (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ""
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n\n\n"
		} else {
			str += ";\n"
		}
	}
	w.Write([]byte(str))
}

//数据中间为空
func middleBlank (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ""
	for i:=0; i<10; i++ {
		if i == 5 {
			str += fmt.Sprintf("           server       127.0.0.%v:8080 weight=2 ;\n", i)
		} else {
			if i == 8 {
				str += fmt.Sprintf("	server	127.0.0.%v:8080	weight=2	;\n", i)
			} else {
			str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
			if 0 == i % 2 {
				str += "down;\n"
			} else {
				str += ";\n"
			}
			}
		}
	}
	w.Write([]byte(str))
}

//数据多个分号
func multiFenhao (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ""
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;;;;;;;\n"
		} else {
			str += ";\n"
		}
	}
	w.Write([]byte(str))
}

//数据注释
func comments (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := "#test\n"
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down; #test----- \n"
		} else {
			str += ";\n"
		}
	}
	str += "#test------ server 127.0.0.%v:8080 weight=2\n"
	w.Write([]byte(str))
}

//所有数据一行
func oneline (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ""
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;"
		} else {
			str += ";"
		}
	}
	w.Write([]byte(str))
}


//超时过10s
func costgt10 (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ""
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += ";\n"
		}
	}
	time.Sleep(15 * time.Second)
	w.Write([]byte(str))
}

//偶尔超时
func randomtimeout (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ""
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", i)
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += ";\n"
		}
	}
	s := rand.Intn(15)
	time.Sleep(time.Duration(s) * time.Second)
	w.Write([]byte(str))
}

//总在变化
func random (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ""
	for i:=0; i<10; i++ {
		str += fmt.Sprintf("server 127.0.0.%v:8080 weight=2 ", rand.Intn(250))
		if 0 == i % 2 {
			str += "down;\n"
		} else {
			str += ";\n"
		}
	}
	w.Write([]byte(str))
}

//带策略
func ipHashApi (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	str := ipHash(10)
	w.Write([]byte(str))
}

func apis (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	r.ParseForm()
	prefix := "zyc"
	var sleep int64 = 0
	degree := "second"
	count := 1
	per := 10
	if 0 < len(r.Form["prefix"]) {
		prefix = r.Form["prefix"][0]
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
		apiTmp := NewApi(name, per, sleep, degree)
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
	server.AddHandler("firstenter", firstEnter)
	server.AddHandler("middlerenter", middleEnter)
	server.AddHandler("firstblank", firstBlank)
	server.AddHandler("middlerblank", middleBlank)
	server.AddHandler("comments", comments)
	server.AddHandler("oneline", oneline)
	server.AddHandler("costgt10", costgt10)
	server.AddHandler("iphash", ipHashApi)
	server.AddHandler("multiFenhao", multiFenhao)
	server.AddHandler("randomtimeout", randomtimeout)
	server.AddHandler("random", random)
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
