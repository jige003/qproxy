/*************************************************************************
   > File Name: proxy.go
   > Author: b0lu
   > Mail: b0lu_xyz@163.com
   > Created Time: 2016年01月13日 星期三 13时51分30秒
************************************************************************/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/elazarl/goproxy"
)

var (
	p     = fmt.Println
	pf    = fmt.Printf
	fpf   = fmt.Fprintf
	debug = new(bool)
)

func isFileExiste(filename string) bool {
	_, err := os.Stat(filename)
	return !(err != nil && os.IsNotExist(err))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Q_Stdout(req *http.Request) {
	p("\n---------------------start---------------------------")
	request_str, err := httputil.DumpRequest(req, false)
	check(err)
	p("URL==>", req.URL)
	p("HOST==>", req.Host)
	p("REQUESTURI==>", req.RequestURI)

	p("PATH==>", req.URL.Path)
	p("RAWQUERY==>", req.URL.RawQuery)

	p("\n---------------------------ALL REQUESTS INFO------------------------------")
	p(string(request_str))
	p("---------------------------ALL REQUESTS INFO------------------------------\n")
	p("---------------------end---------------------------\n")
}

func Q_log(req *http.Request, w io.Writer) {
	request_str, err := httputil.DumpRequest(req, false)
	check(err)
	w.Write(request_str)
}

type HttpHandle struct {
	req   chan *http.Request
	errch chan error
}

func NewHttpHandle() (*HttpHandle, error) {
	httpHandler := &HttpHandle{make(chan *http.Request), make(chan error)}
	var (
		f       *os.File
		logfile string
		err     error
	)
	logfile = "qproxylog"
	if isFileExiste(logfile) {
		f, err = os.OpenFile(logfile, os.O_RDWR|os.O_APPEND, 0777)
		check(err)
	} else {
		f, err = os.Create(logfile)
		check(err)
	}
	go func() {
		for m := range httpHandler.req {
			if *debug {
				Q_Stdout(m)
			}
			Q_log(m, f)
		}
		httpHandler.errch <- f.Close()
	}()
	return httpHandler, nil
}

func (httpHandler *HttpHandle) PutRequest(req *http.Request) {
	httpHandler.req <- req
}

func (httpHandler *HttpHandle) Close() error {
	close(httpHandler.req)
	return <-httpHandler.errch
}

func main() {
	verbose := flag.Bool("v", false, "proxy log info to stdout")
	addr := flag.String("l", ":9010", "proxy listen addr")
	debug = flag.Bool("d", false, "open the debug mode")
	//fmt.Printf("v1 type :%s\n", reflect.TypeOf(debug))
	flag.Usage = func() {
		p("================By b0lu===============")
		p("usage:  " + string(os.Args[0]) + " -v -l :9010 -d")
		p("-d	open then debug mode")
		p("-v	open proxy log info to stdout")
		p("-l	set the proxy listen addr and port. default: -l :9010")
		p("\nPS: proxy will log requests info into qproxylog file")
	}
	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			p(err)
		}
	}()

	p("================By b0lu===============")
	p("Load config.ini to init allowHostsRule and staticResourcesRules")
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
		return
	}
	allowHostsRule, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "allowHostsRule")
	check(err)
	p("allowHostsRule==>", allowHostsRule)
	allowRules := make([]*regexp.Regexp, 0)

	for _, hostRule := range strings.Split(allowHostsRule, ",") {
		allowRules = append(allowRules, regexp.MustCompile(hostRule))

	}
	p(allowRules)

	staticResources, err := cfg.GetValue(goconfig.DEFAULT_SECTION, "staticResources")
	check(err)
	p("staticResources==>", staticResources)
	staticResourcesRules := regexp.MustCompile(staticResources)
	p(staticResourcesRules)

	httpHandler, err := NewHttpHandle()
	check(err)
	proxy := goproxy.NewProxyHttpServer()

	proxy.Verbose = *verbose
	/*
		if err := os.MkdirAll("proxydump", 0755); err != nil {
			log.Fatal("Cant create dir", err)
		}
	*/
	proxy.OnRequest(goproxy.ReqHostMatches(allowRules...)).
		DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if !staticResourcesRules.MatchString(req.URL.Path) {
			httpHandler.PutRequest(req)
		}
		return req, nil
	})
	p("----------------------Proxy Start------------------------")
	http.ListenAndServe(*addr, proxy)

}
